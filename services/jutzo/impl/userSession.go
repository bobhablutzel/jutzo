package impl

import (
	"services/jutzo"
	"time"
)

// UserSessionImpl defines the structure of the user session information that we store in Redis. This
// contains a unique ID, the username, and the rights. In a more secure environment we could
// also include device or IP fingerprint information - we'll make that a TODO
type UserSessionImpl struct {
	ID       string        `json:"ID"`
	Info     *UserInfoImpl `json:"info"`
	Duration time.Duration `json:"duration"`
}

// GetId for the session itself
func (userSession *UserSessionImpl) GetId() string {
	return userSession.ID
}

// GetUserInfo gets the user information from the session
func (userSession *UserSessionImpl) GetUserInfo() jutzo.UserInfo {
	return userSession.Info
}

// GetDuration of the session from when it was initialized
func (userSession *UserSessionImpl) GetDuration() time.Duration {
	return userSession.Duration
}
