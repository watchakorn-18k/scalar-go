## Scalar Go

<img src="_examples/image.png" align="center">


## Overview 🌐

The Scalar package serves as a provider for the [Scalar](https://github.com/scalar/scalar) project. It offers a comprehensive suite of functions designed for generating API references in HTML format, specializing in JSON data handling and web presentation customization. This includes functions to serialize options into JSON, manage HTML escaping, and dynamically handle different types of specification content.

## Installation 📦

```bash
go get -u github.com/watchakorn-18k/scalar-go@latest
```

## Features 🚀

### JSON Serialization and HTML Escaping

- **safeJSONConfiguration** 🔒: Serializes configuration options into JSON format and escapes HTML characters to prevent XSS attacks. This ensures that the JSON can be safely embedded within HTML documents.

### Specification Content Handling

- **specContentHandler** 📝: Dynamically handles different types of specification content. It can process content as a function returning a map, a direct map, or a plain string, converting the specification content into JSON format suitable for web use.

### HTML Generation

- **ApiReferenceHTML** 📄: Generates a complete HTML document for API reference. It allows for extensive customization, including themes, layouts, and CDN configuration. It handles both direct specification content and content fetched from a URL, providing error handling for missing specifications.

### Spec Validation ✅

- **ValidateSpec**: Validates OpenAPI/Swagger specifications before rendering to ensure they are correct and well-formed
- **ValidateSpecFromFile**: Validates spec files from local paths or URLs
- Supports OpenAPI 2.0 (Swagger), 3.0, and 3.1
- Supports both JSON and YAML formats
- Provides detailed error messages with validation failures
- Can be enabled/disabled per request with `ValidateSpec` option

### Authentication Configuration 🔐

Built-in authentication helpers for easy configuration:

- **API Key Authentication**: Header, Query, or Cookie-based
- **HTTP Authentication**: Basic, Bearer (JWT), Digest
- **OAuth 2.0**: Authorization Code, Client Credentials, Implicit, Password flows
- **OpenID Connect**: Full OpenID Connect support
- **Multiple Auth Methods**: Combine multiple authentication methods

### UI Access Protection 🔒

Protect your API documentation UI with Basic Authentication:

- **Password Protection**: Prevent unauthorized access to your API documentation
- **Optional by Default**: No protection required unless explicitly configured
- **Secure Implementation**: Uses constant-time comparison to prevent timing attacks
- **Framework Support**: Works with all supported frameworks (Fiber, Gin, Echo, Chi)

## Customization Options ⚙️

The package allows extensive customization of the generated API reference through the `Options` struct, supporting:

- **CDN**: URL of the CDN to load additional scripts or styles.
- **Theme**: Customizable themes for styling the API reference.
- **Layout**: Choice between modern and classic layout designs.
- **SpecURL**: URL from which the specification content can be fetched.
- **Dark Mode**: Option to enable a dark theme for the API reference.

### Themes and Styles 🎨

- Provides a default set of CSS styles for both light and dark themes.
- Allows custom CSS injections to tailor the appearance to specific branding or aesthetic requirements.

## Error Handling 🛠️

- Robust error handling for scenarios where necessary parameters like `SpecURL` or `SpecContent` are missing.
- Errors during the fetching of specification content from URLs are properly managed and reported.

## Usage 📚

### Framework Middleware 🎯

Scalar Go provides ready-to-use middleware for popular Go web frameworks:

#### Fiber

```go
package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/watchakorn-18k/scalar-go"
	scalarfiber "github.com/watchakorn-18k/scalar-go/middleware/fiber"
)

func main() {
	app := fiber.New()

	app.Use("/docs", scalarfiber.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API Documentation",
		},
		DarkMode: true,
	}))

	log.Fatal(app.Listen(":3000"))
}
```

#### Gin

```go
package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/watchakorn-18k/scalar-go"
	scalargin "github.com/watchakorn-18k/scalar-go/middleware/gin"
)

func main() {
	r := gin.Default()

	r.GET("/docs", scalargin.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API Documentation",
		},
		DarkMode: true,
	}))

	log.Fatal(r.Run(":3000"))
}
```

#### Echo

```go
package main

import (
	"log"
	"github.com/labstack/echo/v4"
	"github.com/watchakorn-18k/scalar-go"
	scalarecho "github.com/watchakorn-18k/scalar-go/middleware/echo"
)

func main() {
	e := echo.New()

	e.GET("/docs", scalarecho.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API Documentation",
		},
		DarkMode: true,
	}))

	log.Fatal(e.Start(":3000"))
}
```

#### Chi

```go
package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/watchakorn-18k/scalar-go"
	scalarchi "github.com/watchakorn-18k/scalar-go/middleware/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/docs", scalarchi.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API Documentation",
		},
		DarkMode: true,
	}))

	log.Fatal(http.ListenAndServe(":3000", r))
}
```

### Manual Usage (without middleware)

You can also use the core `ApiReferenceHTML` function directly:

```go
package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/watchakorn-18k/scalar-go"
)

func main() {
	app := fiber.New()

	app.Use("/api/docs", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.yaml",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
			},
			DarkMode: true,
		})

		if err != nil {
			return err
		}
		c.Type("html")
		return c.SendString(htmlContent)
	})

	log.Fatal(app.Listen(":3000"))
}
```

## Authentication Configuration 🔐

Scalar Go provides easy-to-use authentication helpers:

### API Key Authentication

```go
auth := scalar.APIKeyAuth("X-API-Key", scalar.APIKeyInHeader).
	WithDescription("API key for authentication")

r.GET("/docs", scalargin.Handler(&scalar.Options{
	SpecURL:        "./swagger.yaml",
	Authentication: auth.MustJSON(),
}))
```

### Bearer Token (JWT) Authentication

```go
auth := scalar.BearerAuth().
	WithBearerFormat("JWT").
	WithDescription("JWT Bearer token authentication")

r.GET("/docs", scalargin.Handler(&scalar.Options{
	SpecURL:        "./swagger.yaml",
	Authentication: auth.MustJSON(),
}))
```

### Basic HTTP Authentication

```go
auth := scalar.BasicAuth().
	WithDescription("Basic HTTP authentication")

r.GET("/docs", scalargin.Handler(&scalar.Options{
	SpecURL:        "./swagger.yaml",
	Authentication: auth.MustJSON(),
}))
```

### OAuth 2.0 Authentication

```go
auth := scalar.OAuth2Auth().
	WithAuthorizationCode(
		"https://example.com/oauth/authorize",
		"https://example.com/oauth/token",
		map[string]string{
			"read":  "Read access",
			"write": "Write access",
		},
	).
	WithClientCredentials(
		"https://example.com/oauth/token",
		map[string]string{
			"api": "API access",
		},
	)

r.GET("/docs", scalargin.Handler(&scalar.Options{
	SpecURL:        "./swagger.yaml",
	Authentication: auth.MustJSON(),
}))
```

### OpenID Connect Authentication

```go
auth := scalar.OpenIDConnectAuth(
	"https://example.com/.well-known/openid-configuration",
).WithDescription("OpenID Connect authentication")

r.GET("/docs", scalargin.Handler(&scalar.Options{
	SpecURL:        "./swagger.yaml",
	Authentication: auth.MustJSON(),
}))
```

### Multiple Authentication Methods

```go
multiAuth, err := scalar.MultipleAuth(map[string]*scalar.AuthConfig{
	"apiKey": scalar.APIKeyAuth("X-API-Key", scalar.APIKeyInHeader),
	"bearer": scalar.BearerAuth().WithBearerFormat("JWT"),
	"oauth2": scalar.OAuth2Auth().WithAuthorizationCode(...),
})

r.GET("/docs", scalargin.Handler(&scalar.Options{
	SpecURL:        "./swagger.yaml",
	Authentication: multiAuth,
}))
```

## UI Access Protection 🔒

Protect your Scalar API documentation UI from unauthorized access using Basic Authentication. This is useful when you want to prevent others from viewing your API specifications.

### Basic Usage

Simply add `UIUsername` and `UIPassword` to your options:

```go
app.Use("/docs", scalarfiber.Handler(&scalar.Options{
	SpecURL:    "./swagger.yaml",
	UIUsername: "admin",
	UIPassword: "secret123",
	CustomOptions: scalar.CustomOptions{
		PageTitle: "Protected API Documentation",
	},
}))
```

Now when users access `/docs`, they'll be prompted for credentials via HTTP Basic Authentication.

### Public and Protected Documentation

You can have both public and protected documentation routes:

```go
// Public documentation (no auth required)
app.Use("/docs/public", scalarfiber.Handler(&scalar.Options{
	SpecURL: "./public-api.yaml",
	CustomOptions: scalar.CustomOptions{
		PageTitle: "Public API Documentation",
	},
}))

// Protected documentation (requires auth)
app.Use("/docs/private", scalarfiber.Handler(&scalar.Options{
	SpecURL:    "./internal-api.yaml",
	UIUsername: "admin",
	UIPassword: "secret123",
	CustomOptions: scalar.CustomOptions{
		PageTitle: "Internal API Documentation",
	},
}))
```

### Security Features

- **Constant-Time Comparison**: Prevents timing attacks by using `crypto/subtle`
- **Standard Basic Auth**: Compatible with all browsers and HTTP clients
- **Optional by Default**: No authentication unless explicitly configured
- **Per-Route Configuration**: Different credentials for different documentation routes

### Example with All Frameworks

#### Fiber
```go
app.Use("/docs", scalarfiber.Handler(&scalar.Options{
	SpecURL:    "./swagger.yaml",
	UIUsername: "admin",
	UIPassword: "secret123",
}))
```

#### Gin
```go
r.GET("/docs", scalargin.Handler(&scalar.Options{
	SpecURL:    "./swagger.yaml",
	UIUsername: "admin",
	UIPassword: "secret123",
}))
```

#### Echo
```go
e.GET("/docs", scalarecho.Handler(&scalar.Options{
	SpecURL:    "./swagger.yaml",
	UIUsername: "admin",
	UIPassword: "secret123",
}))
```

#### Chi
```go
r.Get("/docs", scalarchi.Handler(&scalar.Options{
	SpecURL:    "./swagger.yaml",
	UIUsername: "admin",
	UIPassword: "secret123",
}))
```

### Important Notes

- **Not for API Authentication**: This feature protects the documentation UI only, not your actual API endpoints
- **Use Strong Passwords**: Always use strong, unique passwords in production
- **HTTPS Recommended**: Basic Auth sends credentials in base64 encoding, use HTTPS in production
- **Environment Variables**: Store credentials in environment variables, not in code:

```go
app.Use("/docs", scalarfiber.Handler(&scalar.Options{
	SpecURL:    "./swagger.yaml",
	UIUsername: os.Getenv("DOCS_USERNAME"),
	UIPassword: os.Getenv("DOCS_PASSWORD"),
}))
```

### Available Authentication Types

| Type | Helper Function | Description |
|------|----------------|-------------|
| API Key | `APIKeyAuth(name, location)` | API key in header, query, or cookie |
| Bearer | `BearerAuth()` | Bearer token authentication (e.g., JWT) |
| Basic | `BasicAuth()` | Basic HTTP authentication |
| OAuth 2.0 | `OAuth2Auth()` | OAuth 2.0 with multiple flows |
| OpenID Connect | `OpenIDConnectAuth(url)` | OpenID Connect authentication |

### Builder Methods

All authentication configurations support method chaining:

- `.WithDescription(desc)` - Add a description
- `.WithBearerFormat(format)` - Specify bearer token format (e.g., "JWT")
- `.WithAuthorizationCode(authURL, tokenURL, scopes)` - Add OAuth 2.0 Authorization Code flow
- `.WithClientCredentials(tokenURL, scopes)` - Add OAuth 2.0 Client Credentials flow
- `.WithImplicit(authURL, scopes)` - Add OAuth 2.0 Implicit flow
- `.WithPassword(tokenURL, scopes)` - Add OAuth 2.0 Password flow
- `.WithRefreshURL(flow, refreshURL)` - Add refresh URL to OAuth 2.0 flow
- `.ToJSON()` - Convert to JSON string (returns error)
- `.MustJSON()` - Convert to JSON string (panics on error)

## Spec Validation ✅

Scalar Go provides built-in validation for OpenAPI/Swagger specifications to ensure they are correct before rendering. This helps catch errors early and provides clear feedback when something is wrong.

### Supported Formats

- ✅ OpenAPI 3.1
- ✅ OpenAPI 3.0
- ✅ OpenAPI 2.0 (Swagger)
- ✅ JSON format
- ✅ YAML format

### Enable Validation in Middleware

Simply set `ValidateSpec: true` in your options:

```go
app.Use("/docs", scalarfiber.Handler(&scalar.Options{
	SpecURL:      "./swagger.yaml",
	ValidateSpec: true, // Enable validation
	CustomOptions: scalar.CustomOptions{
		PageTitle: "My API Documentation",
	},
}))
```

If the spec is invalid, the middleware will return an error response with details about what's wrong.

### Manual Validation

You can also validate specs programmatically:

#### Validate from File

```go
// Validate a spec file (supports both local paths and URLs)
err := scalar.ValidateSpecFromFile("./swagger.yaml")
if err != nil {
	log.Printf("Spec validation failed: %v", err)
}
```

#### Validate from Content

```go
// Validate spec content directly
specContent := `
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
paths:
  /hello:
    get:
      responses:
        '200':
          description: Success
`

err := scalar.ValidateSpec(specContent)
if err != nil {
	log.Printf("Spec validation failed: %v", err)
}
```

#### Validation Endpoint Example

```go
app.Get("/validate", func(c *fiber.Ctx) error {
	err := scalar.ValidateSpecFromFile("./swagger.yaml")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"valid": false,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"valid":   true,
		"message": "Spec is valid!",
	})
})
```

### Error Messages

When validation fails, you'll get detailed error messages:

```
spec validation failed: OpenAPI 3.x validation failed: field info is required
```

For OpenAPI 2.0:
```
spec validation failed: OpenAPI 2.0 validation failed: missing 'info.title' field, missing 'paths' field
```

### Performance Considerations

- Validation adds a small overhead (typically a few milliseconds)
- For production with static specs, validate once during startup
- For dynamic specs or development, enable `ValidateSpec: true`
- Validation is cached per request, so multiple calls with the same content are fast

### When to Use Validation

**Always validate when:**
- ✅ Developing and testing API documentation
- ✅ Accepting user-uploaded spec files
- ✅ Fetching specs from external sources
- ✅ Running in development/staging environments

**Optional in:**
- ⚠️ Production with static, pre-validated specs (for performance)
- ⚠️ High-traffic endpoints (validate during deployment instead)

```
## License

This project is forked from [MarceloPetrucio/go-scalar-api-reference](https://github.com/MarceloPetrucio/go-scalar-api-reference) and is licensed under the MIT License.