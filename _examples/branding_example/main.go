package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/watchakorn-18k/scalar-go"
	scalarfiber "github.com/watchakorn-18k/scalar-go/middleware/fiber"
)

func main() {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Serve static files for demo
	app.Static("/assets", "./assets")

	// Example 1: Using URL for logo and favicon
	app.Use("/docs/url", scalarfiber.Handler(&scalar.Options{
		SpecURL: "../simple_api/openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle:  "My Branded API - URL Assets",
			LogoURL:    "/assets/logo.png",
			FaviconURL: "/assets/favicon.ico",
		},
		DarkMode: true,
	}))

	// Example 2: Using base64 data URI for logo and favicon
	// This is useful when you want to embed the logo directly without external files
	app.Use("/docs/base64", scalarfiber.Handler(&scalar.Options{
		SpecURL: "../simple_api/openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My Branded API - Base64 Assets",
			// Small SVG logo as base64
			LogoURL: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8Y2lyY2xlIGN4PSI1MCIgY3k9IjUwIiByPSI0MCIgZmlsbD0iIzRBOTBFMiIvPgogIDx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBmb250LXNpemU9IjM2IiBmaWxsPSJ3aGl0ZSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPkE8L3RleHQ+Cjwvc3ZnPg==",
			// Small SVG favicon as base64
			FaviconURL: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIiIGhlaWdodD0iMzIiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgPGNpcmNsZSBjeD0iMTYiIGN5PSIxNiIgcj0iMTQiIGZpbGw9IiM0QTkwRTIiLz4KICA8dGV4dCB4PSI1MCUiIHk9IjUwJSIgZm9udC1zaXplPSIxOCIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiIGR5PSIuM2VtIj5BPC90ZXh0Pgo8L3N2Zz4=",
		},
		DarkMode: true,
	}))

	// Example 3: Only custom logo
	app.Use("/docs/logo-only", scalarfiber.Handler(&scalar.Options{
		SpecURL: "../simple_api/openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "API with Custom Logo Only",
			LogoURL:   "https://raw.githubusercontent.com/scalar/scalar/main/assets/scalar-logo.svg",
		},
		DarkMode: true,
	}))

	// Example 4: Only custom favicon
	app.Use("/docs/favicon-only", scalarfiber.Handler(&scalar.Options{
		SpecURL: "../simple_api/openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle:  "API with Custom Favicon Only",
			FaviconURL: "data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>🚀</text></svg>",
		},
		DarkMode: true,
	}))

	// Example 5: No custom branding (default)
	app.Use("/docs/default", scalarfiber.Handler(&scalar.Options{
		SpecURL: "../simple_api/openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Default API Documentation",
		},
		DarkMode: true,
	}))

	// API routes
	api := app.Group("/api")
	api.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from Branded API!",
		})
	})

	log.Println("🚀 Server running on http://localhost:3000")
	log.Println("")
	log.Println("📚 API Documentation Examples:")
	log.Println("  1. URL Assets:      http://localhost:3000/docs/url")
	log.Println("  2. Base64 Assets:   http://localhost:3000/docs/base64")
	log.Println("  3. Logo Only:       http://localhost:3000/docs/logo-only")
	log.Println("  4. Favicon Only:    http://localhost:3000/docs/favicon-only")
	log.Println("  5. Default (no branding): http://localhost:3000/docs/default")
	log.Println("")
	log.Fatal(app.Listen(":3000"))
}
