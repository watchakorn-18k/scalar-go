package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/watchakorn-18k/scalar-go"
	scalarecho "github.com/watchakorn-18k/scalar-go/middleware/echo"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Scalar API Documentation
	e.GET("/docs", scalarecho.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Echo API Documentation",
		},
		DarkMode: true,
	}))

	// API routes
	api := e.Group("/api")
	{
		api.GET("/hello", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"message": "Hello from Echo!",
			})
		})
	}

	log.Println("Server running on http://localhost:3000")
	log.Println("API Docs available at http://localhost:3000/docs")
	log.Fatal(e.Start(":3000"))
}
