package scalar

import (
	"strings"
	"testing"
)

const testSpec = `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
  description: A simple test API
  contact:
    name: API Support
    email: support@example.com
    url: https://example.com/support
servers:
  - url: https://api.example.com/v1
    description: Production server
paths:
  /users:
    get:
      summary: List users
      description: Get a list of all users
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
          required: false
          description: Maximum number of users to return
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              example:
                users:
                  - id: 1
                    name: John Doe
    post:
      summary: Create user
      description: Create a new user
      requestBody:
        required: true
        content:
          application/json:
            example:
              name: Jane Doe
              email: jane@example.com
      responses:
        '201':
          description: User created
components:
  schemas:
    User:
      type: object
      description: User object
      properties:
        id:
          type: integer
        name:
          type: string
        email:
          type: string
`

func TestExportToMarkdown(t *testing.T) {
	options := &ExportOptions{
		Format:          ExportFormatMarkdown,
		IncludeTOC:      true,
		IncludeExamples: true,
	}

	result, err := ExportSpec(testSpec, options)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Check for key elements
	if !strings.Contains(result, "# Test API") {
		t.Error("Missing title")
	}
	if !strings.Contains(result, "## Table of Contents") {
		t.Error("Missing table of contents")
	}
	if !strings.Contains(result, "## Base URL") {
		t.Error("Missing base URL section")
	}
	if !strings.Contains(result, "## Endpoints") {
		t.Error("Missing endpoints section")
	}
	if !strings.Contains(result, "/users") {
		t.Error("Missing /users endpoint")
	}
	if !strings.Contains(result, "GET") {
		t.Error("Missing GET method")
	}
	if !strings.Contains(result, "POST") {
		t.Error("Missing POST method")
	}
	if !strings.Contains(result, "## Schemas") {
		t.Error("Missing schemas section")
	}
}

func TestExportToMarkdownWithoutTOC(t *testing.T) {
	options := &ExportOptions{
		Format:     ExportFormatMarkdown,
		IncludeTOC: false,
	}

	result, err := ExportSpec(testSpec, options)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if strings.Contains(result, "## Table of Contents") {
		t.Error("Should not include table of contents")
	}
}

func TestExportToHTML(t *testing.T) {
	options := &ExportOptions{
		Format: ExportFormatHTML,
	}

	result, err := ExportSpec(testSpec, options)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !strings.Contains(result, "<!DOCTYPE html>") {
		t.Error("Missing HTML doctype")
	}
	if !strings.Contains(result, "<html>") {
		t.Error("Missing HTML tag")
	}
	if !strings.Contains(result, "Test API") {
		t.Error("Missing title in HTML")
	}
}

func TestExportToJSON(t *testing.T) {
	options := &ExportOptions{
		Format: ExportFormatJSON,
	}

	result, err := ExportSpec(testSpec, options)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !strings.Contains(result, `"openapi"`) {
		t.Error("Missing OpenAPI field in JSON")
	}
	if !strings.Contains(result, `"Test API"`) {
		t.Error("Missing title in JSON")
	}
}

func TestExportToYAML(t *testing.T) {
	options := &ExportOptions{
		Format: ExportFormatYAML,
	}

	result, err := ExportSpec(testSpec, options)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !strings.Contains(result, "openapi:") {
		t.Error("Missing OpenAPI field in YAML")
	}
	if !strings.Contains(result, "Test API") {
		t.Error("Missing title in YAML")
	}
}

func TestDefaultExportOptions(t *testing.T) {
	options := DefaultExportOptions()

	if options.Format != ExportFormatMarkdown {
		t.Errorf("Expected default format %s, got %s", ExportFormatMarkdown, options.Format)
	}
	if !options.IncludeTOC {
		t.Error("Expected IncludeTOC to be true by default")
	}
	if !options.IncludeExamples {
		t.Error("Expected IncludeExamples to be true by default")
	}
}

func TestExportWithCustomTitle(t *testing.T) {
	options := &ExportOptions{
		Format: ExportFormatMarkdown,
		Title:  "Custom Title",
	}

	result, err := ExportSpec(testSpec, options)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !strings.Contains(result, "# Custom Title") {
		t.Error("Missing custom title")
	}
}

func TestExportInvalidSpec(t *testing.T) {
	invalidSpec := `invalid yaml content`

	options := DefaultExportOptions()
	_, err := ExportSpec(invalidSpec, options)

	if err == nil {
		t.Error("Expected error for invalid spec")
	}
}

func TestExportOpenAPI2(t *testing.T) {
	swagger2Spec := `
swagger: "2.0"
info:
  title: Swagger 2.0 API
  version: 1.0.0
  description: A Swagger 2.0 spec
host: api.example.com
basePath: /v1
paths:
  /test:
    get:
      summary: Test endpoint
      responses:
        200:
          description: Success
`

	options := DefaultExportOptions()
	result, err := ExportSpec(swagger2Spec, options)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !strings.Contains(result, "Swagger 2.0 API") {
		t.Error("Missing title from Swagger 2.0 spec")
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	options := &ExportOptions{
		Format: "unsupported",
	}

	_, err := ExportSpec(testSpec, options)
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported export format") {
		t.Errorf("Unexpected error message: %v", err)
	}
}
