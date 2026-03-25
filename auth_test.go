package scalar

import (
	"encoding/json"
	"testing"
)

func TestAPIKeyAuth(t *testing.T) {
	tests := []struct {
		name        string
		keyName     string
		location    APIKeyLocation
		description string
	}{
		{
			name:        "API Key in Header",
			keyName:     "X-API-Key",
			location:    APIKeyInHeader,
			description: "API key for authentication",
		},
		{
			name:     "API Key in Query",
			keyName:  "api_key",
			location: APIKeyInQuery,
		},
		{
			name:     "API Key in Cookie",
			keyName:  "session",
			location: APIKeyInCookie,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := APIKeyAuth(tt.keyName, tt.location)
			if tt.description != "" {
				auth.WithDescription(tt.description)
			}

			if auth.Type != AuthTypeAPIKey {
				t.Errorf("Expected type %s, got %s", AuthTypeAPIKey, auth.Type)
			}
			if auth.Name != tt.keyName {
				t.Errorf("Expected name %s, got %s", tt.keyName, auth.Name)
			}
			if auth.In != tt.location {
				t.Errorf("Expected location %s, got %s", tt.location, auth.In)
			}

			// Test JSON conversion
			jsonStr, err := auth.ToJSON()
			if err != nil {
				t.Errorf("Failed to convert to JSON: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
				t.Errorf("Failed to parse JSON: %v", err)
			}
		})
	}
}

func TestBearerAuth(t *testing.T) {
	auth := BearerAuth().
		WithBearerFormat("JWT").
		WithDescription("JWT Bearer token authentication")

	if auth.Type != AuthTypeHTTP {
		t.Errorf("Expected type %s, got %s", AuthTypeHTTP, auth.Type)
	}
	if auth.Scheme != HTTPSchemeBearer {
		t.Errorf("Expected scheme %s, got %s", HTTPSchemeBearer, auth.Scheme)
	}
	if auth.BearerFormat != "JWT" {
		t.Errorf("Expected bearer format JWT, got %s", auth.BearerFormat)
	}

	jsonStr, err := auth.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
}

func TestBasicAuth(t *testing.T) {
	auth := BasicAuth().WithDescription("Basic HTTP authentication")

	if auth.Type != AuthTypeHTTP {
		t.Errorf("Expected type %s, got %s", AuthTypeHTTP, auth.Type)
	}
	if auth.Scheme != HTTPSchemeBasic {
		t.Errorf("Expected scheme %s, got %s", HTTPSchemeBasic, auth.Scheme)
	}

	jsonStr, err := auth.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
}

func TestOAuth2Auth(t *testing.T) {
	scopes := map[string]string{
		"read":  "Read access",
		"write": "Write access",
		"admin": "Admin access",
	}

	auth := OAuth2Auth().
		WithAuthorizationCode(
			"https://example.com/oauth/authorize",
			"https://example.com/oauth/token",
			scopes,
		).
		WithClientCredentials(
			"https://example.com/oauth/token",
			scopes,
		).
		WithDescription("OAuth 2.0 authentication")

	if auth.Type != AuthTypeOAuth2 {
		t.Errorf("Expected type %s, got %s", AuthTypeOAuth2, auth.Type)
	}

	if len(auth.Flows) != 2 {
		t.Errorf("Expected 2 flows, got %d", len(auth.Flows))
	}

	jsonStr, err := auth.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
}

func TestOAuth2WithRefreshURL(t *testing.T) {
	scopes := map[string]string{
		"read": "Read access",
	}

	auth := OAuth2Auth().
		WithAuthorizationCode(
			"https://example.com/oauth/authorize",
			"https://example.com/oauth/token",
			scopes,
		).
		WithRefreshURL(OAuth2FlowAuthorizationCode, "https://example.com/oauth/refresh")

	flow := auth.Flows[string(OAuth2FlowAuthorizationCode)]
	if flow.RefreshURL != "https://example.com/oauth/refresh" {
		t.Errorf("Expected refresh URL to be set, got %s", flow.RefreshURL)
	}
}

func TestOpenIDConnectAuth(t *testing.T) {
	url := "https://example.com/.well-known/openid-configuration"
	auth := OpenIDConnectAuth(url).
		WithDescription("OpenID Connect authentication")

	if auth.Type != AuthTypeOpenID {
		t.Errorf("Expected type %s, got %s", AuthTypeOpenID, auth.Type)
	}
	if auth.OpenIDConnectURL != url {
		t.Errorf("Expected URL %s, got %s", url, auth.OpenIDConnectURL)
	}

	jsonStr, err := auth.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
}

func TestMultipleAuth(t *testing.T) {
	configs := map[string]*AuthConfig{
		"apiKey": APIKeyAuth("X-API-Key", APIKeyInHeader).
			WithDescription("API Key authentication"),
		"bearer": BearerAuth().
			WithBearerFormat("JWT").
			WithDescription("JWT Bearer token"),
		"oauth2": OAuth2Auth().
			WithAuthorizationCode(
				"https://example.com/oauth/authorize",
				"https://example.com/oauth/token",
				map[string]string{
					"read":  "Read access",
					"write": "Write access",
				},
			),
	}

	jsonStr, err := MultipleAuth(configs)
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 auth methods, got %d", len(result))
	}
}

func TestMustJSON(t *testing.T) {
	auth := BearerAuth().WithBearerFormat("JWT")

	// Should not panic
	jsonStr := auth.MustJSON()
	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
}

func TestAuthConfigChaining(t *testing.T) {
	// Test method chaining
	auth := BearerAuth().
		WithBearerFormat("JWT").
		WithDescription("Test authentication")

	if auth.BearerFormat != "JWT" {
		t.Errorf("Expected bearer format JWT, got %s", auth.BearerFormat)
	}
	if auth.Description != "Test authentication" {
		t.Errorf("Expected description to be set, got %s", auth.Description)
	}
}
