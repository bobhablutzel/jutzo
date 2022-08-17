package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"net/url"
	"services/jutzo"
	"time"
)

type RedisCache struct {
	config jutzo.ConfigurationProvider
	client *redis.Client
}

// NewRedisCache will create a new Redis cache for the engine to use
func NewRedisCache(configurationProvider jutzo.ConfigurationProvider) (jutzo.UserSessionCache, error) {
	result := new(RedisCache)
	result.config = configurationProvider
	err := result.Connect()
	return result, err
}

// Connect to the cache so that user sessions can be managed
func (cache *RedisCache) Connect() error {
	if cache.client == nil {

		// Connect to the redis service. We are given the connect information
		// as a redis:// connection URL, but the Redis middleware we're using
		// doesn't accept those URLs. So we parse it here and put it into the
		// format we can use. (The reason we accept in the URL format is that
		// it is the format natively provided by the hosting provider we initially
		// targeted, and because it allows for a single environment variable
		// rather than multiple ones.)
		//
		// We default to the default port on the local machine.
		var connectURL string
		var present bool
		if connectURL, present = cache.config.GetConfigurationString("JUTZO_REDIS_URL"); !present {
			connectURL = "redis://localhost"
		}

		// Use that to get the connection to the Redis server
		if parsed, err := url.Parse(connectURL); err == nil {

			// Sanity check that we got a Redis URL
			if parsed.Scheme == "redis" {

				// Get the username and host from the parsed information, build
				// the connection information, and connect
				password, _ := parsed.User.Password()
				cache.client = redis.NewClient(&redis.Options{
					Addr:     parsed.Host,
					Password: password,
					DB:       0,
				})
				return nil
			} else {
				return errors.New("invalid scheme for Redis configuration variable JUTZO_REDIS_URL")
			}
		} else {
			return err
		}
	} else {
		return nil
	}
}

// Disconnect from the cache when we're done with it
func (cache *RedisCache) Disconnect() error {
	cache.client = nil
	return nil
}

// GetUserSessionByID will look for the ID in the cache and return it
// if it exists and isn't expired. It will return nil if the session
// cannot be found, and an error if there is a problem communicating
// with the cache
func (cache *RedisCache) GetUserSessionByID(uniqueID string) (jutzo.UserSession, error) {

	// Get the actual user session from the Redis server
	if marshalledUserSession, err := cache.client.Get(context.Background(), uniqueID).Result(); err == nil {

		// Decode the user session. Note we need to have an allocated UserSession so that
		// it doesn't go out of scope when this function ends
		userSession := new(UserSessionImpl)
		if err = json.Unmarshal([]byte(marshalledUserSession), userSession); err == nil {
			return userSession, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

}

// CacheUserSession so that it can be retrieved again by the unique ID
func (cache *RedisCache) CacheUserSession(userInfo jutzo.UserInfo) (jutzo.UserSession, error) {
	// Create a new uuid to store that user session in the Redis cache with
	if uniqueID, err := uuid.NewRandom(); err == nil {

		// Get the new UUID as a string
		id := uniqueID.String()

		// Create a new user session
		userSession := new(UserSessionImpl)
		userSession.Info = userInfo.(*UserInfoImpl)
		userSession.ID = id
		userSession.Duration = time.Duration(8) * time.Hour

		if marshalledSession, err := json.Marshal(userSession); err == nil {
			return userSession, cache.client.Set(context.Background(), id, marshalledSession, userSession.Duration).Err()
		} else {
			return nil, errors.New(fmt.Sprintf("could not marshal the userSession: %s", err.Error()))
		}
	} else {
		return nil, errors.New("could not create unique Redis key")
	}

}

// InvalidateUserSession by removing the session from the cache
func (cache *RedisCache) InvalidateUserSession(uniqueID string) error {
	return cache.client.Del(context.Background(), uniqueID).Err()
}
