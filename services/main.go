package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Routine to set up the GIN router. This
// both creates the router, configures it, and
// defines the endpoints that we support.
func setupRouter(db *sql.DB) *gin.Engine {

	// Create the router with default parameters. Note
	// this depends on the environment variable GIN_MODE
	// to be set to "release" in order to create a production
	// GIN server
	router := gin.Default()

	// If we are in production, then use
	// the trusted cloudflare platform
	if gin.Mode() == gin.ReleaseMode {
		router.TrustedPlatform = gin.PlatformCloudflare
	}

	// Set up CORS middleware options.
	router.Use(cors.Default())

	// Set up a group so that all endpoints are within the /v1
	// namespace. This give us flexibility for breaking changes
	// in the future if needed.
	v1 := router.Group("/v1")

	// Ping test
	v1.GET("/ping", ping)

	// User management
	v1.POST("/user/register", func(c *gin.Context) { handleRegisterUser(c, db) })
	v1.POST("/user/login", func(c *gin.Context) { handleLogin(c, db) })
	v1.GET("/user/validateEmail/:key", func(c *gin.Context) { handleValidateEmail(c, db) })

	// Define a group for endpoints that require authentication but no specific rights
	authenticated := router.Group("/v1", requireValidJWTToken)
	authenticated.GET("/user/getValidationLink",
		func(c *gin.Context) { handleResendValidateEmailLink(c, db) })
	authenticated.GET("/user/logoff", handleLogoff)

	// Define a group for endpoints that require specific rights to access
	granted := router.Group("/v1", requireGrants([]string{"admin"}))
	granted.GET("/user/list", func(c *gin.Context) { handleListUsers(c, db) })

	// Blog methods
	v1.GET("/blog/newest", newest)
	v1.GET("/blog/entry/:id", blogEntry)

	return router
}

// Routine to handle a ping test
func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func handleRegisterUser(c *gin.Context, db *sql.DB) {

	// registerUserPayload is used with the registration POST request
	type registerUserPayload struct {
		User  string `json:"user" binding:"required"`
		Email string `json:"email" binding:"required"`
		Pass  string `json:"pass" binding:"required"`
	}

	// Get the JSON payload from the database
	var payload registerUserPayload
	err := c.BindJSON(&payload)
	if checkValidPayload(c, err) {
		if status, err := registerUser(db, payload.User, payload.Pass, payload.Email); err == nil {

			// Check the status
			switch status {
			case StatusInserted:
				{
					// Successfully inserted, create the email verification url
					if uuid, err := createUniqueValidationForUser(db, payload.User); err == nil {
						c.String(http.StatusOK, createHATEOASURL(c, "/v1/user/validateEmail/%s", uuid))
					} else {
						c.String(http.StatusInternalServerError, "User inserted; error creating validation uuid")
					}
				}
			case StatusDuplicateUsername:
				{
					c.String(http.StatusConflict, "Username is already in use")
				}
			case StatusDuplicateEmail:
				{
					c.String(http.StatusConflict, "Email is already in use")
				}
			}
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
	}
}

// Routine to get the validation URL for an existing account. This
// can be requested for a user that lost or didn't receive the validation URL
//
// This routine should be in a route protected by the requireValidJWTToken middleware
// in order to ensure that the user is logged in
func handleResendValidateEmailLink(c *gin.Context, db *sql.DB) {
	if user, ok := getUserFromContext(c); ok {
		if uuid, err := createUniqueValidationForUser(db, user); err == nil {
			c.String(http.StatusOK, createHATEOASURL(c, "/v1/user/validateEmail/%s ", uuid))
		} else {
			c.String(http.StatusInternalServerError, "Error getting validation uuid: %s", err.Error())
		}
	} else {
		c.String(http.StatusInternalServerError, "Should have user context from JTX middleware")
	}
}

// Process a request to validate an email
func handleValidateEmail(c *gin.Context, db *sql.DB) {
	uniqueID := c.Param("key")
	if validUUID, err := validateEmail(db, uniqueID); err == nil {
		if !validUUID {
			c.String(http.StatusBadRequest, "Invalid or expired validation request")
		} else {
			// TODO change to redirect to
			c.String(http.StatusOK, "OK")
		}
	} else {
		c.String(http.StatusInternalServerError, "Error processing validation: %s", err.Error())
	}
}

// Routine to log the user off, destroying the session token
func handleLogoff(c *gin.Context) {

	// Attempt to logoff the user
	if userSession, ok := getUserSessionFromContext(c); ok {
		if err := destroyUserSession(userSession.ID); err == nil {
			c.String(http.StatusOK, "OK")
		} else {
			c.String(http.StatusInternalServerError, "Error logging off")
		}
	} else {
		// User isn't authorized, and wants to logoff. Is this an
		// error or not? We decide it is
		c.String(http.StatusUnauthorized, "Invalid session")
	}
}

// Routine to validate that the password provided in clear text matches
// the hashed password in the database
func handleLogin(c *gin.Context, db *sql.DB) {

	// validatePasswordPayload is used with the HATEOAS link back to
	// validate that the user provided a valid email address
	type validatePasswordPayload struct {
		User string `json:"user" binding:"required"`
		Pass string `json:"pass" binding:"required"`
	}

	// Get the JSON payload from the database
	var payload validatePasswordPayload
	if err := c.BindJSON(&payload); err == nil {

		if token, err := login(db, payload.User, payload.Pass); err == nil {
			c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
			c.String(http.StatusOK, "OK")
		} else {
			c.String(http.StatusUnauthorized, "Invalid username or password")
		}

	} else {
		c.String(http.StatusBadRequest, "Malformed request: %s", err.Error())
	}
}

func handleListUsers(c *gin.Context, db *sql.DB) {
	if users, err := listUsers(db); err == nil {
		c.JSON(http.StatusOK, users)
	} else {
		c.String(http.StatusBadRequest, "Unable to list users: %s", err.Error())
	}
}

// Entry point for the server.
func main() {

	// See if we have to create the database
	initDB := slices.Contains(os.Args[1:], "--initDB")
	force := slices.Contains(os.Args[1:], "--force")

	// Open up our database connection
	if db, err := connectToDatabase(initDB, force); err == nil {

		// Close the database when we exit. (See below for
		// graceful exit on signals, e.g. SIGINT)
		defer shutdownDatabase(db)

		// Ensure that we have at least an administrative user
		if valid, err := ensureAdministrativeUser(db); err == nil {
			if valid {
				runServer(db)
			} else {
				log.Printf("Exiting because no administrative user has been defined")
			}
		} else {
			log.Printf("Exiting because the administrative user count couldn't be determined: %s", err.Error())
		}
	} else {
		log.Printf("Exiting because database could not be opened: %s", err.Error())
	}
}

func runServer(db *sql.DB) {
	// Initialize the JWT token module
	if err := initToken(); err == nil {

		// Now we can configure our router
		router := setupRouter(db)

		// Launch a new http server (using the port from the
		// PORT environment variable)
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", getConfigItemPort()),
			Handler: router,
		}

		// Handle incoming traffic
		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutdown Server ...")

		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown: ", err)
		}
	} else {
		log.Printf("Token system could not be initialized: %s", err.Error())
	}
	log.Println("Server exiting")

}
