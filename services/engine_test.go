package main

import (
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"services/jutzo"
	"services/jutzo/impl"
	"strconv"
	"testing"
)

type TestConfig struct {
	vars map[string]string
}

var PostgresURL = os.Getenv("TestPostgresURL")
var RedisURL = os.Getenv("TestRedisURL")
var NormalConfig = map[string]string{
	"JUTZO_DB_URL":      PostgresURL,
	"JUTZO_REDIS_URL":   RedisURL,
	"JUTZO_ADMIN_EMAIL": "bob@hablutzel.com",
	"JUTZO_ADMIN_PASS":  "test-pass",
	"JUTZO_ADMIN_USER":  "bob",
}

// GetConfigurationString returns the configuration item
// value. The second return determines if the item was
// actually provided. If the item was not provided the
// routine should return ("", false); otherwise it should
// provide (value, true)
func (c TestConfig) GetConfigurationString(name string) (string, bool) {
	value, ok := c.vars[name]
	return value, ok
}

// GetConfigurationInt is similar to GetConfigurationString
// but validates that the item is both provided and an integer
// If the item is not provided or not an integer, this should
// return (0, false); otherwise it should return (value, true)
func (c TestConfig) GetConfigurationInt(name string) (int, bool) {

	if value, isPresent := c.GetConfigurationString(name); isPresent {
		if retValue, err := strconv.Atoi(value); err == nil {
			return retValue, true
		}
	}
	return 0, false
}

func TestDBConnection(t *testing.T) {
	configuration := TestConfig{NormalConfig}
	_ = impl.NewPostgresConnection(configuration)
}

func TestMissingDBURL(t *testing.T) {
	configuration := TestConfig{map[string]string{}}
	db := impl.NewPostgresConnection(configuration)
	if err := db.Connect(); err == nil {
		t.Errorf("Testing lack of database URL failed")
	}
}

func TestUserCache(t *testing.T) {
	configuration := TestConfig{NormalConfig}
	if _, err := impl.NewRedisCache(configuration); err != nil {
		t.Errorf("Test creating cache failed: %s", err.Error())
	}

}

func deleteTables(directConnect *sql.DB, t *testing.T) {
	tablesToDelete := []string{
		"drop table if exists jutzo_database_info cascade",
		"drop table if exists jutzo_pending_validation cascade",
		"drop table if exists jutzo_registered_user cascade "}
	for _, statement := range tablesToDelete {
		if _, err := directConnect.Exec(statement); err != nil {
			t.Errorf("Unable to delete table %s: %s", statement, err.Error())
		}
	}
}

func directConnect(t *testing.T) *sql.DB {
	// Directly connect to the remote server
	if directConnect, err := sql.Open("postgres", PostgresURL); err != nil {
		t.Errorf("Could not connect to database on %s", PostgresURL)
		return nil
	} else {
		return directConnect
	}
}

func connectAndWipe(t *testing.T) *sql.DB {
	directConnect := directConnect(t)
	deleteTables(directConnect, t)
	return directConnect
}

