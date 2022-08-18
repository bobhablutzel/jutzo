package impl

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"services/jutzo"
)

// EngineImpl provides the implementation structure for the
// implementation of a Jutzo engine
type EngineImpl struct {
	config jutzo.ConfigurationProvider
	db     jutzo.DatabaseConnection
	cache  jutzo.UserSessionCache
}

// NewJutzoEngine sets up the Jutzo environment with the configuration information provided.
// The routine expects that the database and cache connections are already established
// but will connect if not
func NewJutzoEngine(config jutzo.ConfigurationProvider, connection jutzo.DatabaseConnection, cache jutzo.UserSessionCache) (jutzo.Engine, error) {

	// Create a new Jutzo engine we can pass back
	engine := new(EngineImpl)
	engine.config = config
	engine.db = connection
	engine.cache = cache

	// Connect to the database. Note this is should be a no-op if already connected.
	if err := connection.Connect(); err != nil {
		return nil, err
	}

	// Initialize the user session cache. Again, a no-op if already connected
	if err := cache.Connect(); err != nil {
		return nil, err
	}

	// Ensure that there is an administrative user registered
	if count, err := connection.GetAdminCount(); err != nil {
		return nil, err
	} else {
		if count == 0 {
			user, userPresent := config.GetConfigurationString("JUTZO_ADMIN_USER")
			password, passwordPresent := config.GetConfigurationString("JUTZO_ADMIN_PASS")
			email, emailPresent := config.GetConfigurationString("JUTZO_ADMIN_EMAIL")
			if userPresent && passwordPresent && emailPresent {

				log.Printf("Creating administrator %s @ %s", user, email)
				if registerUser, userInfo, err := engine.RegisterUser(user, password, email); err != nil {
					return nil, err
				} else {
					if registerUser != jutzo.Success {
						return nil, errors.New("unable to find or register admin account due to conflict")
					} else {

						// Force update the rights for that user. This will also allow them
						// to log in with an email that hasn't been validated.
						userInfo.GrantRight("admin")
						if err := connection.UpdateUserInfo(userInfo); err != nil {

							log.Printf("Update of grants failed")
							return nil, err
						}

					}
				}
			} else {
				return nil, errors.New("no administrator found, and missing configuration to create one")
			}
		}
	}

	return engine, nil
}

func (engine *EngineImpl) GetConfigProvider() jutzo.ConfigurationProvider {
	return engine.config
}

func (engine *EngineImpl) GetDatabase() jutzo.DatabaseConnection {
	return engine.db
}

func (engine *EngineImpl) GetUserSessionCache() jutzo.UserSessionCache {
	return engine.cache
}

// Shutdown closes down the Jutzo engine gracefully
func (engine *EngineImpl) Shutdown() {
	err := engine.db.Shutdown()
	if err != nil {
		log.Printf("Error shutting down db: %s", err.Error())
	}
}

func (engine *EngineImpl) Login(user string, password string) (jutzo.UserSession, error) {

	if userInfo, err := engine.db.RetrieveUserInformation(user); err != nil {
		return nil, err
	} else {

		// Make sure that the user is allowed to log in
		admin := userInfo.HasRights([]string{"admin"})
		canLogin := userInfo.HasRights([]string{"login"})
		if (userInfo.IsEmailValidated() && canLogin) || admin {

			// User is able to log in, compare the password hash. We do this second because it's
			// a more expensive operation than just the testing done above
			if err = bcrypt.CompareHashAndPassword(userInfo.GetPasswordHash(), []byte(password)); err == nil {

				// Successful login, create a user session
				return engine.cache.CacheUserSession(userInfo)
			} else {
				return nil, err
			}
		} else {
			return nil, errors.New(fmt.Sprintf("User %s is not active - unvalidated or no login rights", user))
		}
	}
}

func (engine *EngineImpl) RegisterUser(user string, password string, email string) (int, jutzo.UserInfo, error) {

	// See if the username or email already exists
	if usernameExists, emailExists, err := engine.db.CheckForUsernameOrEmail(user, email); err == nil {

		if usernameExists {
			return jutzo.DuplicateUsername, nil, nil
		} else if emailExists {
			return jutzo.DuplicateEmail, nil, nil
		} else {
			// Get the cost for creating the password. This cost can be
			// set as an environment variable or defaulted, and it will
			// be stored with the password for validation later
			passwordCost, isPresent := engine.config.GetConfigurationInt("JUTZO_HASH_COST")
			if !isPresent {
				passwordCost = 15
			}

			// Encrypt the password provided
			if passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost); err != nil {
				return 0, nil, err
			} else {

				// Store the user
				if userInfo, err := engine.db.StoreUser(user, email, passwordHash); err != nil {
					return 0, nil, err
				} else {
					return jutzo.Success, userInfo, nil
				}
			}
		}
	} else {
		return 0, nil, err
	}
}

func (engine *EngineImpl) CreateUniqueValidationForUser(user string) (string, string, error) {
	return engine.db.CreateValidationFor(user)
}

func (engine *EngineImpl) ValidateEmail(uniqueID string) error {
	return engine.db.CompleteValidationFor(uniqueID)
}

func (engine *EngineImpl) LoadUserSession(uniqueID string) (jutzo.UserSession, error) {
	return engine.cache.GetUserSessionByID(uniqueID)
}

func (engine *EngineImpl) DestroyUserSession(uniqueID string) error {
	return engine.cache.InvalidateUserSession(uniqueID)
}

// ListUsers can be called by an admin to list users. This call will
// return up to maxUsers users at a time; if you want to retrieve the
// remaining users call ListUsers again passing in the username from the
// last record returned. If no more users can be found, an empty slice
// will be returned with no errors.
// If you want to start at the beginning of the user list, pass in an
// empty string as startingAt
func (engine *EngineImpl) ListUsers(startingAt string, maxUsers int) ([]jutzo.UserInfo, error) {
	return engine.db.ListUsers(startingAt, maxUsers)
}
