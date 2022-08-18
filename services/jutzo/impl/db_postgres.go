package impl

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math"
	"services/jutzo"
	"strings"
	"time"
)

type PostgresConnection struct {
	config jutzo.ConfigurationProvider
	db     *sql.DB
}

func NewPostgresConnection(configurationProvider jutzo.ConfigurationProvider) *PostgresConnection {
	result := new(PostgresConnection)
	result.config = configurationProvider
	return result
}

const SupportedSchema = 1

var UpgradeStatements = [...][]string{

	// Upgrade from schema 0 (non-existent) to schema 1
	{
		`drop table if exists jutzo_database_info cascade `,
		`drop table if exists jutzo_pending_validation cascade `,
		`drop table if exists jutzo_registered_user cascade `,
		`create table if not exists jutzo_registered_user
			(
			username        varchar(256)            not null
				constraint username_key
				primary key,
			email           varchar(256)            not null,
			email_validated boolean   default false not null,
			creation_time   timestamp default now() not null,
			password_hash   bytea                   not null,
			rights          text      default 'blog,login'::text
			)`,
		`alter table jutzo_registered_user owner to jutzo`,
		`create table if not exists jutzo_pending_validation
			(
			uuid     uuid default gen_random_uuid() not null
				constraint uuid_key
				primary key,
			username varchar(256)                   not null
				constraint foreign_key_name
				references jutzo_registered_user
				on update cascade on delete cascade
			)`,
		`alter table jutzo_pending_validation owner to jutzo`,
		`create unique index if not exists email_idx on jutzo_registered_user (email)`,
		`create table jutzo_database_info
			(
				schema_ordinal integer default 1
			)`,
		`alter table jutzo_database_info owner to jutzo`,
		`insert into jutzo_database_info (schema_ordinal) values (1)`,
	},
}

// Connect to the database. This should also do all structural
// validations and forced updates required in order for other
// routines to work successfully, including making sure there
// is at least one administrative user registered.
func (connection *PostgresConnection) Connect() error {

	// Do nothing if we are already connected
	if connection.db == nil {

		var err error
		if dbURL, isPresent := connection.config.GetConfigurationString("JUTZO_DB_URL"); isPresent {
			connection.db, err = sql.Open("postgres", dbURL)
			if err != nil {
				return err
			}

			// Log the connection with a redacted password
			redacted, _ := redactPassword(dbURL)
			log.Printf("Connected to database: %s", redacted)

			if version, err := connection.getEffectiveSchemaVersion(); err != nil {
				return err
			} else {

				// Make sure we have the right schema for this version of the system
				if version < SupportedSchema {
					if err = connection.upgradeFrom(version); err != nil {
						return err
					}
				}
				return nil
			}
		} else {
			return errors.New("required connection string for database JUTZO_DB_URL is missing")
		}
	} else {
		return nil
	}
}

// Shutdown the database and clean up any resources used
func (connection *PostgresConnection) Shutdown() error {
	db := connection.db
	if db != nil {
		connection.db = nil
		return db.Close()
	} else {
		return nil
	}
}

// getEffectiveSchemaVersion will return the effective version of the database,
// which will either be the value returned from the information table or zero
// if the information table cannot be found
func (connection *PostgresConnection) getEffectiveSchemaVersion() (int, error) {

	// See if we can get the version number from the information table.
	// The information table is the one table we can depend on always
	// being present; if the information table doesn't exist then we
	// are at effective schema zero (in other words, no schema exists yet)
	if connection.doesInfoTableExist() {

		// We have the information table, get the version number
		row := connection.db.QueryRow("select schema_ordinal from jutzo_database_info")
		var schemaOrdinal int
		switch err := row.Scan(&schemaOrdinal); err {
		case nil:
			return schemaOrdinal, nil
		default:
			return 0, errors.New("database version could not be retrieved")
		}
	} else {
		return 0, nil
	}
}

