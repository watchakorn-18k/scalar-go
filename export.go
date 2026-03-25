package scalar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

// ExportFormat represents the export format type
type ExportFormat string

const (
	ExportFormatMarkdown ExportFormat = "markdown"
	ExportFormatHTML     ExportFormat = "html"
	ExportFormatJSON     ExportFormat = "json"
	ExportFormatYAML     ExportFormat = "yaml"
)

// ExportOptions configures export behavior
type ExportOptions struct {
	Format          ExportFormat
	IncludeTOC      bool   // Include Table of Contents (for Markdown)
	IncludeExamples bool   // Include request/response examples
	Title           string // Custom title for the export
}

// DefaultExportOptions returns default export options
func DefaultExportOptions() *ExportOptions {
	return &ExportOptions{
		Format:          ExportFormatMarkdown,
		IncludeTOC:      true,
		IncludeExamples: true,
	}
}

// ExportSpec exports an OpenAPI specification to the specified format
func ExportSpec(specContent string, options *ExportOptions) (string, error) {
	if options == nil {
		options = DefaultExportOptions()
	}

	// Parse the spec
	var doc *openapi3.T
	var err error

	// Try to load as OpenAPI 3.x first
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	doc, err = loader.LoadFromData([]byte(specContent))

	if err != nil {
		// For OpenAPI 2.0, we'll handle it differently
		// Just return an error for now, or we could implement separate handling
		return "", fmt.Errorf("failed to parse spec as OpenAPI 3.x. OpenAPI 2.0 specs should be converted to 3.x first: %w", err)
	}

	switch options.Format {
	case ExportFormatMarkdown:
		return exportToMarkdown(doc, options)
	case ExportFormatHTML:
		return exportToHTML(doc, options)
	case ExportFormatJSON:
		return exportToJSON(doc)
	case ExportFormatYAML:
		return exportToYAML(doc)
	default:
		return "", fmt.Errorf("unsupported export format: %s", options.Format)
	}
}

// ExportSpecFromFile reads a spec file and exports it
func ExportSpecFromFile(filePath string, options *ExportOptions) (string, error) {
	content, err := readSpecFile(filePath)
	if err != nil {
		return "", err
	}
	return ExportSpec(string(content), options)
}

// exportToMarkdown converts OpenAPI spec to Markdown
func exportToMarkdown(doc *openapi3.T, options *ExportOptions) (string, error) {
	var buf bytes.Buffer

	// Title
	title := options.Title
	if title == "" && doc.Info != nil {
		title = doc.Info.Title
	}
	if title != "" {
		buf.WriteString(fmt.Sprintf("# %s\n\n", title))
	}

	// Info section
	if doc.Info != nil {
		if doc.Info.Description != "" {
			buf.WriteString(fmt.Sprintf("%s\n\n", doc.Info.Description))
		}
		if doc.Info.Version != "" {
			buf.WriteString(fmt.Sprintf("**Version:** %s\n\n", doc.Info.Version))
		}
		if doc.Info.Contact != nil {
			buf.WriteString("**Contact:**\n")
			if doc.Info.Contact.Name != "" {
				buf.WriteString(fmt.Sprintf("- Name: %s\n", doc.Info.Contact.Name))
			}
			if doc.Info.Contact.Email != "" {
				buf.WriteString(fmt.Sprintf("- Email: %s\n", doc.Info.Contact.Email))
			}
			if doc.Info.Contact.URL != "" {
				buf.WriteString(fmt.Sprintf("- URL: %s\n", doc.Info.Contact.URL))
			}
			buf.WriteString("\n")
		}
	}

	// Table of Contents
	if options.IncludeTOC && doc.Paths != nil {
		buf.WriteString("## Table of Contents\n\n")
		for path := range doc.Paths.Map() {
			buf.WriteString(fmt.Sprintf("- [%s](#%s)\n", path, strings.ReplaceAll(path, "/", "")))
		}
		buf.WriteString("\n")
	}

	// Base URL
	if len(doc.Servers) > 0 {
		buf.WriteString("## Base URL\n\n")
		for _, server := range doc.Servers {
			buf.WriteString(fmt.Sprintf("- `%s`", server.URL))
			if server.Description != "" {
				buf.WriteString(fmt.Sprintf(" - %s", server.Description))
			}
			buf.WriteString("\n")
		}
		buf.WriteString("\n")
	}

	// Paths/Endpoints
	if doc.Paths != nil {
		buf.WriteString("## Endpoints\n\n")

		for path, pathItem := range doc.Paths.Map() {
			buf.WriteString(fmt.Sprintf("### %s\n\n", path))

			if pathItem.Get != nil {
				writeOperation(&buf, "GET", path, pathItem.Get, options)
			}
			if pathItem.Post != nil {
				writeOperation(&buf, "POST", path, pathItem.Post, options)
			}
			if pathItem.Put != nil {
				writeOperation(&buf, "PUT", path, pathItem.Put, options)
			}
			if pathItem.Delete != nil {
				writeOperation(&buf, "DELETE", path, pathItem.Delete, options)
			}
			if pathItem.Patch != nil {
				writeOperation(&buf, "PATCH", path, pathItem.Patch, options)
			}
		}
	}

	// Schemas
	if doc.Components != nil && doc.Components.Schemas != nil && len(doc.Components.Schemas) > 0 {
		buf.WriteString("## Schemas\n\n")
		for name, schemaRef := range doc.Components.Schemas {
			buf.WriteString(fmt.Sprintf("### %s\n\n", name))
			if schemaRef.Value != nil && schemaRef.Value.Description != "" {
				buf.WriteString(fmt.Sprintf("%s\n\n", schemaRef.Value.Description))
			}
			buf.WriteString("```json\n")
			if schemaRef.Value != nil {
				schemaJSON, _ := json.MarshalIndent(schemaRef.Value, "", "  ")
				buf.Write(schemaJSON)
			}
			buf.WriteString("\n```\n\n")
		}
	}

	return buf.String(), nil
}

