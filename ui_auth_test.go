package scalar

import (
	"encoding/base64"
	"testing"
)

func TestIsUIAuthEnabled(t *testing.T) {
	tests := []struct {
		name     string
		options  *Options
		expected bool
	}{
		{
			name: "Auth enabled with both username and password",
			options: &Options{
				UIUsername: "admin",
				UIPassword: "secret",
			},
			expected: true,
		},
		{
			name: "Auth disabled with empty username",
			options: &Options{
				UIUsername: "",
				UIPassword: "secret",
			},
			expected: false,
		},
		{
			name: "Auth disabled with empty password",
			options: &Options{
				UIUsername: "admin",
				UIPassword: "",
			},
			expected: false,
		},
		{
			name: "Auth disabled with both empty",
			options: &Options{
				UIUsername: "",
				UIPassword: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.IsUIAuthEnabled()
			if result != tt.expected {
				t.Errorf("IsUIAuthEnabled() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateUICredentials(t *testing.T) {
	tests := []struct {
		name     string
		options  *Options
		username string
		password string
		expected bool
	}{
		{
			name: "Valid credentials",
			options: &Options{
				UIUsername: "admin",
				UIPassword: "secret123",
			},
			username: "admin",
			password: "secret123",
			expected: true,
		},
		{
			name: "Invalid username",
			options: &Options{
				UIUsername: "admin",
				UIPassword: "secret123",
			},
			username: "wrong",
			password: "secret123",
			expected: false,
		},
		{
			name: "Invalid password",
			options: &Options{
				UIUsername: "admin",
				UIPassword: "secret123",
			},
			username: "admin",
			password: "wrong",
			expected: false,
		},
		{
			name: "Auth not configured - should allow access",
			options: &Options{
				UIUsername: "",
				UIPassword: "",
			},
			username: "any",
			password: "any",
			expected: true,
		},
		{
			name: "Empty credentials against configured auth",
			options: &Options{
				UIUsername: "admin",
				UIPassword: "secret123",
			},
			username: "",
			password: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.ValidateUICredentials(tt.username, tt.password)
			if result != tt.expected {
				t.Errorf("ValidateUICredentials() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseBasicAuth(t *testing.T) {
	tests := []struct {
		name             string
		authHeader       string
		expectedUsername string
		expectedPassword string
		expectedOk       bool
	}{
		{
			name:             "Valid Basic Auth",
			authHeader:       "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:secret123")),
			expectedUsername: "admin",
			expectedPassword: "secret123",
			expectedOk:       true,
		},
		{
			name:             "Valid Basic Auth with colon in password",
			authHeader:       "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass:word")),
			expectedUsername: "user",
			expectedPassword: "pass:word",
			expectedOk:       true,
		},
		{
			name:             "Invalid - no Basic prefix",
			authHeader:       base64.StdEncoding.EncodeToString([]byte("admin:secret123")),
			expectedUsername: "",
			expectedPassword: "",
			expectedOk:       false,
		},
		{
			name:             "Invalid - wrong scheme",
			authHeader:       "Bearer token123",
			expectedUsername: "",
			expectedPassword: "",
			expectedOk:       false,
		},
		{
			name:             "Invalid - malformed base64",
			authHeader:       "Basic invalid!!!",
			expectedUsername: "",
			expectedPassword: "",
			expectedOk:       false,
		},
		{
			name:             "Invalid - no colon separator",
			authHeader:       "Basic " + base64.StdEncoding.EncodeToString([]byte("adminonly")),
			expectedUsername: "",
			expectedPassword: "",
			expectedOk:       false,
		},
		{
			name:             "Empty auth header",
			authHeader:       "",
			expectedUsername: "",
			expectedPassword: "",
			expectedOk:       false,
		},
		{
			name:             "Case insensitive Basic prefix",
			authHeader:       "basic " + base64.StdEncoding.EncodeToString([]byte("admin:secret")),
			expectedUsername: "admin",
			expectedPassword: "secret",
			expectedOk:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, password, ok := ParseBasicAuth(tt.authHeader)
			if ok != tt.expectedOk {
				t.Errorf("ParseBasicAuth() ok = %v, want %v", ok, tt.expectedOk)
			}
			if username != tt.expectedUsername {
				t.Errorf("ParseBasicAuth() username = %v, want %v", username, tt.expectedUsername)
			}
			if password != tt.expectedPassword {
				t.Errorf("ParseBasicAuth() password = %v, want %v", password, tt.expectedPassword)
			}
		})
	}
}
