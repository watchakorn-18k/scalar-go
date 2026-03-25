package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/watchakorn-18k/scalar-go"
	scalargin "github.com/watchakorn-18k/scalar-go/middleware/gin"
)

func main() {
	r := gin.Default()

	// Example 1: API Key Authentication
	r.GET("/docs/apikey", scalargin.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "API Key Auth Example",
		},
		Authentication: scalar.APIKeyAuth("X-API-Key", scalar.APIKeyInHeader).
			WithDescription("API key for authentication").
			MustJSON(),
		DarkMode: true,
	}))

	// Example 2: Bearer Token (JWT) Authentication
	r.GET("/docs/bearer", scalargin.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Bearer Auth Example",
		},
		Authentication: scalar.BearerAuth().
			WithBearerFormat("JWT").
			WithDescription("JWT Bearer token authentication").
			MustJSON(),
		DarkMode: true,
	}))

	// Example 3: Basic HTTP Authentication
	r.GET("/docs/basic", scalargin.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Basic Auth Example",
		},
		Authentication: scalar.BasicAuth().
			WithDescription("Basic HTTP authentication").
			MustJSON(),
		DarkMode: true,
	}))

	// Example 4: OAuth 2.0 Authentication
	r.GET("/docs/oauth2", scalargin.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "OAuth 2.0 Example",
		},
		Authentication: scalar.OAuth2Auth().
			WithAuthorizationCode(
				"https://example.com/oauth/authorize",
				"https://example.com/oauth/token",
				map[string]string{
					"read":  "Read access to protected resources",
					"write": "Write access to protected resources",
					"admin": "Admin access to all resources",
				},
			).
			WithClientCredentials(
				"https://example.com/oauth/token",
				map[string]string{
					"api": "API access",
				},
			).
			WithDescription("OAuth 2.0 authentication with multiple flows").
			MustJSON(),
		DarkMode: true,
	}))

	// Example 5: OpenID Connect Authentication
	r.GET("/docs/openid", scalargin.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "OpenID Connect Example",
		},
		Authentication: scalar.OpenIDConnectAuth(
			"https://example.com/.well-known/openid-configuration",
		).
			WithDescription("OpenID Connect authentication").
			MustJSON(),
		DarkMode: true,
	}))

	// Example 6: Multiple Authentication Methods
	multiAuth, err := scalar.MultipleAuth(map[string]*scalar.AuthConfig{
		"apiKey": scalar.APIKeyAuth("X-API-Key", scalar.APIKeyInHeader).
			WithDescription("API Key authentication"),
		"bearer": scalar.BearerAuth().
			WithBearerFormat("JWT").
			WithDescription("JWT Bearer token"),
		"oauth2": scalar.OAuth2Auth().
			WithAuthorizationCode(
				"https://example.com/oauth/authorize",
				"https://example.com/oauth/token",
				map[string]string{
					"read":  "Read access",
					"write": "Write access",
				},
			),
	})
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/docs/multi", scalargin.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Multiple Auth Methods",
		},
		Authentication: multiAuth,
		DarkMode:       true,
	}))

	// API routes
	api := r.Group("/api")
	{
		api.GET("/public", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Public endpoint"})
		})

		api.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Protected endpoint"})
		})
	}

	log.Println("Server running on http://localhost:3000")
	log.Println("Available documentation:")
	log.Println("  - API Key Auth:    http://localhost:3000/docs/apikey")
	log.Println("  - Bearer Auth:     http://localhost:3000/docs/bearer")
	log.Println("  - Basic Auth:      http://localhost:3000/docs/basic")
	log.Println("  - OAuth 2.0:       http://localhost:3000/docs/oauth2")
	log.Println("  - OpenID Connect:  http://localhost:3000/docs/openid")
	log.Println("  - Multiple Auth:   http://localhost:3000/docs/multi")

	log.Fatal(r.Run(":3000"))
}
