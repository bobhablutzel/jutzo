package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"reflect"
)

const InvalidDatabase = -2

type columnDesc map[string]string

type schema map[string][]columnDesc

// Define an expected schema for each version of the database
// known (currently just V1). This is used when validating the
// structure of the database.
var expectedSchemas = [...]schema{schemaForV1}

// DDL for creating the current schema if it doesn't exist
var currentDDL = ddlV1

// Current version of the schema - found by introspection after the
// database is loaded
var currentSchemaVersion int

// Connect to the remote database. This will return the database connection
// or an error of the database couldn't be opened for some reason. Note that
// one of the reasons the database could fail to open is a schema validation failure
func connectToDatabase(initDB bool, forceInit bool) (*sql.DB, error) {

	dbURL := getRequiredConfigString("JUTZO_DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.New("unable to connect to database host")
	}

	// Log the connection with a redacted password
	redacted, _ := redactPassword(dbURL)
	log.Printf("Connected to database: %s", redacted)

	if currentSchemaVersion, err = introspectDBVersion(db); err != nil {
		if initDB {
			err = recreateDatabase(db, forceInit)
			return db, err
		} else {
			shutdownDatabase(db)
			return nil, err
		}
	} else if isKnownSchema(currentSchemaVersion) {
		return db, nil
	} else {
		if initDB {
			err = recreateDatabase(db, forceInit)
			return db, err
		} else {
			shutdownDatabase(db)
			return nil, errors.New("database has an unsupported schema version - server might need to be upgraded")
		}
	}
}

func doRecreationDDL(db *sql.DB) error {
	log.Printf("Recreating database...")
	for _, query := range currentDDL {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

// Used to recreate the database when there is a schema mismatch.
// This is a highly destructive operation; it does not attempt
// to save any existing data.
func recreateDatabase(db *sql.DB, forceInit bool) error {
	if !forceInit {
		var answer string
		fmt.Print("Initializing the database may lose data, including content and users. Proceed? [N/y] ")
		_, err := fmt.Scanln(&answer)

		if err == nil && answer == "y" {
			return doRecreationDDL(db)
		} else {
			return errors.New("user declined to create database")
		}
	} else {
		return doRecreationDDL(db)
	}
}

// Check to see if the schema is a known schema
func isKnownSchema(schema int) bool {
	return schema == schemaV1
}

// All the Jutzo data schemas will have a table 'jutzo_database_info' when
// minimally a column named 'schema_ordinal' of type integer. This field will give
// the schema ordinal (zero based) that the database should represent. If this field cannot
// be found, then the database is corrupt
func validateInfoTableExists(db *sql.DB) bool {

	sqlStatement := `SELECT data_type from information_schema.columns
                            where table_name = 'jutzo_database_info' 
                              and column_name = 'schema_ordinal'`
	row := db.QueryRow(sqlStatement)
	var dataType string
	switch err := row.Scan(&dataType); err {
	case nil:
		return dataType == "integer"
	default:
		log.Printf("Error validating information table in DB: %s", err.Error())
		return false
	}
}

// Using the 'jutzo_database_info' table, get the version of the
// database schema that should match this database. This routine
// will return either the schema version, or InvalidDatabase if the
// version number cannot be obtained.
func getDatabaseVersionFromDatabase(db *sql.DB) (int, error) {

	if validateInfoTableExists(db) {
		row := db.QueryRow("select schema_ordinal from jutzo_database_info")
		var schemaOrdinal int
		switch err := row.Scan(&schemaOrdinal); err {
		case nil:
			return schemaOrdinal, nil
		default:
			return InvalidDatabase, errors.New("database version could not be retrieved")
		}
	} else {
		return InvalidDatabase, errors.New("database info table missing - database might need to be initialized")
	}

}

// Validates the schema and returns the schema version number
//
// This works by first attempting to get a schema version from
// the database. Once that is obtained, the table/column information
// will be returned for all tables named 'jutzo%' (all Jutzo
// tables will follow that naming convention). The table/column names
// returned will then be used to build up a schema representation that
// can be used to compare with the expected schema.
func introspectDBVersion(db *sql.DB) (int, error) {

	if schemaOrdinal, err := getDatabaseVersionFromDatabase(db); err != nil {
		return InvalidDatabase, err
	} else {
		// Make sure it's a schema version we're aware of
		if schemaOrdinal >= 0 && schemaOrdinal < len(expectedSchemas) {

			// Get the expected schema definitions for this version of the database
			expectedSchema := expectedSchemas[schemaOrdinal]

			// Make the structure we'll save the fields we find in
			actualSchema := make(schema)

			// Get the schema from the database information tables
			sqlStatement := `SELECT table_name, column_name, data_type 
                           from information_schema.columns 
                          where table_name like 'jutzo%' 
                          order by table_name, column_name`

			// Loop through all the information about the tables and build the actual
			// schema from the information in the database
			if rows, err := db.Query(sqlStatement); err == nil {
				defer closeRows(rows)
				for rows.Next() {
					var tableName, columnName, dataType string
					if err = rows.Scan(&tableName, &columnName, &dataType); err == nil {
						if columns, ok := actualSchema[tableName]; ok {
							actualSchema[tableName] = append(columns, columnDesc{columnName: dataType})
						} else {
							actualSchema[tableName] = []columnDesc{{columnName: dataType}}
						}
					} else {
						log.Printf("Error asking for schema structure: %s", err.Error())
					}
				}

				// If we didn't get any errors iterating, then check to see if the
				// schemas match
				if rows.Err() == nil {
					if reflect.DeepEqual(expectedSchema, actualSchema) {
						log.Printf("Schema validated as ordinal %d", schemaOrdinal)
						return schemaOrdinal, nil
					} else {
						return InvalidDatabase, errors.New(
							fmt.Sprintf("schema does not match expectations for schema #%d", schemaOrdinal))
					}
				} else {
					return InvalidDatabase, errors.New(
						fmt.Sprintf("Error asking for schema structure: %s", rows.Err().Error()))
				}
			} else {
				return InvalidDatabase, errors.New(
					fmt.Sprintf("Error asking for schema structure: %s", err.Error()))
			}
		} else {
			return InvalidDatabase, errors.New(fmt.Sprintf("schema %d is invalid", schemaOrdinal))
		}
	}
}

// Routine to shutdown the database
func shutdownDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("Unable to successfully close DB")
	} else {
		log.Println("Shut down DB")
	}
}

// Routine to safely close down a Rows instance
func closeRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Printf("Unable to close a row: %s", err.Error())
	}
}

// Generic function to see if there is one and only row in the database
func checkForUniqueRow(db *sql.DB, query string, args ...any) (bool, error) {
	return checkRowCount(db, func(count int) bool { return count == 1 }, query, args...)
}

func checkForAtLeastOneRow(db *sql.DB, query string, args ...any) (bool, error) {
	return checkRowCount(db, func(count int) bool { return count > 0 }, query, args...)
}

func checkRowCount(db *sql.DB, tester func(int) bool, query string, args ...any) (bool, error) {
	row := db.QueryRow(query, args...)
	var count int
	switch err := row.Scan(&count); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return tester(count), nil
	default:
		return false, err
	}
}
