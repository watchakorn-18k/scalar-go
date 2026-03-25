package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/watchakorn-18k/scalar-go"
)

// ExportHandler creates a Gin handler for exporting API documentation
//
// Example usage:
//
//	r.GET("/docs/export/markdown", gin.ExportHandler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//	}, scalar.ExportFormatMarkdown))
//
//	r.GET("/docs/export/pdf", gin.ExportHandler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//	}, scalar.ExportFormatHTML))
func ExportHandler(options *scalar.Options, format scalar.ExportFormat) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read spec content
		var specContent string
		var err error

		if options.SpecContent != nil {
			// Use provided spec content
			switch content := options.SpecContent.(type) {
			case string:
				specContent = content
			case func() map[string]interface{}:
				// Handle function that returns spec
				result := content()
				jsonData, _ := scalar.ExportToJSON(result)
				specContent = jsonData
			case map[string]interface{}:
				jsonData, _ := scalar.ExportToJSON(content)
				specContent = jsonData
			}
		} else if options.SpecURL != "" {
			// Read from file/URL
			content, err := scalar.ReadSpecContent(options.SpecURL)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			specContent = content
		} else {
			c.String(http.StatusBadRequest, "No spec content or URL provided")
			return
		}

		// Export options
		exportOptions := &scalar.ExportOptions{
			Format:          format,
			IncludeTOC:      true,
			IncludeExamples: true,
			Title:           options.CustomOptions.PageTitle,
		}

		// Export spec
		result, err := scalar.ExportSpec(specContent, exportOptions)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Set content type based on format
		switch format {
		case scalar.ExportFormatMarkdown:
			c.Header("Content-Type", "text/markdown; charset=utf-8")
			c.Header("Content-Disposition", "attachment; filename=api-documentation.md")
		case scalar.ExportFormatHTML:
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Header("Content-Disposition", "attachment; filename=api-documentation.html")
		case scalar.ExportFormatJSON:
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.Header("Content-Disposition", "attachment; filename=openapi.json")
		case scalar.ExportFormatYAML:
			c.Header("Content-Type", "application/x-yaml; charset=utf-8")
			c.Header("Content-Disposition", "attachment; filename=openapi.yaml")
		}

		c.String(http.StatusOK, result)
	}
}
