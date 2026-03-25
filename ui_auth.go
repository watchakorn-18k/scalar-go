package scalar

import (
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

// UIAuthConfig holds the UI authentication configuration
type UIAuthConfig struct {
	Username string
	Password string
}

// IsUIAuthEnabled checks if UI authentication is enabled
func (o *Options) IsUIAuthEnabled() bool {
	return o.UIUsername != "" && o.UIPassword != ""
}

// ValidateUICredentials validates the provided username and password against configured credentials
// Uses constant-time comparison to prevent timing attacks
func (o *Options) ValidateUICredentials(username, password string) bool {
	if !o.IsUIAuthEnabled() {
		return true // No auth configured, allow access
	}

	usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(o.UIUsername)) == 1
	passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(o.UIPassword)) == 1

	return usernameMatch && passwordMatch
}

// ParseBasicAuth parses an HTTP Basic Authentication string from Authorization header
// Returns username and password, or empty strings if parsing fails
func ParseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", "", false
	}

	encoded := auth[len(prefix):]
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", false
	}

	credentials := string(decoded)
	colonIndex := strings.Index(credentials, ":")
	if colonIndex == -1 {
		return "", "", false
	}

	return credentials[:colonIndex], credentials[colonIndex+1:], true
}
