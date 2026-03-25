package scalar

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

// ValidationError represents an error that occurred during spec validation
type ValidationError struct {
	Message string
	Details []string
}

func (e *ValidationError) Error() string {
	if len(e.Details) == 0 {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Message, strings.Join(e.Details, ", "))
}

// ValidateSpec validates an OpenAPI specification
// Supports OpenAPI 2.0 (Swagger), 3.0, and 3.1
func ValidateSpec(content string) error {
	// Try to detect the spec version
	var versionCheck map[string]interface{}

	// Try JSON first
	if err := json.Unmarshal([]byte(content), &versionCheck); err != nil {
		// Try YAML
		if err := yaml.Unmarshal([]byte(content), &versionCheck); err != nil {
			return &ValidationError{
				Message: "Invalid spec format",
				Details: []string{"Content must be valid JSON or YAML"},
			}
		}
	}

	// Check for OpenAPI version
	if swagger, ok := versionCheck["swagger"].(string); ok && strings.HasPrefix(swagger, "2.") {
		return validateOpenAPI2(content)
	}

	if openapi, ok := versionCheck["openapi"].(string); ok && (strings.HasPrefix(openapi, "3.0") || strings.HasPrefix(openapi, "3.1")) {
		return validateOpenAPI3(content)
	}

	return &ValidationError{
		Message: "Unsupported or missing spec version",
		Details: []string{"Spec must contain 'swagger: 2.0' or 'openapi: 3.0.x/3.1.x'"},
	}
}

// validateOpenAPI2 validates an OpenAPI 2.0 (Swagger) specification
func validateOpenAPI2(content string) error {
	var doc openapi2.T

	// Try JSON first
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		// Try YAML
		if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
			return &ValidationError{
				Message: "Failed to parse OpenAPI 2.0 spec",
				Details: []string{err.Error()},
			}
		}
	}

	// Validate required fields
	var errors []string

	if doc.Swagger == "" {
		errors = append(errors, "missing 'swagger' field")
	}

	if doc.Info.Title == "" {
		errors = append(errors, "missing 'info.title' field")
	}
	if doc.Info.Version == "" {
		errors = append(errors, "missing 'info.version' field")
	}

	if doc.Paths == nil || len(doc.Paths) == 0 {
		errors = append(errors, "missing or empty 'paths' field")
	}

	if len(errors) > 0 {
		return &ValidationError{
			Message: "OpenAPI 2.0 validation failed",
			Details: errors,
		}
	}

	return nil
}

// validateOpenAPI3 validates an OpenAPI 3.0 or 3.1 specification
func validateOpenAPI3(content string) error {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	var doc *openapi3.T
	var err error

	// Try JSON first
	doc, err = loader.LoadFromData([]byte(content))
	if err != nil {
		return &ValidationError{
			Message: "Failed to parse OpenAPI 3.x spec",
			Details: []string{err.Error()},
		}
	}

	// Validate the document
	if err := doc.Validate(loader.Context); err != nil {
		return &ValidationError{
			Message: "OpenAPI 3.x validation failed",
			Details: []string{err.Error()},
		}
	}

	return nil
}

// ValidateSpecFromFile reads and validates a spec file
func ValidateSpecFromFile(filePath string) error {
	content, err := readSpecFile(filePath)
	if err != nil {
		return err
	}
	return ValidateSpec(string(content))
}

// readSpecFile reads a spec file from local path or URL
func readSpecFile(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		content, err := fetchContentFromURL(path)
		if err != nil {
			return nil, err
		}
		return []byte(content), nil
	}

	// Handle local file
	fileURL, err := ensureFileURL(path)
	if err != nil {
		return nil, err
	}

	return readFileFromURL(fileURL)
}
