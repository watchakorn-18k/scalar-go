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
//
// With UI authentication:
//
//	app.Use("/docs", fiber.Handler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//		UIUsername: "admin",
//		UIPassword: "secret",
//	}))
func Handler(options *scalar.Options) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if UI authentication is enabled
		if options.IsUIAuthEnabled() {
			auth := c.Get("Authorization")
			username, password, ok := scalar.ParseBasicAuth(auth)

			if !ok || !options.ValidateUICredentials(username, password) {
				c.Set("WWW-Authenticate", `Basic realm="Scalar API Documentation"`)
				return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
			}
		}

		htmlContent, err := scalar.ApiReferenceHTML(options)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(htmlContent)
	}
}
