package impl

import (
	"golang.org/x/exp/slices"
	"services/jutzo"
	"time"
)

type UserInfoImpl struct {
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PasswordHash   []byte    `json:"passwordHash"`
	PasswordCost   int       `json:"passwordCost"`
	EmailValidated bool      `json:"emailValidated"`
	Rights         []string  `json:"rights"`
	CreationTime   time.Time `json:"creationTime"`
}

func NewUserInfo(username string, email string, passwordHash []byte, passwordCost int, emailValidated bool, rights []string, creationTime time.Time) jutzo.UserInfo {
	result := new(UserInfoImpl)
	result.Username = username
	result.Email = email
	result.PasswordHash = passwordHash
	result.PasswordCost = passwordCost
	result.EmailValidated = emailValidated
	result.Rights = rights
	result.CreationTime = creationTime
	return result
}

// GetUsername associated with this user
func (userInfo *UserInfoImpl) GetUsername() string {
	return userInfo.Username
}

// GetEmail associated with this user
func (userInfo *UserInfoImpl) GetEmail() string {
	return userInfo.Email
}

// GetPasswordHash associated with this user
func (userInfo *UserInfoImpl) GetPasswordHash() []byte {
	return userInfo.PasswordHash
}

// GetPasswordCost associated with this user
func (userInfo *UserInfoImpl) GetPasswordCost() int {
	return userInfo.PasswordCost
}

// IsEmailValidated for this user
func (userInfo *UserInfoImpl) IsEmailValidated() bool {
	return userInfo.EmailValidated
}

// HasRights will return true if the user contains ALL the rights requested
func (userInfo *UserInfoImpl) HasRights(requested []string) bool {
	rightCount := 0
	for _, requestedRight := range requested {
		for _, item := range userInfo.Rights {
			if item == requestedRight {
				rightCount++
			}
		}
	}
	return rightCount == len(requested)
}

// HasAnyRight will return true if the user has ANY of the rights requested
func (userInfo *UserInfoImpl) HasAnyRight(requested []string) bool {
	for _, requestedRight := range requested {
		for _, item := range userInfo.Rights {
			if item == requestedRight {
				return true
			}
		}
	}
	return false
}

// GetAllRights that the user has been granted as a string array
func (userInfo *UserInfoImpl) GetAllRights() []string {
	return userInfo.Rights
}

// GrantRight to the user. A no-op if the right is already in place
func (userInfo *UserInfoImpl) GrantRight(right string) {
	if !slices.Contains(userInfo.Rights, right) {
		userInfo.Rights = append(userInfo.Rights, right)
	}
}

// GetCreationTime for the user
func (userInfo *UserInfoImpl) GetCreationTime() time.Time {
	return userInfo.CreationTime
}
