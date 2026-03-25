package scalar

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSafeJSONConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		options  *Options
		contains []string
	}{
		{
			name: "escapes double quotes",
			options: &Options{
				Theme:   ThemeDefault,
				Layout:  LayoutModern,
				SpecURL: "https://example.com/spec.yaml",
			},
			contains: []string{"&quot;", "theme&quot;", "layout&quot;"},
		},
		{
			name: "handles empty options",
			options: &Options{
				DarkMode: true,
			},
			contains: []string{"&quot;darkMode&quot;:true"},
		},
		{
			name: "handles custom options",
			options: &Options{
				CustomOptions: CustomOptions{
					PageTitle: "Test API",
				},
			},
			contains: []string{"&quot;pageTitle&quot;", "Test API"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := safeJSONConfiguration(tt.options)

			// Check if result is valid escaped JSON
			if strings.Contains(result, `"`) && !strings.Contains(result, `&quot;`) {
				t.Error("JSON should not contain unescaped double quotes")
			}

			// Check expected content
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, but it didn't. Got: %s", expected, result)
				}
			}
		})
	}
}

func TestSpecContentHandler(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "handles function returning map",
			input: func() map[string]interface{} {
				return map[string]interface{}{
					"openapi": "3.0.0",
					"info": map[string]interface{}{
						"title": "Test API",
					},
				}
			},
			expected: `{"info":{"title":"Test API"},"openapi":"3.0.0"}`,
		},
		{
			name: "handles map directly",
			input: map[string]interface{}{
				"swagger": "2.0",
				"host":    "example.com",
			},
			expected: `{"host":"example.com","swagger":"2.0"}`,
		},
		{
			name:     "handles string directly",
			input:    `{"openapi":"3.0.0"}`,
			expected: `{"openapi":"3.0.0"}`,
		},
		{
			name:     "handles nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "handles invalid type",
			input:    12345,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := specContentHandler(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestApiReferenceHTML(t *testing.T) {
	t.Run("returns error when both SpecURL and SpecContent are missing", func(t *testing.T) {
		options := &Options{}
		_, err := ApiReferenceHTML(options)

		if err == nil {
			t.Error("Expected error when both SpecURL and SpecContent are missing")
		}

		expectedErr := "specURL or specContent must be provided"
		if err.Error() != expectedErr {
			t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
		}
	})

	t.Run("generates HTML with SpecContent", func(t *testing.T) {
		options := &Options{
			SpecContent: `{"openapi":"3.0.0"}`,
			CustomOptions: CustomOptions{
				PageTitle: "My API Docs",
			},
		}

		html, err := ApiReferenceHTML(options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify HTML structure
		expectedParts := []string{
			"<!DOCTYPE html>",
			"<html>",
			"<title>My API Docs</title>",
			`<script id="api-reference"`,
			`{"openapi":"3.0.0"}`,
			DefaultCDN,
		}

		for _, part := range expectedParts {
			if !strings.Contains(html, part) {
				t.Errorf("HTML should contain %q", part)
			}
		}
	})

	t.Run("uses default page title", func(t *testing.T) {
		options := &Options{
			SpecContent: `{"openapi":"3.0.0"}`,
		}

		html, err := ApiReferenceHTML(options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !strings.Contains(html, "<title>Scalar API Reference</title>") {
			t.Error("HTML should contain default page title")
		}
	})

	t.Run("applies custom CDN", func(t *testing.T) {
		customCDN := "https://custom-cdn.example.com/scalar"
		options := &Options{
			SpecContent: `{"openapi":"3.0.0"}`,
			CDN:         customCDN,
		}

		html, err := ApiReferenceHTML(options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !strings.Contains(html, customCDN) {
			t.Errorf("HTML should contain custom CDN: %s", customCDN)
		}
	})

	t.Run("removes custom theme CSS when theme is set", func(t *testing.T) {
		options := &Options{
			SpecContent: `{"openapi":"3.0.0"}`,
			Theme:       ThemeMoon,
		}

		html, err := ApiReferenceHTML(options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// When theme is set, custom CSS should be empty
		if strings.Contains(html, "--scalar-color-1") {
			t.Error("HTML should not contain custom theme CSS when theme is explicitly set")
		}
	})

	t.Run("includes custom theme CSS when no theme is set", func(t *testing.T) {
		options := &Options{
			SpecContent: `{"openapi":"3.0.0"}`,
		}

		html, err := ApiReferenceHTML(options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !strings.Contains(html, "--scalar-color-1") {
			t.Error("HTML should contain custom theme CSS when no theme is set")
		}
	})

	t.Run("fetches content from HTTP URL", func(t *testing.T) {
		// Create a test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"openapi": "3.0.0",
				"info":    "Test API",
			})
		}))
		defer server.Close()

		options := &Options{
			SpecURL: server.URL,
		}

		html, err := ApiReferenceHTML(options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !strings.Contains(html, "openapi") {
			t.Error("HTML should contain content fetched from URL")
		}
	})

	t.Run("fetches content from local file", func(t *testing.T) {
		// Create a temporary file
		tmpDir := t.TempDir()
		specFile := filepath.Join(tmpDir, "test-spec.yaml")
		specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0`

		err := os.WriteFile(specFile, []byte(specContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Get absolute path
		absPath, err := filepath.Abs(specFile)
		if err != nil {
			t.Fatalf("Failed to get absolute path: %v", err)
		}

		options := &Options{
			SpecURL: absPath,
		}

		html, err := ApiReferenceHTML(options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !strings.Contains(html, "Test API") {
			t.Error("HTML should contain content from local file")
		}
	})

	t.Run("handles dark mode", func(t *testing.T) {
		options := &Options{
			SpecContent: `{"openapi":"3.0.0"}`,
			DarkMode:    true,
		}

		html, err := ApiReferenceHTML(options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Check that darkMode is in the configuration
		if !strings.Contains(html, "darkMode&quot;:true") {
			t.Error("HTML should include darkMode configuration")
		}
	})

	t.Run("handles all layout types", func(t *testing.T) {
		layouts := []ReferenceLayoutType{LayoutModern, LayoutClassic}

		for _, layout := range layouts {
			options := &Options{
				SpecContent: `{"openapi":"3.0.0"}`,
				Layout:      layout,
			}

			html, err := ApiReferenceHTML(options)
			if err != nil {
				t.Fatalf("Unexpected error for layout %s: %v", layout, err)
			}

			if !strings.Contains(html, string(layout)) {
				t.Errorf("HTML should contain layout: %s", layout)
			}
		}
	})
}

func TestApiReferenceHTML_SpecContentTypes(t *testing.T) {
	t.Run("handles map SpecContent", func(t *testing.T) {
		options := &Options{
			SpecContent: map[string]interface{}{
				"openapi": "3.0.0",
				"info": map[string]interface{}{
					"title": "Map Test",
				},
			},
		}

		html, err := ApiReferenceHTML(options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !strings.Contains(html, "Map Test") {
			t.Error("HTML should contain content from map")
		}
	})

	t.Run("handles function SpecContent", func(t *testing.T) {
		options := &Options{
			SpecContent: func() map[string]interface{} {
				return map[string]interface{}{
					"openapi": "3.0.0",
					"info": map[string]interface{}{
						"title": "Function Test",
					},
				}
			},
		}

		html, err := ApiReferenceHTML(options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !strings.Contains(html, "Function Test") {
			t.Error("HTML should contain content from function")
		}
	})
}
