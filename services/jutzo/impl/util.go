package impl

import (
	"fmt"
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
