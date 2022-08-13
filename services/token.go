// Defines a set of routines for managing the user session in a GIN
// server. This module has the functionality around creating the
// JWT token, decoding a JWT token, and managing the GIN middleware
// around rights based on identity

package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt"
)

// The JWT secret, loaded from the JUTZO_JWT_SECRET environment variable.
var jwtSecret []byte

// The redis connection
var redisClient *redis.Client

// Get the JWT secret from the environment variables
func initToken() error {
	rawString := getRequiredConfigString("JUTZO_JWT_SECRET")
	var err error
	if jwtSecret, err = hex.DecodeString(rawString); err == nil {

		// Connect to the redis service
		redisConnectionString := getRequiredConfigString("JUTZO_REDIS_URL")
		if parsed, err := url.Parse(redisConnectionString); err == nil {
			if parsed.Scheme == "redis" {
				port := parsed.Port()
				if port == "" {
					port = "6379"
				}
				password, _ := parsed.User.Password()

				redisClient = redis.NewClient(&redis.Options{
					Addr:     fmt.Sprintf("%s:%s", parsed.Host, port),
					Password: password,
					DB:       0,
				})
			} else {
				return errors.New("invalid scheme for redis connection environment variable JUTZO_REDIS_URL")
			}
		} else {
			return err
		}
	}
	return err
}

type UserSession struct {
	User   string      `json:"user"`
	Rights *UserRights `json:"rights"`
}

// The UserSessionCookieToken structure is used to contain the information
// about the user session in the JWT token
type UserSessionCookieToken struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

// Creates a new JWT token for the user with the provided user rights
func createToken(user string, rights *UserRights) (string, error) {

	// Create a new UserSession element with the rights provided
	claims := UserSession{Rights: rights}

	// Create a new uuid to store that user session in the Redis cache with
	if uniqueID, err := uuid.NewRandom(); err == nil {

		// Store the claims
		duration := time.Duration(8) * time.Hour
		id := uniqueID.String()

		if marshalledSession, err := json.Marshal(claims); err == nil {
			if err = redisClient.Set(context.Background(), id, marshalledSession, duration).Err(); err == nil {
				cookie := UserSessionCookieToken{
					ID: id,
					StandardClaims: jwt.StandardClaims{
						Issuer:    "Jutzo Service",
						Subject:   user,
						ExpiresAt: time.Now().Local().Add(duration).Unix(),
					},
				}

				// Create the token
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, cookie)

				// Sign and get the complete encoded token as a string using the secret
				return token.SignedString(jwtSecret)
			} else {
				return "", errors.New("unable to store session to Redis server")
			}
		} else {
			return "", errors.New(fmt.Sprintf("could not marshal the claims: %s", err.Error()))
		}
	} else {
		return "", errors.New("could not create unique Redis key")
	}
}

// Decodes a JWT token to get the user and rights
func decodeToken(encoded string) (*UserSessionCookieToken, error) {

	// Attempt to parse the token
	token, err := jwt.ParseWithClaims(encoded, &UserSessionCookieToken{}, func(token *jwt.Token) (interface{}, error) {

		// Validate that we got the right signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	} else {

		claims, ok := token.Claims.(*UserSessionCookieToken)
		if ok && token.Valid {
			return claims, nil
		} else {
			return nil, errors.New("unable to extract JWT claims")
		}
	}
}

// Find the bearer token and parse it to the context
//
// This is a helper function for a GIN middleware
// function. See requireValidJWTToken or requireGrant for the
// middleware functions
func parseJWTToken(c *gin.Context) {

	// Get the authorization header
	token := c.Request.Header["Authorization"]
	if token == nil || 1 != len(token) {
		c.AbortWithStatus(http.StatusUnauthorized)
	} else {

		// Decode the authorization header to get the JWT encoded token
		var encoded string
		if _, err := fmt.Sscanf(token[0], "Bearer %s", &encoded); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			// Decode the JWT token to get the user session
			if sessionCookieToken, err := decodeToken(encoded); err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {

				// Get the actual user session from the Redis server
				if marshalledUserSession, err := redisClient.Get(context.Background(), sessionCookieToken.ID).Result(); err == nil {

					// Decode the user session. Note we need to have an allocated UserSession so that
					// it doesn't go out of scope when this function ends
					userSession := new(UserSession)
					if err = json.Unmarshal([]byte(marshalledUserSession), userSession); err == nil {
						// Save the JWT token for downstream usage

						ctx := context.WithValue(c.Request.Context(), "userSession", userSession)
						c.Request = c.Request.WithContext(ctx)
					} else {
						c.AbortWithStatus(http.StatusUnauthorized)
					}
				} else {
					c.AbortWithStatus(http.StatusUnauthorized)
				}
			}
		}
	}
}

// Attempt to get the user session. This routine depends on the
// context having been set up in one of the GIN middleware methods
// (either requireValidJWTToken or requireGrants), so this should only be
// called in a routine that is part of a group that has this middleware
// invoked before the endpoint.
func getUserSessionFromContext(c *gin.Context) (*UserSession, bool) {
	if userSessionObj := c.Request.Context().Value("userSession"); userSessionObj != nil {
		userSession, ok := userSessionObj.(*UserSession)
		return userSession, ok
	} else {
		return nil, false
	}
}

func getUserFromContext(c *gin.Context) (string, bool) {
	if userSession, ok := getUserSessionFromContext(c); ok {
		return userSession.User, true
	} else {
		return "", false
	}
}

func getRights(c *gin.Context) (*UserRights, bool) {
	if userSession, ok := getUserSessionFromContext(c); ok {
		return userSession.Rights, true
	} else {
		return nil, false
	}
}

func requireGrants(grants []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		parseJWTToken(c)
		if !c.IsAborted() {
			if userRights, ok := getRights(c); ok {
				for _, grant := range grants {
					if !slices.Contains(userRights.Grants, grant) {
						c.AbortWithStatus(http.StatusForbidden)
					}
				}
				c.Next()
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
		}
	}
}

func requireValidJWTToken(c *gin.Context) {
	parseJWTToken(c)
	if !c.IsAborted() {
		c.Next()
	}
}