// doesInfoTableExist answers the existential question on whether the database exists
// or not. All the Jutzo data schemas will have a table 'jutzo_database_info' when
// minimally a column named 'schema_ordinal' of type integer. This field will give
// the schema ordinal (zero based) that the database should represent. If this field cannot
// be found, then the database is corrupt
func (connection *PostgresConnection) doesInfoTableExist() bool {

	sqlStatement := `SELECT data_type from information_schema.columns
                            where table_name = 'jutzo_database_info' 
                              and column_name = 'schema_ordinal'`
	row := connection.db.QueryRow(sqlStatement)
	var dataType string
	switch err := row.Scan(&dataType); err {
	case nil:
		return dataType == "integer"
	case sql.ErrNoRows:
		return false
	default:
		log.Printf("Error validating information table in DB: %s", err.Error())
		return false
	}
}

// upgradeFrom the reported version to the current version of the
// schema. This will execute a set of statements that changes the
// database into the version required for this version of the server
func (connection *PostgresConnection) upgradeFrom(version int) error {

	db := connection.db
	log.Printf("Upgrading database, please wait...")
	for targetVersion, statements := range UpgradeStatements[version:] {

		// Execute all the statements in that upgrade set
		for _, statement := range statements {
			_, err := db.Exec(statement)
			if err != nil {
				return err
			}
		}
		log.Printf("Upgraded to version %d", targetVersion+1)
	}
	return nil
}

// CheckForUsernameOrEmail in the database so that we don't use the
// same username or email twice
func (connection *PostgresConnection) CheckForUsernameOrEmail(username string, email string) (bool, bool, error) {

	// Look for either the username or email in a single-shot query
	query := `select userQuery.userCount, emailQuery.emailCount from
    			(select count(*) userCount 
    			   from jutzo_registered_user j1 
    			  where j1.username = $1) as userQuery,
    			(select count(*) emailCount 
    			   from jutzo_registered_user j2 
    			  where j2.email = $2) as emailQuery`

	// Execute the query
	row := connection.db.QueryRow(query, username, email)
	var userCount int
	var emailCount int
	if err := row.Scan(&userCount, &emailCount); err == nil {
		return userCount > 0, emailCount > 0, nil
	} else {
		return false, false, err
	}
}

// StoreUser in the database with the given username, email and password hash. This
// routine will return an error if either the email or username are already in the
// database, so checking first with CheckForUsernameOrEmail is a good idea
func (connection *PostgresConnection) StoreUser(username string, email string, passwordHash []byte) (jutzo.UserInfo, error) {
	statement := `insert into jutzo_registered_user
                              (username, email,password_hash)
                       values ($1, $2, $3)
                     returning rights, creation_time`

	row := connection.db.QueryRow(statement, username, email, passwordHash)
	var rightsString string
	var creationTime time.Time
	if err := row.Scan(&rightsString, &creationTime); err == nil {
		return NewUserInfo(username, email, passwordHash, false, strings.Split(rightsString, ","), creationTime), err
	} else {
		return nil, err
	}
}

// UpdateUserInfo that has changed with what is stored in the database
func (connection *PostgresConnection) UpdateUserInfo(userInfo jutzo.UserInfo) error {

	statement := `update jutzo_registered_user set rights = $1 where username = $2`

	// TODO - anything more than rights
	rightsString := ""
	separator := ""
	for _, right := range userInfo.GetAllRights() {
		rightsString = fmt.Sprintf("%s%s%s", rightsString, separator, right)
		separator = ","
	}
	_, err := connection.db.Exec(statement, rightsString, userInfo.GetUsername())
	return err
}

// RetrieveUserInformation for the specified username so that the user credentials can
// be validated
func (connection *PostgresConnection) RetrieveUserInformation(username string) (jutzo.UserInfo, error) {
	statement := `SELECT email, email_validated, creation_time, 
                            password_hash, email_validated, rights 
                       from jutzo_registered_user
                      where username = $1`
	row := connection.db.QueryRow(statement, username)

	var email string
	var emailValidated bool
	var creationTime time.Time
	var passwordHash []byte
	var rightsString string
	if err := row.Scan(&email, &emailValidated, &creationTime, &passwordHash, &emailValidated, &rightsString); err != nil {
		return nil, err
	} else {
		// Build and return the user info
		return NewUserInfo(username, email, passwordHash, emailValidated, strings.Split(rightsString, ","), creationTime), nil
	}
}

