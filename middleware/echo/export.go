package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/watchakorn-18k/scalar-go"
)

// ExportHandler creates an Echo handler for exporting API documentation
//
// Example usage:
//
//	e.GET("/docs/export/markdown", echo.ExportHandler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//	}, scalar.ExportFormatMarkdown))
//
//	e.GET("/docs/export/pdf", echo.ExportHandler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//	}, scalar.ExportFormatHTML))
func ExportHandler(options *scalar.Options, format scalar.ExportFormat) echo.HandlerFunc {
	return func(c echo.Context) error {
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
				return c.String(http.StatusInternalServerError, err.Error())
			}
			specContent = content
		} else {
			return c.String(http.StatusBadRequest, "No spec content or URL provided")
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
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Set content type based on format
		switch format {
		case scalar.ExportFormatMarkdown:
			c.Response().Header().Set("Content-Type", "text/markdown; charset=utf-8")
			c.Response().Header().Set("Content-Disposition", "attachment; filename=api-documentation.md")
		case scalar.ExportFormatHTML:
			c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
			c.Response().Header().Set("Content-Disposition", "attachment; filename=api-documentation.html")
		case scalar.ExportFormatJSON:
			c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
			c.Response().Header().Set("Content-Disposition", "attachment; filename=openapi.json")
		case scalar.ExportFormatYAML:
			c.Response().Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
			c.Response().Header().Set("Content-Disposition", "attachment; filename=openapi.yaml")
		}

		return c.String(http.StatusOK, result)
	}
}
