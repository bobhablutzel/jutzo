// Defines a set of routines for managing the user session in a GIN
// server. This module has the functionality around creating the
// JWT token, decoding a JWT token, and managing the GIN middleware
// around rights based on identity

package jutzo

import (
	"time"
)

// UserSession defines the information that we know about a logged in user
type UserSession interface {

	// GetId for the session itself
	GetId() string

	// GetUserInfo gets the user information from the session
	GetUserInfo() UserInfo

	// GetDuration of the session from when it was initialized
	GetDuration() time.Duration
}

type UserSessionCache interface {

	// Connect to the cache so that user sessions can be managed
	Connect() error

	// Disconnect from the cache when we're done with it
	Disconnect() error

	// GetUserSessionByID will look for the ID in the cache and return it
	// if it exists and isn't expired. It will return nil if the session
	// cannot be found, and an error if there is a problem communicating
	// with the cache
	GetUserSessionByID(uniqueID string) (UserSession, error)

	// CacheUserSession so that it can be retrieved again by the unique ID
	CacheUserSession(userInfo UserInfo) (UserSession, error)

	// InvalidateUserSession by removing the session from the cache
	InvalidateUserSession(uniqueID string) error
}