// CreateValidationFor the user specified, so that the user can
// validate they actually have access to the email
func (connection *PostgresConnection) CreateValidationFor(username string) (uniqueID string, email string, err error) {

	deleteStatement := `delete from jutzo_pending_validation where username = $1`
	insertStatement := `insert into jutzo_pending_validation (username) 
                             values ($1) 
                          returning uuid, (select email 
                                             from jutzo_registered_user 
                                            where jutzo_registered_user.username = $1)`

	db := connection.db

	// First delete any pending validations for this user
	if _, err = db.Exec(deleteStatement, username); err == nil {

		// Now we can insert into the validation table and get the UUID that is generated. We
		// structure this as a query because we're going to return the UUID that gets created
		// automatically by the DB server as a part of the insert
		row := db.QueryRow(insertStatement, username)
		err = row.Scan(&uniqueID, &email)
		return
	} else {
		return "", "", err
	}

}

// CompleteValidationFor the uniqueID created with the
// CreateValidationFor method
func (connection *PostgresConnection) CompleteValidationFor(uniqueID string) error {

	db := connection.db
	updateStatement := `update jutzo_registered_user set email_validated = true where username = 
                               (select username from jutzo_pending_validation where uuid = $1)`
	deleteStatement := `delete from jutzo_pending_validation where uuid = $1`

	// Mark the user's email as valid based on the unique ID provided
	if _, err := db.Exec(updateStatement, uniqueID); err == nil {

		// Delete the pending validation record
		_, err = db.Exec(deleteStatement, uniqueID)
		return err
	} else {
		return err
	}

}

// ListUsers in the database, starting with the specified user, until maxUsers are returned.
// If the starting user is specified as "", then the list will start at the beginning of the
// users in the database; otherwise it can be used for pagination through the set of users.
// The first user returned will be the first user AFTER the one specified, so duplicate records
// will not occur. If no more users can be found, a nil slice will be returned with no error
func (connection *PostgresConnection) ListUsers(startingAt string, maxUsers int) ([]jutzo.UserInfo, error) {
	sqlStatement := `SELECT username, email, email_validated, creation_time, email_validated, rights
                       from jutzo_registered_user
                       where username > $1
                       order by username
                       limit $2`

	// Get the max users to return - or MaxInt if the user specifies an invalid one
	if maxUsers <= 0 {
		maxUsers = math.MaxInt32
	}

	if rows, err := connection.db.Query(sqlStatement, startingAt, maxUsers); err == nil {

		defer func(rows *sql.Rows) {
			if err := rows.Close(); err != nil {
				log.Printf("Error closing row: %s", err.Error())
			}
		}(rows)
		var username string
		var result []jutzo.UserInfo
		for rows.Next() {
			var email string
			var emailValidated bool
			var creationTime time.Time
			var rightsString string

			if err := rows.Scan(&username, &email, &emailValidated, &creationTime, &emailValidated, &rightsString); err != nil {
				return nil, err
			} else {
				result = append(result, NewUserInfo(username, email, []byte{},
					emailValidated, strings.Split(rightsString, ","), creationTime))
			}
		}

		if len(result) < maxUsers {
			return result, nil
		} else {
			return result, nil
		}
	} else {
		return nil, err
	}

}

// GetAdminCount returns the number of administrator users
func (connection *PostgresConnection) GetAdminCount() (count int, err error) {
	query := `select count(*) from jutzo_registered_user where rights like '%admin%'`

	// Count the admins
	row := connection.db.QueryRow(query)
	err = row.Scan(&count)
	return
}
