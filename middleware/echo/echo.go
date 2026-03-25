package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/watchakorn-18k/scalar-go"
)

// Handler creates an Echo middleware handler for Scalar API documentation
//
// Example usage:
//
//	e := echo.New()
//	e.GET("/docs", echo.Handler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//		CustomOptions: scalar.CustomOptions{
//			PageTitle: "My API Documentation",
//		},
//	}))
func Handler(options *scalar.Options) echo.HandlerFunc {
	return func(c echo.Context) error {
		htmlContent, err := scalar.ApiReferenceHTML(options)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.HTML(http.StatusOK, htmlContent)
	}
}
