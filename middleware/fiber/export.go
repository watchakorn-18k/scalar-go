package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/watchakorn-18k/scalar-go"
)

// ExportHandler creates a Fiber handler for exporting API documentation
//
// Example usage:
//
//	app.Get("/docs/export/markdown", fiber.ExportHandler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//	}, scalar.ExportFormatMarkdown))
//
//	app.Get("/docs/export/pdf", fiber.ExportHandler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//	}, scalar.ExportFormatHTML))
func ExportHandler(options *scalar.Options, format scalar.ExportFormat) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			specContent = content
		} else {
			return c.Status(fiber.StatusBadRequest).SendString("No spec content or URL provided")
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
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Set content type based on format
		switch format {
		case scalar.ExportFormatMarkdown:
			c.Set("Content-Type", "text/markdown; charset=utf-8")
			c.Set("Content-Disposition", "attachment; filename=api-documentation.md")
		case scalar.ExportFormatHTML:
			c.Set("Content-Type", "text/html; charset=utf-8")
			c.Set("Content-Disposition", "attachment; filename=api-documentation.html")
		case scalar.ExportFormatJSON:
			c.Set("Content-Type", "application/json; charset=utf-8")
			c.Set("Content-Disposition", "attachment; filename=openapi.json")
		case scalar.ExportFormatYAML:
			c.Set("Content-Type", "application/x-yaml; charset=utf-8")
			c.Set("Content-Disposition", "attachment; filename=openapi.yaml")
		}

		return c.SendString(result)
	}
}
