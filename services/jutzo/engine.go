// Engine for the Jutzo environment. Encapsulates the internal connections
// to the database, Redis server, etc.

package jutzo

const (
	Success           = 0
	DuplicateEmail    = 1
	DuplicateUsername = 2
)

type IterationStatus interface {
}

// Engine defines the engine interface
type Engine interface {

	// Shutdown ensures the engine has a chance to close all it's internal connections
	Shutdown()

	Login(user string, password string) (UserSession, error)

	// RegisterUser sets up a new user in the user management system. This will return
	// one of (Success, DuplicateEmail, DuplicateUsername) depending on whether the
	// username and email are unique, or whether one of the email or username is duplicated.
	// If both the username and email are duplicated the DuplicateUsername will be
	// returned.
	RegisterUser(user string, password string, email string) (result int, userInfo UserInfo, err error)

	// CreateUniqueValidationForUser creates a new validation request record
	// that can be satisfied by a call to ValidateEmail
	CreateUniqueValidationForUser(user string) (uniqueID string, email string, err error)

	// ValidateEmail is called when a user responds to an email to the given
	// email address that contains the unique ID created by CreateUniqueValidationForUser
	ValidateEmail(uniqueID string) error

	// DestroyUserSession kills an active user session
	DestroyUserSession(uniqueID string) error

	// LoadUserSession returns the UserSession instance associated with the
	// given unique identifier.
	LoadUserSession(uniqueID string) (UserSession, error)

	// ListUsers can be called by an admin to list users. This call will
	// return up to maxUsers users at a time; if you want to retrieve the
	// remaining users call ListUsers again passing in the username from the
	// last record returned. If no more users can be found, an empty slice
	// will be returned with no errors.
	// If you want to start at the beginning of the user list, pass in an
	// empty string as startingAt
	ListUsers(startingAt string, maxUsers int) ([]UserInfo, error)

	// GetConfigProvider that was used to create the engine
	GetConfigProvider() ConfigurationProvider

	// GetDatabase that was used to create the engine
	GetDatabase() DatabaseConnection

	// GetUserSessionCache that was used to create the engine
	GetUserSessionCache() UserSessionCache
}
