package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

// Use the GIN context to create a link back to the service that we're running
func createHATEOASURL(c *gin.Context, format string, args ...any) string {

	// Determine HTTP or HTTPS
	scheme := ""
	if c.Request.TLS != nil {
		scheme = "s"
	}

	// Now we can build the root from the request host and scheme
	root := fmt.Sprintf("http%s://%s", scheme, c.Request.Host)
	return fmt.Sprintf("%s%s", root, fmt.Sprintf(format, args...))
}
