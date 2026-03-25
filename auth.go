package scalar

import "encoding/json"

// AuthType represents the type of authentication
type AuthType string

const (
	AuthTypeAPIKey AuthType = "apiKey"
	AuthTypeHTTP   AuthType = "http"
	AuthTypeOAuth2 AuthType = "oauth2"
	AuthTypeOpenID AuthType = "openIdConnect"
)

// APIKeyLocation represents where the API key should be sent
type APIKeyLocation string

const (
	APIKeyInHeader APIKeyLocation = "header"
	APIKeyInQuery  APIKeyLocation = "query"
	APIKeyInCookie APIKeyLocation = "cookie"
)

// HTTPScheme represents HTTP authentication schemes
type HTTPScheme string

const (
	HTTPSchemeBasic  HTTPScheme = "basic"
	HTTPSchemeBearer HTTPScheme = "bearer"
	HTTPSchemeDigest HTTPScheme = "digest"
)

// OAuth2Flow represents OAuth 2.0 flow types
type OAuth2Flow string

const (
	OAuth2FlowAuthorizationCode OAuth2Flow = "authorizationCode"
	OAuth2FlowClientCredentials OAuth2Flow = "clientCredentials"
	OAuth2FlowImplicit          OAuth2Flow = "implicit"
	OAuth2FlowPassword          OAuth2Flow = "password"
)

// AuthConfig represents the base authentication configuration
type AuthConfig struct {
	Type        AuthType               `json:"type"`
	Name        string                 `json:"name,omitempty"`
	In          APIKeyLocation         `json:"in,omitempty"`
	Scheme      HTTPScheme             `json:"scheme,omitempty"`
	BearerFormat string                `json:"bearerFormat,omitempty"`
	Flows       map[string]OAuth2FlowConfig `json:"flows,omitempty"`
	OpenIDConnectURL string            `json:"openIdConnectUrl,omitempty"`
	Description string                 `json:"description,omitempty"`
}

// OAuth2FlowConfig represents OAuth 2.0 flow configuration
type OAuth2FlowConfig struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

// APIKeyAuth creates an API Key authentication configuration
//
// Example:
//
//	auth := scalar.APIKeyAuth("X-API-Key", scalar.APIKeyInHeader).
//		WithDescription("API key for authentication")
func APIKeyAuth(name string, location APIKeyLocation) *AuthConfig {
	return &AuthConfig{
		Type: AuthTypeAPIKey,
		Name: name,
		In:   location,
	}
}

// BearerAuth creates a Bearer token authentication configuration
//
// Example:
//
//	auth := scalar.BearerAuth().
//		WithBearerFormat("JWT").
//		WithDescription("JWT Bearer token")
func BearerAuth() *AuthConfig {
	return &AuthConfig{
		Type:   AuthTypeHTTP,
		Scheme: HTTPSchemeBearer,
	}
}

// BasicAuth creates a Basic HTTP authentication configuration
//
// Example:
//
//	auth := scalar.BasicAuth().
//		WithDescription("Basic HTTP authentication")
func BasicAuth() *AuthConfig {
	return &AuthConfig{
		Type:   AuthTypeHTTP,
		Scheme: HTTPSchemeBasic,
	}
}

// OAuth2Auth creates an OAuth 2.0 authentication configuration
//
// Example:
//
//	auth := scalar.OAuth2Auth().
//		WithAuthorizationCode(
//			"https://example.com/oauth/authorize",
//			"https://example.com/oauth/token",
//			map[string]string{
//				"read":  "Read access",
//				"write": "Write access",
//			},
//		)
func OAuth2Auth() *AuthConfig {
	return &AuthConfig{
		Type:  AuthTypeOAuth2,
		Flows: make(map[string]OAuth2FlowConfig),
	}
}

// OpenIDConnectAuth creates an OpenID Connect authentication configuration
//
// Example:
//
//	auth := scalar.OpenIDConnectAuth("https://example.com/.well-known/openid-configuration").
//		WithDescription("OpenID Connect authentication")
func OpenIDConnectAuth(url string) *AuthConfig {
	return &AuthConfig{
		Type:             AuthTypeOpenID,
		OpenIDConnectURL: url,
	}
}

// WithDescription adds a description to the authentication configuration
func (a *AuthConfig) WithDescription(description string) *AuthConfig {
	a.Description = description
	return a
}

// WithBearerFormat specifies the format of the bearer token (e.g., "JWT")
func (a *AuthConfig) WithBearerFormat(format string) *AuthConfig {
	a.BearerFormat = format
	return a
}

// WithAuthorizationCode adds OAuth 2.0 Authorization Code flow
func (a *AuthConfig) WithAuthorizationCode(authURL, tokenURL string, scopes map[string]string) *AuthConfig {
	if a.Flows == nil {
		a.Flows = make(map[string]OAuth2FlowConfig)
	}
	a.Flows[string(OAuth2FlowAuthorizationCode)] = OAuth2FlowConfig{
		AuthorizationURL: authURL,
		TokenURL:         tokenURL,
		Scopes:           scopes,
	}
	return a
}

// WithClientCredentials adds OAuth 2.0 Client Credentials flow
func (a *AuthConfig) WithClientCredentials(tokenURL string, scopes map[string]string) *AuthConfig {
	if a.Flows == nil {
		a.Flows = make(map[string]OAuth2FlowConfig)
	}
	a.Flows[string(OAuth2FlowClientCredentials)] = OAuth2FlowConfig{
		TokenURL: tokenURL,
		Scopes:   scopes,
	}
	return a
}

// WithImplicit adds OAuth 2.0 Implicit flow
func (a *AuthConfig) WithImplicit(authURL string, scopes map[string]string) *AuthConfig {
	if a.Flows == nil {
		a.Flows = make(map[string]OAuth2FlowConfig)
	}
	a.Flows[string(OAuth2FlowImplicit)] = OAuth2FlowConfig{
		AuthorizationURL: authURL,
		Scopes:           scopes,
	}
	return a
}

// WithPassword adds OAuth 2.0 Password flow
func (a *AuthConfig) WithPassword(tokenURL string, scopes map[string]string) *AuthConfig {
	if a.Flows == nil {
		a.Flows = make(map[string]OAuth2FlowConfig)
	}
	a.Flows[string(OAuth2FlowPassword)] = OAuth2FlowConfig{
		TokenURL: tokenURL,
		Scopes:   scopes,
	}
	return a
}

// WithRefreshURL adds a refresh URL to an OAuth 2.0 flow
func (a *AuthConfig) WithRefreshURL(flow OAuth2Flow, refreshURL string) *AuthConfig {
	if a.Flows != nil {
		if flowConfig, exists := a.Flows[string(flow)]; exists {
			flowConfig.RefreshURL = refreshURL
			a.Flows[string(flow)] = flowConfig
		}
	}
	return a
}

// ToJSON converts the authentication configuration to a JSON string
func (a *AuthConfig) ToJSON() (string, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MustJSON converts the authentication configuration to a JSON string, panics on error
func (a *AuthConfig) MustJSON() string {
	data, err := a.ToJSON()
	if err != nil {
		panic(err)
	}
	return data
}

// MultipleAuth creates a configuration with multiple authentication methods
//
// Example:
//
//	auth := scalar.MultipleAuth(map[string]*AuthConfig{
//		"apiKey": scalar.APIKeyAuth("X-API-Key", scalar.APIKeyInHeader),
//		"bearer": scalar.BearerAuth().WithBearerFormat("JWT"),
//	})
func MultipleAuth(configs map[string]*AuthConfig) (string, error) {
	data, err := json.Marshal(configs)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
