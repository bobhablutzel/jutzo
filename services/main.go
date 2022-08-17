package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"services/jutzo"
	"services/jutzo/impl"
	"strconv"
	"syscall"
	"time"
)

// Routine to set up the GIN router. This
// both creates the router, configures it, and
// defines the endpoints that we support.
func setupRouter(engine jutzo.Engine, configuration Configuration) (*gin.Engine, error) {

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

	// Get our JWT token engine
	if tokenEngine, err := NewTokenEngine(configuration); err == nil {

		// Set up a group so that all endpoints are within the /v1
		// namespace. This give us flexibility for breaking changes
		// in the future if needed.
		v1 := router.Group("/v1")

		// Ping test
		v1.GET("/ping", ping)

		// User management
		v1.POST("/user/register", func(c *gin.Context) { handleRegisterUser(c, engine) })
		v1.POST("/user/login", func(c *gin.Context) { handleLogin(c, tokenEngine, engine) })
		v1.GET("/user/validateEmail/:key", func(c *gin.Context) { handleValidateEmail(c, engine) })

		// Define a group for endpoints that require authentication but no specific rights
		authenticated := router.Group("/v1", requireValidJWTToken(engine, tokenEngine))
		authenticated.GET("/user/getValidationLink",
			func(c *gin.Context) { handleResendValidateEmailLink(c, engine) })
		authenticated.GET("/user/logoff", func(c *gin.Context) { handleLogoff(c, engine) })

		// Define a group for endpoints that require specific rights to access
		granted := router.Group("/v1", requireGrants(engine, tokenEngine, []string{"admin"}))
		granted.GET("/user/list", func(c *gin.Context) { handleListUsers(c, engine) })

		// Blog methods
		v1.GET("/blog/newest", newest)
		v1.GET("/blog/entry/:id", blogEntry)

		return router, nil
	} else {
		return nil, err
	}
}

