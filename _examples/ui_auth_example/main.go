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

	// Public API Documentation (no authentication)
	app.Use("/docs/public", scalarfiber.Handler(&scalar.Options{
		SpecURL: "../simple_api/openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Public API Documentation",
		},
		DarkMode: true,
	}))

	// Protected API Documentation (requires authentication)
	app.Use("/docs/private", scalarfiber.Handler(&scalar.Options{
		SpecURL: "../simple_api/openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Private API Documentation",
		},
		DarkMode:   true,
		UIUsername: "admin",
		UIPassword: "secret123",
	}))

	// API routes
	api := app.Group("/api")
	api.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from Fiber!",
		})
	})

	log.Println("Server running on http://localhost:3000")
	log.Println("Public API Docs: http://localhost:3000/docs/public")
	log.Println("Private API Docs: http://localhost:3000/docs/private (username: admin, password: secret123)")
	log.Fatal(app.Listen(":3000"))
}
