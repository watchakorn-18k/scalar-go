package main

import (
	"log"
	"simple_api/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/watchakorn-18k/scalar-go"
	scalarfiber "github.com/watchakorn-18k/scalar-go/middleware/fiber"
)

func main() {
	app := fiber.New()
	middlewares.Logger(app)

	// Using Scalar Fiber middleware
	app.Use("/api/docs", scalarfiber.Handler(&scalar.Options{
		SpecURL: "./docs/swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Simple API",
		},
		DarkMode: true,
	}))

	app.Use(recover.New())
	app.Use(cors.New())
	api := app.Group("/api/")
	api.Get("/hello", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"message": "Hello World",
		})
	})
	log.Fatal(app.Listen(":3000"))
}
