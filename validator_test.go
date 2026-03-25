package scalar

import (
	"strings"
	"testing"
)

func TestValidateSpec_OpenAPI3_Valid(t *testing.T) {
	validSpec := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: OK
`
	err := ValidateSpec(validSpec)
	if err != nil {
		t.Errorf("Expected valid spec to pass validation, got error: %v", err)
	}
}

func TestValidateSpec_OpenAPI3_ValidJSON(t *testing.T) {
	validSpec := `{
  "openapi": "3.0.0",
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "paths": {
    "/test": {
      "get": {
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}`
	err := ValidateSpec(validSpec)
	if err != nil {
		t.Errorf("Expected valid JSON spec to pass validation, got error: %v", err)
	}
}

func TestValidateSpec_OpenAPI2_Valid(t *testing.T) {
	validSpec := `{
  "swagger": "2.0",
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "paths": {
    "/test": {
      "get": {
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}`
	err := ValidateSpec(validSpec)
	if err != nil {
		t.Errorf("Expected valid Swagger 2.0 spec to pass validation, got error: %v", err)
	}
}

func TestValidateSpec_MissingVersion(t *testing.T) {
	invalidSpec := `{
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "paths": {}
}`
	err := ValidateSpec(invalidSpec)
	if err == nil {
		t.Error("Expected error for spec missing version field")
	}
	if !strings.Contains(err.Error(), "Unsupported or missing spec version") {
		t.Errorf("Expected error about missing version, got: %v", err)
	}
}

func TestValidateSpec_InvalidJSON(t *testing.T) {
	invalidSpec := `{this is not valid json or yaml`
	err := ValidateSpec(invalidSpec)
	if err == nil {
		t.Error("Expected error for invalid JSON/YAML")
	}
}

func TestValidateSpec_OpenAPI3_MissingInfo(t *testing.T) {
	invalidSpec := `{
  "openapi": "3.0.0",
  "paths": {
    "/test": {
      "get": {
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}`
	err := ValidateSpec(invalidSpec)
	if err == nil {
		t.Error("Expected error for spec missing info field")
	}
}

func TestValidateSpec_OpenAPI2_MissingTitle(t *testing.T) {
	invalidSpec := `{
  "swagger": "2.0",
  "info": {
    "version": "1.0.0"
  },
  "paths": {
    "/test": {
      "get": {
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}`
	err := ValidateSpec(invalidSpec)
	if err == nil {
		t.Error("Expected error for Swagger 2.0 spec missing title")
	}
	if !strings.Contains(err.Error(), "info.title") {
		t.Errorf("Expected error about missing title, got: %v", err)
	}
}

func TestValidateSpec_OpenAPI2_MissingPaths(t *testing.T) {
	invalidSpec := `{
  "swagger": "2.0",
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  }
}`
	err := ValidateSpec(invalidSpec)
	if err == nil {
		t.Error("Expected error for spec missing paths")
	}
	if !strings.Contains(err.Error(), "paths") {
		t.Errorf("Expected error about missing paths, got: %v", err)
	}
}

func TestValidationError_Error(t *testing.T) {
	ve := &ValidationError{
		Message: "Test error",
		Details: []string{"detail1", "detail2"},
	}
	errMsg := ve.Error()
	if !strings.Contains(errMsg, "Test error") {
		t.Errorf("Expected error message to contain 'Test error', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "detail1") || !strings.Contains(errMsg, "detail2") {
		t.Errorf("Expected error message to contain details, got: %s", errMsg)
	}
}

func TestValidationError_ErrorNoDetails(t *testing.T) {
	ve := &ValidationError{
		Message: "Test error",
	}
	errMsg := ve.Error()
	if errMsg != "Test error" {
		t.Errorf("Expected error message 'Test error', got: %s", errMsg)
	}
}