// writeOperation writes a single operation to the buffer
func writeOperation(buf *bytes.Buffer, method, path string, op *openapi3.Operation, options *ExportOptions) {
	buf.WriteString(fmt.Sprintf("#### `%s %s`\n\n", method, path))

	if op.Summary != "" {
		buf.WriteString(fmt.Sprintf("**Summary:** %s\n\n", op.Summary))
	}

	if op.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", op.Description))
	}

	// Parameters
	if len(op.Parameters) > 0 {
		buf.WriteString("**Parameters:**\n\n")
		buf.WriteString("| Name | In | Type | Required | Description |\n")
		buf.WriteString("|------|-----|------|----------|-------------|\n")
		for _, paramRef := range op.Parameters {
			if paramRef.Value != nil {
				param := paramRef.Value
				required := "No"
				if param.Required {
					required = "Yes"
				}
				paramType := "string"
				if param.Schema != nil && param.Schema.Value != nil && param.Schema.Value.Type != nil {
					paramType = param.Schema.Value.Type.Slice()[0]
				}
				description := param.Description
				buf.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %s |\n",
					param.Name, param.In, paramType, required, description))
			}
		}
		buf.WriteString("\n")
	}

	// Request Body
	if op.RequestBody != nil && op.RequestBody.Value != nil {
		buf.WriteString("**Request Body:**\n\n")
		if op.RequestBody.Value.Description != "" {
			buf.WriteString(fmt.Sprintf("%s\n\n", op.RequestBody.Value.Description))
		}
		if options.IncludeExamples && op.RequestBody.Value.Content != nil {
			for contentType, mediaType := range op.RequestBody.Value.Content {
				buf.WriteString(fmt.Sprintf("Content-Type: `%s`\n\n", contentType))
				if mediaType.Example != nil {
					buf.WriteString("```json\n")
					exampleJSON, _ := json.MarshalIndent(mediaType.Example, "", "  ")
					buf.Write(exampleJSON)
					buf.WriteString("\n```\n\n")
				}
			}
		}
	}

	// Responses
	if op.Responses != nil && len(op.Responses.Map()) > 0 {
		buf.WriteString("**Responses:**\n\n")
		for status, responseRef := range op.Responses.Map() {
			if responseRef.Value != nil {
				response := responseRef.Value
				description := ""
				if response.Description != nil {
					description = *response.Description
				}
				buf.WriteString(fmt.Sprintf("- **%s**: %s\n", status, description))
				if options.IncludeExamples && response.Content != nil {
					for contentType, mediaType := range response.Content {
						buf.WriteString(fmt.Sprintf("  - Content-Type: `%s`\n", contentType))
						if mediaType.Example != nil {
							buf.WriteString("  ```json\n  ")
							exampleJSON, _ := json.MarshalIndent(mediaType.Example, "  ", "  ")
							buf.Write(exampleJSON)
							buf.WriteString("\n  ```\n")
						}
					}
				}
			}
		}
		buf.WriteString("\n")
	}

	buf.WriteString("---\n\n")
}

// exportToHTML converts OpenAPI spec to HTML
func exportToHTML(doc *openapi3.T, options *ExportOptions) (string, error) {
	// Convert to markdown first, then wrap in HTML
	markdown, err := exportToMarkdown(doc, options)
	if err != nil {
		return "", err
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>%s</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 1200px; margin: 0 auto; padding: 20px; }
        h1 { color: #333; border-bottom: 2px solid #e0234d; padding-bottom: 10px; }
        h2 { color: #555; margin-top: 30px; }
        h3 { color: #666; }
        code { background: #f4f4f4; padding: 2px 6px; border-radius: 3px; }
        pre { background: #f4f4f4; padding: 15px; border-radius: 5px; overflow-x: auto; }
        table { border-collapse: collapse; width: 100%%; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background-color: #f8f8f8; }
    </style>
</head>
<body>
<pre>%s</pre>
</body>
</html>`, doc.Info.Title, markdown)

	return html, nil
}

// exportToJSON converts OpenAPI spec to formatted JSON
func exportToJSON(doc *openapi3.T) (string, error) {
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// exportToYAML converts OpenAPI spec to YAML
func exportToYAML(doc *openapi3.T) (string, error) {
	data, err := yaml.Marshal(doc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
