package main

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Users struct {
	Username string `json:"username"`
}

// Validate if a user exists or not.
func checkIfUsernameIsInDatabase(db *sql.DB, user string) (bool, error) {
	return checkForUniqueRow(db, "select count(*) from jutzo_registered_user where username = $1", user)
}

// Validate if an email is already in use in the system
func checkIfEmailIsInDatabase(db *sql.DB, email string) (bool, error) {
	return checkForUniqueRow(db, "select count(*) from jutzo_registered_user where email = $1", email)
}

// Ensure there is at least one administrative user in the database. Force it if we need to.
func ensureAdministrativeUser(db *sql.DB) (bool, error) {
	exists, err := checkForAtLeastOneRow(db, `select count(*) from jutzo_registered_user where rights like '%admin%'`)
	if err != nil {
		return false, err
	} else {
		if !exists {
			return true, forceRegisterAdminUser(db)
		} else {
			return true, nil
		}
	}
}

// Routine to register a new user in the database
//
// Return values:
//    Status: StatusInserted ==> User was inserted
//    Status: StatusDuplicateUsername ==> User not inserted, username exists
//    Status: StatusDuplicateEmail ==> User not inserted, email exists
const (
	StatusInserted = iota
	StatusDuplicateUsername
	StatusDuplicateEmail
)

func registerUser(db *sql.DB, user string, password string, email string) (int, error) {
	// See if the user already exists
	if exists, err := checkIfUsernameIsInDatabase(db, user); err == nil {
		if exists {
			return StatusDuplicateUsername, nil
		} else {

			// Check to see if the email already exists
			if exists, err = checkIfEmailIsInDatabase(db, email); err == nil {
				if exists {
					return StatusDuplicateEmail, nil
				} else {
					err = insertUser(db, user, password, email)
					return StatusInserted, err
				}
			} else {
				return 0, err
			}
		}
	} else {
		return 0, err
	}
}

func forceRegisterAdminUser(db *sql.DB) error {

	// Register the admin user
	user := getRequiredConfigString("JUTZO_ADMIN_USER")
	_, err := registerUser(db,
		user,
		getRequiredConfigString("JUTZO_ADMIN_PASS"),
		getRequiredConfigString("JUTZO_ADMIN_EMAIL"))
	if err == nil {
		log.Printf("Making %s an administrator", user)
		// Update the admin user to show that it's valid and an administrator
		sqlStatement := `update jutzo_registered_user 
                            set email_validated = true,
                                rights = 'admin'
                         where username = $1`
		_, err = db.Exec(sqlStatement, user)
	}
	return err

}

// Actually insert a user into the database. This will
// hash the password, insert the user, and trigger the
// email validation routine.
//
// This routine DOES NOT validate whether the user is already
// in the database, and so should only be called after that
// validation check is performed.
func insertUser(db *sql.DB, user string, pass string, email string) error {

	// Get the cost for creating the password. This cost can be
	// set as an environment variable or defaulted, and it will
	// be stored with the password for validation later
	cost := getConfigInt("JUTZO_HASH_COST", 15)

	// Encrypt the password provided
	if bytes, err := bcrypt.GenerateFromPassword([]byte(pass), cost); err != nil {
		return err
	} else {

		// Now we can insert the user into the registered user table
		sqlStatement := `insert into jutzo_registered_user
                                     (username, email, password_cost, password_hash)
                         values ($1, $2, $3, $4)`
		_, err := db.Exec(sqlStatement, user, email, cost, bytes)
		return err
	}
}

// Used to list all users (for an administrator)
// Note this does not check if the caller is an administrator
func listUsers(db *sql.DB) ([]Users, error) {
	sqlStatement := `SELECT username from jutzo_registered_user;`

	var users []Users
	if rows, err := db.Query(sqlStatement); err == nil {
		defer closeRows(rows)
		for rows.Next() {
			var username string
			if err = rows.Scan(&username); err == nil {
				users = append(users, Users{username})
			} else {
				return nil, err
			}
		}

		// If we didn't get any errors iterating, the user array
		return users, rows.Err()
	} else {
		return nil, err
	}
}

// Routine to create a new pending validation record and return that to the
// requester as UUID that can be used to satisfy the validation
func createUniqueValidationForUser(db *sql.DB, user string) (string, error) {

	// First delete any pending validations for this user
	if _, err := db.Exec("delete from jutzo_pending_validation where username = $1", user); err == nil {

		// Now we can insert into the validation table and get the UUID that is generated. We
		// structure this as a query because we're going to return the UUID that gets created
		// automatically by the DB server as a part of the insert
		row := db.QueryRow("insert into jutzo_pending_validation (username) values ($1) returning uuid", user)
		var uuid string
		switch err := row.Scan(&uuid); err {
		case nil:
			return uuid, nil
		default:
			return "", err
		}
	} else {
		return "", err
	}
}

// Routine to validate an email. This essentially takes the UUID given and
// (1) deletes the pending validation record, and (b) marks the email as validated
func validateEmail(db *sql.DB, uuid string) (bool, error) {

	// Get the username associated with this uuid
	row := db.QueryRow("select username from jutzo_pending_validation where uuid = $1", uuid)
	var username string
	switch err := row.Scan(&username); err {
	case nil:

		// Valid request for validation, so mark the user's email as valid.
		if _, err := db.Exec("update jutzo_registered_user set email_validated = true where username = $1", username); err == nil {
			// Delete the pending validation record
			_, err := db.Exec("delete from jutzo_pending_validation where uuid = $1", uuid)
			return err == nil, err
		} else {
			return false, err
		}
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

// Routine to validate that the password provided in clear text matches
// the hashed password in the database
func login(db *sql.DB, user string, password string) (string, error) {

	// Get the password hash for the user
	row := db.QueryRow("select password_hash, email_validated, rights from jutzo_registered_user where username = $1", user)

	var passwordHash []byte
	var emailValidated bool
	var rightsString string
	switch err := row.Scan(&passwordHash, &emailValidated, &rightsString); err {
	case sql.ErrNoRows:
		return "", errors.New(fmt.Sprintf("User not found in database: %s", user))
	case nil:
		{
			// Make sure that the user is allowed to log in
			rights := NewRights(rightsString)
			canLogin := HasLogin(rights) || HasAdmin(rights)
			if emailValidated && canLogin {

				// User is able to log in, compare the password hash
				if err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password)); err == nil {
					if token, err := createUserSession(user, rights); err == nil {
						return token, nil
					} else {
						return "", err
					}
				} else {
					return "", err
				}
			} else {
				return "", errors.New(fmt.Sprintf("User %s is not active - unvalidated or no login rights", user))
			}
		}
	default:
		return "", err
	}
}