func TestConnectingToEmptyDB(t *testing.T) {
	configuration := TestConfig{NormalConfig}

	// Directly connect to the remote server; blow away
	// any of the existing tables
	directConnect := connectAndWipe(t)

	// Now connect to the database with the engine object
	db := impl.NewPostgresConnection(configuration)
	if err := db.Connect(); err != nil {
		t.Errorf("Failed to connect to database: %s", err.Error())
	}

	// Once we do this the tables should be recreated. Check them.
	rows, err := directConnect.Query(`select table_name, column_name, data_type
												   from information_schema.columns
                                                  where table_name like 'jutzo%'
                                                  order by table_name, column_name`)
	if err != nil {
		t.Errorf("Could not get table information from db, %s", err.Error())
	} else {

		// Defer the row cleanup
		defer func(rows *sql.Rows) {
			if err := rows.Close(); err != nil {
				t.Errorf("Error closing row: %s", err.Error())
			}
		}(rows)

		// Define the expected results
		expectedResults := []map[string]string{
			{"table_name": "jutzo_database_info", "column_name": "schema_ordinal", "data_type": "integer"},
			{"table_name": "jutzo_pending_validation", "column_name": "username", "data_type": "character varying"},
			{"table_name": "jutzo_pending_validation", "column_name": "uuid", "data_type": "uuid"},
			{"table_name": "jutzo_registered_user", "column_name": "creation_time", "data_type": "timestamp without time zone"},
			{"table_name": "jutzo_registered_user", "column_name": "email", "data_type": "character varying"},
			{"table_name": "jutzo_registered_user", "column_name": "email_validated", "data_type": "boolean"},
			{"table_name": "jutzo_registered_user", "column_name": "password_hash", "data_type": "bytea"},
			{"table_name": "jutzo_registered_user", "column_name": "rights", "data_type": "text"},
			{"table_name": "jutzo_registered_user", "column_name": "username", "data_type": "character varying"},
		}

		//  Loop through the actual results
		for rows.Next() {
			var tableName, columnName, dataType string

			if len(expectedResults) == 0 {
				t.Errorf("Got more results than expected")
			}

			if err := rows.Scan(&tableName, &columnName, &dataType); err != nil {
				t.Errorf("Error getting results: %s", err.Error())
			} else {
				expectedResult := expectedResults[0]
				expectedResults = expectedResults[1:]
				if expectedResult["table_name"] != tableName ||
					expectedResult["column_name"] != columnName ||
					expectedResult["data_type"] != dataType {
					t.Errorf("Expected %s, %s, %s but got %s, %s, %s",
						expectedResult["table_name"], expectedResult["column_name"], expectedResult["data_type"],
						tableName, columnName, dataType)
				}
			}
		}

		// Make sure we got all our expected results
		if len(expectedResults) != 0 {
			t.Errorf("Got fewer results than expected")
		}

		// There shouldn't be an admin yet
		if adminCount, err := db.GetAdminCount(); err != nil {
			t.Errorf("Count not get admin count: %s", err.Error())
		} else {
			if adminCount != 0 {
				t.Errorf("Where did that admin come from ?")
			}
		}
	}
}

func TestConnectingToNewEngineWithoutAdminConfig(t *testing.T) {

	// Create an incomplete configuration. This doesn't have the
	// admin config required
	configuration := TestConfig{map[string]string{
		"JUTZO_DB_URL":    PostgresURL,
		"JUTZO_REDIS_URL": RedisURL,
	}}

	// Directly connect to the remote server; blow away
	// any of the existing tables
	_ = connectAndWipe(t)

	// Now connect to the database with the engine object
	db := impl.NewPostgresConnection(configuration)
	if err := db.Connect(); err != nil {
		t.Errorf("Failed to connect to database: %s", err.Error())
	}

	if cache, err := impl.NewRedisCache(configuration); err != nil {
		t.Errorf("Test creating cache failed")
	} else {

		// Connect to the database
		if _, err := impl.NewJutzoEngine(configuration, db, cache); err == nil {
			t.Errorf("The engine didn't report the problem creating the admin")
		}
	}
}

func TestConnectingToNewEngine(t *testing.T) {
	configuration := TestConfig{NormalConfig}

	// Directly connect to the remote server; blow away
	// any of the existing tables
	directConnect := connectAndWipe(t)

	// Now connect to the database with the engine object
	db := impl.NewPostgresConnection(configuration)
	if err := db.Connect(); err != nil {
		t.Errorf("Failed to connect to database: %s", err.Error())
	}

	if cache, err := impl.NewRedisCache(configuration); err != nil {
		t.Errorf("Test creating cache failed")
	} else {

		// Connect to the database
		if engine, err := impl.NewJutzoEngine(configuration, db, cache); err != nil {
			t.Errorf("Test creating engine failed")
		} else {
			// We should have an admin now
			if adminCount, err := db.GetAdminCount(); err != nil {
				t.Errorf("Count not get admin count: %s", err.Error())
			} else {
				if adminCount == 0 {
					t.Errorf("Admin was not properly created")
				} else {

					// Basic tests passed, we can call some more detailed one
					if userInfo := testRegisterUser(t, directConnect, engine, "test"); userInfo != nil {
						testValidateEmailFor(t, directConnect, engine, userInfo)
					}
				}
			}
		}
	}
}

