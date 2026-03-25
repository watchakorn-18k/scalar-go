package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/watchakorn-18k/scalar-go"
)

// Handler creates a Gin middleware handler for Scalar API documentation
//
// Example usage:
//
//	r := gin.Default()
//	r.GET("/docs", gin.Handler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//		CustomOptions: scalar.CustomOptions{
//			PageTitle: "My API Documentation",
//		},
//	}))
//
// With UI authentication:
//
//	r.GET("/docs", gin.Handler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//		UIUsername: "admin",
//		UIPassword: "secret",
//	}))
func Handler(options *scalar.Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if UI authentication is enabled
		if options.IsUIAuthEnabled() {
			auth := c.GetHeader("Authorization")
			username, password, ok := scalar.ParseBasicAuth(auth)

			if !ok || !options.ValidateUICredentials(username, password) {
				c.Header("WWW-Authenticate", `Basic realm="Scalar API Documentation"`)
				c.String(http.StatusUnauthorized, "Unauthorized")
				return
			}
		}

		htmlContent, err := scalar.ApiReferenceHTML(options)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, htmlContent)
	}
}