// Find the bearer token and parse it to the context
//
// This is a helper function for a GIN middleware
// function. See requireValidJWTToken or requireGrant for the
// middleware functions
func parseJWTToken(engine jutzo.Engine, tokenEngine TokenEngine, c *gin.Context) {

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
			if uniqueID, err := tokenEngine.Decode(encoded); err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {

				// Load the user session from the engine
				if userSession, err := engine.LoadUserSession(uniqueID); err == nil {

					// Save the user session for downstream usage. We do this by creating a new
					// request with an updated context that contains the user session as a value
					ctx := context.WithValue(c.Request.Context(), "userSession", userSession)
					c.Request = c.Request.WithContext(ctx)
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
func getUserSessionFromContext(c *gin.Context) (jutzo.UserSession, bool) {
	if userSessionObj := c.Request.Context().Value("userSession"); userSessionObj != nil {
		userSession, ok := userSessionObj.(jutzo.UserSession)
		return userSession, ok
	} else {
		return nil, false
	}
}

func requireGrants(engine jutzo.Engine, tokenEngine TokenEngine, requiredRights []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		parseJWTToken(engine, tokenEngine, c)
		if !c.IsAborted() {
			if userSession, ok := getUserSessionFromContext(c); ok {
				if userSession.GetUserInfo().HasRights(requiredRights) {
					c.Next()
				} else {
					c.AbortWithStatus(http.StatusForbidden)
				}
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
		}
	}
}

func requireValidJWTToken(engine jutzo.Engine, tokenEngine TokenEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		parseJWTToken(engine, tokenEngine, c)
		if !c.IsAborted() {
			c.Next()
		}
	}
}

// Routine to handle a ping test
func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func handleRegisterUser(c *gin.Context, engine jutzo.Engine) {

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
		if status, _, err := engine.RegisterUser(payload.User, payload.Pass, payload.Email); err == nil {

			// Check the status
			switch status {
			case jutzo.Success:
				{
					// Successfully inserted, create the email verification url
					if uuid, _, err := engine.CreateUniqueValidationForUser(payload.User); err == nil {
						c.String(http.StatusOK, createHATEOASURL(c, "/v1/user/validateEmail/%s", uuid))
					} else {
						c.String(http.StatusInternalServerError, "User inserted; error creating validation uuid")
					}
				}
			case jutzo.DuplicateUsername:
				{
					c.String(http.StatusConflict, "Username is already in use")
				}
			case jutzo.DuplicateEmail:
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
func handleResendValidateEmailLink(c *gin.Context, engine jutzo.Engine) {
	if userSession, ok := getUserSessionFromContext(c); ok {
		if uuid, _, err := engine.CreateUniqueValidationForUser(userSession.GetUserInfo().GetUsername()); err == nil {
			c.String(http.StatusOK, createHATEOASURL(c, "/v1/userSession/validateEmail/%s ", uuid))
		} else {
			c.String(http.StatusInternalServerError, "Error getting validation uuid: %s", err.Error())
		}
	} else {
		c.String(http.StatusInternalServerError, "Should have userSession context from JTX middleware")
	}
}

// Process a request to validate an email
func handleValidateEmail(c *gin.Context, engine jutzo.Engine) {
	uniqueID := c.Param("key")
	if err := engine.ValidateEmail(uniqueID); err == nil {
		// TODO change to redirect to
		c.String(http.StatusOK, "OK")
	} else {
		c.String(http.StatusInternalServerError, "Error processing validation: %s", err.Error())
	}
}

// Routine to log the user off, destroying the session token
func handleLogoff(c *gin.Context, engine jutzo.Engine) {

	// Attempt to logoff the user
	if userSession, ok := getUserSessionFromContext(c); ok {
		if err := engine.DestroyUserSession(userSession.GetId()); err == nil {
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
func handleLogin(c *gin.Context, tokenEngine TokenEngine, engine jutzo.Engine) {

	// validatePasswordPayload is used with the HATEOAS link back to
	// validate that the user provided a valid email address
	type validatePasswordPayload struct {
		User string `json:"user" binding:"required"`
		Pass string `json:"pass" binding:"required"`
	}

	// Get the JSON payload from the database
	var payload validatePasswordPayload
	if err := c.BindJSON(&payload); err == nil {

		if userSession, err := engine.Login(payload.User, payload.Pass); err == nil {

			log.Printf("User password accepted, session id: %s", userSession.GetId())
			user, id, duration := userSession.GetUserInfo().GetUsername(), userSession.GetId(), userSession.GetDuration()
			if token, err := tokenEngine.Encode(user, id, duration); err == nil {
				c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
				c.String(http.StatusOK, "OK")
			}
		} else {
			c.String(http.StatusUnauthorized, "Invalid username or password")
		}

	} else {
		c.String(http.StatusBadRequest, "Malformed request: %s", err.Error())
	}
}

func handleListUsers(c *gin.Context, engine jutzo.Engine) {

	// Get the starting username (if provided)
	startingUsername := c.DefaultQuery("from", "")
	maxUsersString := c.DefaultQuery("count", "0")
	if maxUsers, err := strconv.Atoi(maxUsersString); err == nil {
		log.Printf("Query starting at %s for %d", startingUsername, maxUsers)

		if users, err := engine.ListUsers(startingUsername, maxUsers); err == nil {
			if users == nil {
				c.String(http.StatusNoContent, "")
			} else {
				c.JSON(http.StatusOK, users)
			}
		} else {
			c.String(http.StatusInternalServerError, "Unable to list users: %s", err.Error())
		}
	} else {
		c.String(http.StatusBadRequest, "Malformed count: %s", err.Error())
	}

}

func runServer(engine jutzo.Engine, configuration Configuration) {

	// Now we can configure our router
	if router, err := setupRouter(engine, configuration); err == nil {

		// Launch a new http server (using the port from the
		// PORT environment variable)
		port := os.Getenv("JUTZO_SERVER_PORT")
		if port == "" {
			port = "8080"
		}
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
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
		log.Fatal("Error configuring router: ", err)
	}
}

// Entry point for the server.
func main() {

	//// See if we have to create the database
	//initDB := slices.Contains(os.Args[1:], "--initDB")
	//force := slices.Contains(os.Args[1:], "--force")

	// Create a configuration provider
	configuration := Configuration{}

	db := impl.NewPostgresConnection(configuration)
	if err := db.Connect(); err == nil {

		if cache, err := impl.NewRedisCache(configuration); err == nil {

			// Initialize the jutzo engine
			if engine, err := impl.NewJutzoEngine(configuration, db, cache); err == nil {
				defer engine.Shutdown()

				// Run the engine
				runServer(engine, configuration)

			} else {
				log.Printf("Jutzo system could not be initialized: %s", err.Error())
			}
		} else {
			log.Printf("Unable to create cache")
		}
	} else {
		log.Printf("Error connecting to database: %s", err.Error())
	}
	log.Println("Server exiting")
}
