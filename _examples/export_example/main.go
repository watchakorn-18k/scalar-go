package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/watchakorn-18k/scalar-go"
	scalargin "github.com/watchakorn-18k/scalar-go/middleware/gin"
)

func main() {
	r := gin.Default()

	// Regular API documentation
	r.GET("/docs", scalargin.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API Documentation",
		},
		DarkMode: true,
	}))

	// Export endpoints
	options := &scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API",
		},
	}

	// Export as Markdown
	r.GET("/docs/export/markdown", scalargin.ExportHandler(options, scalar.ExportFormatMarkdown))

	// Export as HTML
	r.GET("/docs/export/html", scalargin.ExportHandler(options, scalar.ExportFormatHTML))

	// Export as JSON (OpenAPI spec)
	r.GET("/docs/export/json", scalargin.ExportHandler(options, scalar.ExportFormatJSON))

	// Export as YAML (OpenAPI spec)
	r.GET("/docs/export/yaml", scalargin.ExportHandler(options, scalar.ExportFormatYAML))

	// API routes
	api := r.Group("/api")
	{
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello World!",
			})
		})

		api.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"users": []gin.H{
					{"id": 1, "name": "John Doe"},
					{"id": 2, "name": "Jane Smith"},
				},
			})
		})

		api.POST("/users", func(c *gin.Context) {
			c.JSON(201, gin.H{
				"message": "User created",
			})
		})
	}

	log.Println("Server running on http://localhost:3000")
	log.Println("Available endpoints:")
	log.Println("  - API Documentation:     http://localhost:3000/docs")
	log.Println("  - Export Markdown:       http://localhost:3000/docs/export/markdown")
	log.Println("  - Export HTML:           http://localhost:3000/docs/export/html")
	log.Println("  - Export JSON (OpenAPI): http://localhost:3000/docs/export/json")
	log.Println("  - Export YAML (OpenAPI): http://localhost:3000/docs/export/yaml")

	log.Fatal(r.Run(":3000"))
}
