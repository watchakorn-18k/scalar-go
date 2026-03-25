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

	// Example 1: With validation enabled (recommended for production)
	app.Get("/docs", scalarFiber.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "API Docs with Validation",
		},
		DarkMode:     true,
		ValidateSpec: true, // Enable spec validation
	}))

	// Example 2: Without validation (faster, but no safety checks)
	app.Get("/docs-no-validation", scalarFiber.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "API Docs without Validation",
		},
		DarkMode:     false,
		ValidateSpec: false, // Disable validation
	}))

	// Example 3: Manual validation
	app.Get("/validate", func(c *fiber.Ctx) error {
		// Validate spec file
		err := scalar.ValidateSpecFromFile("./swagger.yaml")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Spec is valid!",
		})
	})

	// Example 4: Validate spec content
	app.Get("/validate-content", func(c *fiber.Ctx) error {
		specContent := `
openapi: 3.0.0
info:
  title: Sample API
  version: 1.0.0
paths:
  /hello:
    get:
      summary: Returns a greeting
      responses:
        '200':
          description: A greeting message
`
		err := scalar.ValidateSpec(specContent)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Spec content is valid!",
		})
	})

	fmt.Println("🚀 Server started on http://localhost:3000")
	fmt.Println("📚 API Docs (with validation): http://localhost:3000/docs")
	fmt.Println("📚 API Docs (no validation): http://localhost:3000/docs-no-validation")
	fmt.Println("✅ Validate endpoint: http://localhost:3000/validate")
	fmt.Println("✅ Validate content endpoint: http://localhost:3000/validate-content")

	log.Fatal(app.Listen(":3000"))
}
