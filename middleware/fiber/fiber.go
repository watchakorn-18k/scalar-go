package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/watchakorn-18k/scalar-go"
)

// Handler creates a Fiber middleware handler for Scalar API documentation
//
// Example usage:
//
//	app := fiber.New()
//	app.Use("/docs", fiber.Handler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//		CustomOptions: scalar.CustomOptions{
//			PageTitle: "My API Documentation",
//		},
//	}))
func Handler(options *scalar.Options) fiber.Handler {
	return func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(options)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(htmlContent)
	}
}
