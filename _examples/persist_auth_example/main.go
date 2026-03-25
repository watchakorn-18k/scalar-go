package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/watchakorn-18k/scalar-go"
	scalarFiber "github.com/watchakorn-18k/scalar-go/middleware/fiber"
)

func main() {
	app := fiber.New()

	// Example 1: With Auth Persistence (Bearer tokens saved to localStorage)
	app.Get("/docs", scalarFiber.Handler(&scalar.Options{
		SpecURL:     "./swagger.yaml",
		PersistAuth: true, // 🔑 Enable auth persistence!
		CustomOptions: scalar.CustomOptions{
			PageTitle: "API Docs with Persistent Auth",
		},
		DarkMode: true,
	}))

	// Example 2: Without Auth Persistence (tokens lost on refresh)
	app.Get("/docs-no-persist", scalarFiber.Handler(&scalar.Options{
		SpecURL:     "./swagger.yaml",
		PersistAuth: false, // Tokens will be lost on refresh
		CustomOptions: scalar.CustomOptions{
			PageTitle: "API Docs without Persistence",
		},
		DarkMode: false,
	}))

	// Example 3: With validation AND persistence
	app.Get("/docs-full", scalarFiber.Handler(&scalar.Options{
		SpecURL:      "./swagger.yaml",
		ValidateSpec: true,  // Validate spec
		PersistAuth:  true,  // Persist auth tokens
		UIUsername:   "admin", // Protect UI
		UIPassword:   "secret",
		CustomOptions: scalar.CustomOptions{
			PageTitle:  "Full-Featured API Docs",
			FaviconURL: "/favicon.ico",
		},
	}))

	// Simple API endpoint
	api := app.Group("/api")

	// Protected endpoint - requires Bearer token
	api.Get("/protected", func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		// Simple validation (in production, verify the token properly)
		if auth != "Bearer test-token-123" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		return c.JSON(fiber.Map{
			"message": "You are authenticated!",
			"user":    "test-user",
		})
	})

	// Public endpoint
	api.Get("/public", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "This is a public endpoint",
		})
	})

	fmt.Println("🚀 Server started on http://localhost:3000")
	fmt.Println("\n📚 Documentation URLs:")
	fmt.Println("  ✅ With Auth Persistence:    http://localhost:3000/docs")
	fmt.Println("  ❌ Without Auth Persistence: http://localhost:3000/docs-no-persist")
	fmt.Println("  🔒 Full Featured (Protected): http://localhost:3000/docs-full (admin/secret)")
	fmt.Println("\n🧪 Test API Endpoints:")
	fmt.Println("  🔓 Public:    GET http://localhost:3000/api/public")
	fmt.Println("  🔐 Protected: GET http://localhost:3000/api/protected")
	fmt.Println("                (requires: Authorization: Bearer test-token-123)")
	fmt.Println("\n💡 Tips:")
	fmt.Println("  1. Go to /docs")
	fmt.Println("  2. Try the /api/protected endpoint and add bearer token: test-token-123")
	fmt.Println("  3. Refresh the page - your token is still there! ✨")
	fmt.Println("  4. Compare with /docs-no-persist where token disappears on refresh")
	fmt.Println("  5. Open browser console and run: scalarClearAuth() to clear saved tokens")

	log.Fatal(app.Listen(":3000"))
}
