package jutzo

import (
	"time"
)

type UserInfo interface {

	// GetUsername associated with this user
	GetUsername() string

	// GetEmail associated with this user
	GetEmail() string

	// GetPasswordHash associated with this user
	GetPasswordHash() []byte

	// IsEmailValidated for this user
	IsEmailValidated() bool

	// HasRights will return true if the user contains ALL the rights requested
	HasRights(requested []string) bool

	// HasAnyRight will return true if the user has ANY of the rights requested
	HasAnyRight(requested []string) bool

	// GetAllRights that the user has been granted as a string array
	GetAllRights() []string

	// GrantRight to the user. A no-op if the right is already in place
	GrantRight(right string)

	// GetCreationTime for the user
	GetCreationTime() time.Time
}
