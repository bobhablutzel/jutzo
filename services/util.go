package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
)

// Redacts a password from a standard service connection URL
func redactPassword(connectionString string) (string, error) {
	if parsed, err := url.Parse(connectionString); err == nil {
		return fmt.Sprintf("%s://%s:*****@%s:%s)",
			parsed.Scheme, parsed.User.Username(), parsed.Host, parsed.Port()), nil
	} else {
		return connectionString, err
	}
}

// Make sure that the payload provided by the user matches the expected JSON payload
func checkValidPayload(c *gin.Context, err error) bool {
	if err != nil {
		errorMessage := fmt.Sprintf("Malformed request: %s", err.Error())
		c.String(http.StatusBadRequest, errorMessage)
		return false
	} else {
		return true
	}
}

// Generic method to return an internal server message if there is an error
func checkValidDBRequest(c *gin.Context, err error) bool {
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to access database: %s", err.Error())
		c.String(http.StatusInternalServerError, errorMessage)
		return false
	} else {
		return true
	}
}

// Use the GIN context to create a link back to the service that we're running
func createHATEOASURL(c *gin.Context, format string, args ...any) string {

	// Determine HTTP or HTTPS
	scheme := ""
	if c.Request.TLS != nil {
		scheme = "s"
	}

	// Now we can build the root from the request host and scheme
	root := fmt.Sprintf("http%s://%s", scheme, c.Request.Host)
	//root := getConfigString("JUTZO_HATEOAS_ROOT", fmt.Sprintf("localhost:%s", getConfigItemPort()))

	log.Printf("Using %s as the HATEOAS root url", root)
	return fmt.Sprintf("%s%s", root, fmt.Sprintf(format, args...))
}
