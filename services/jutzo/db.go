package jutzo

import (
	_ "github.com/lib/pq"
)

type DatabaseConnection interface {

	// Connect to the database. This should also do all structural
	// validations and forced updates required in order for other
	// routines to work successfully, including making sure there
	// is at least one administrative user registered. This is a no-op
	// if already connected
	Connect() error

	// Shutdown the database and clean up any resources used
	Shutdown() error

	// CheckForUsernameOrEmail in the database so that we don't use the
	// same username or email twice
	CheckForUsernameOrEmail(username string, email string) (userExists bool, emailExists bool, err error)

	// StoreUser in the database with the given username, email and password hash
	StoreUser(username string, email string, passwordHash []byte) (UserInfo, error)

	// UpdateUserInfo that has changed with what is stored in the database
	UpdateUserInfo(userInfo UserInfo) error

	// RetrieveUserInformation for the specified username so that the user credentials can
	// be validated
	RetrieveUserInformation(username string) (userInfo UserInfo, err error)

	// GetAdminCount returns the number of administrator users
	GetAdminCount() (int, error)

	// CreateValidationFor the user specified, so that the user can
	// validate they actually have access to the email. The email that is
	// in the database will be returned, along with a unique identifier that
	// can be passed to CompleteValidationFor
	CreateValidationFor(username string) (uniqueID string, email string, err error)

	// CompleteValidationFor the uniqueID created with the
	// CreateValidationFor method
	CompleteValidationFor(uniqueID string) error

	// ListUsers in the database, starting with the specified user, until maxUsers are returned.
	// If the starting user is specified as "", then the list will start at the beginning of the
	// users in the database; otherwise it can be used for pagination through the set of users.
	// The first user returned will be the first user AFTER the one specified, so duplicate records
	// will not occur. If no more users can be found, a nil slice will be returned with no error
	ListUsers(startingAt string, maxUsers int) ([]UserInfo, error)
}