func testValidateEmailFor(t *testing.T, db *sql.DB, engine jutzo.Engine, info jutzo.UserInfo) {

	// Te user should exist, so attempt to get the validation link
	if validationString, email, err := engine.CreateUniqueValidationForUser(info.GetUsername()); err != nil {
		t.Errorf("Unable to create a validation string: %s", err.Error())
	} else {

		// Make sure the email is valid
		if email != info.GetEmail() {
			t.Errorf("EMail %s from user info does not match %s from database", info.GetEmail(), email)
		}

		// Validate that it's a uuid
		if _, err := uuid.Parse(validationString); err != nil {
			t.Errorf("Error - validation string not uuid: %s", err.Error())
		} else {

			// Make sure that call didn't validate the email
			statement := `select email_validated from jutzo_registered_user where username = $1`
			row := db.QueryRow(statement, info.GetUsername())
			var emailValid bool
			if err = row.Scan(&emailValid); err != nil {
				t.Errorf("Failed to retrieve the results: %s", err.Error())
			} else {
				if emailValid {
					t.Errorf("For some reason the email is valid")
				}
			}

			// Now validate the user
			if err = engine.ValidateEmail(validationString); err != nil {
				t.Errorf("Could not validate user %s, uniqueID = %s", info.GetUsername(), validationString)
			} else {

				// Now it should be validated
				statement := `select email_validated from jutzo_registered_user where username = $1`
				row := db.QueryRow(statement, info.GetUsername())
				var emailValid bool
				if err = row.Scan(&emailValid); err != nil {
					t.Errorf("Failed to retrieve the results: %s", err.Error())
				} else {
					if !emailValid {
						t.Errorf("For some reason the email is valid")
					}
				}
			}

		}
	}

}

func testRegisterUser(t *testing.T, db *sql.DB, engine jutzo.Engine, username string) jutzo.UserInfo {

	t.Logf("Testing register user")
	const password = "test-pass"
	const email = "test@test.com"

	if status, userInfo, err := engine.RegisterUser(username, password, email); err != nil {
		t.Errorf("Error registering user %s: %s", username, err.Error())
		return nil
	} else {
		if status != jutzo.Success {
			t.Errorf("Unable to create %s: status = %d", username, status)
			return nil
		} else {
			if userInfo.GetUsername() != username || userInfo.GetEmail() != email {
				t.Errorf("User info returned is invalid")
				return nil
			}

			// Direct validation
			statement := `select username, email, password_hash, rights, email_validated
                            from jutzo_registered_user 
                           where username = $1`
			row := db.QueryRow(statement, username)
			var dbUsername, dbEmail, rights string
			var passwordHash []byte
			var emailValid bool
			if err = row.Scan(&dbUsername, &dbEmail, &passwordHash, &rights, &emailValid); err != nil {
				t.Errorf("Failed to retrieve the results: %s", err.Error())
			} else {
				if username != dbUsername {
					t.Errorf("Username does not match")
				}
				if email != dbEmail {
					t.Errorf("Email does not match")
				}
				if rights != "blog,login" {
					t.Errorf("Rights is not what we expected")
				}
				if err = bcrypt.CompareHashAndPassword(userInfo.GetPasswordHash(), []byte(password)); err != nil {
					t.Errorf("Password does not validate")
				}
				if emailValid {
					t.Errorf("For some reason the email is valid")
				}
			}

			return userInfo
		}
	}

}
