package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/watchakorn-18k/scalar-go"
	scalargin "github.com/watchakorn-18k/scalar-go/middleware/gin"
)

func main() {
	r := gin.Default()

	// Scalar API Documentation
	r.GET("/docs", scalargin.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Gin API Documentation",
		},
		DarkMode: true,
	}))

	// API routes
	api := r.Group("/api")
	{
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello from Gin!",
			})
		})
	}

	log.Println("Server running on http://localhost:3000")
	log.Println("API Docs available at http://localhost:3000/docs")
	log.Fatal(r.Run(":3000"))
}
