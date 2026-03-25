package chi

import (
	"net/http"

	"github.com/watchakorn-18k/scalar-go"
)

// ExportHandler creates a Chi/net/http handler for exporting API documentation
//
// Example usage:
//
//	r.Get("/docs/export/markdown", chi.ExportHandler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//	}, scalar.ExportFormatMarkdown))
//
//	r.Get("/docs/export/pdf", chi.ExportHandler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//	}, scalar.ExportFormatHTML))
func ExportHandler(options *scalar.Options, format scalar.ExportFormat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			specContent = content
		} else {
			http.Error(w, "No spec content or URL provided", http.StatusBadRequest)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set content type based on format
		switch format {
		case scalar.ExportFormatMarkdown:
			w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
			w.Header().Set("Content-Disposition", "attachment; filename=api-documentation.md")
		case scalar.ExportFormatHTML:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Content-Disposition", "attachment; filename=api-documentation.html")
		case scalar.ExportFormatJSON:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Content-Disposition", "attachment; filename=openapi.json")
		case scalar.ExportFormatYAML:
			w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
			w.Header().Set("Content-Disposition", "attachment; filename=openapi.yaml")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}
