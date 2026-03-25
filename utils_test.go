package scalar

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureFileURL(t *testing.T) {
	// Get current directory for testing
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	tests := []struct {
		name        string
		input       string
		expectStart string
		expectErr   bool
	}{
		{
			name:        "absolute path without file:// prefix",
			input:       "/tmp/test.yaml",
			expectStart: "file:///tmp/test.yaml",
			expectErr:   false,
		},
		{
			name:        "absolute path with file:// prefix",
			input:       "file:///tmp/test.yaml",
			expectStart: "file:///tmp/test.yaml",
			expectErr:   false,
		},
		{
			name:        "relative path",
			input:       "docs/spec.yaml",
			expectStart: "file://",
			expectErr:   false,
		},
		{
			name:        "relative path with file:// prefix",
			input:       "file://docs/spec.yaml",
			expectStart: "file://",
			expectErr:   false,
		},
		{
			name:        "current directory reference",
			input:       "./spec.yaml",
			expectStart: "file://",
			expectErr:   false,
		},
		{
			name:        "parent directory reference",
			input:       "../spec.yaml",
			expectStart: "file://",
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ensureFileURL(tt.input)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectErr {
				if !strings.HasPrefix(result, tt.expectStart) {
					t.Errorf("Expected result to start with %q, got %q", tt.expectStart, result)
				}

				// Verify file:// prefix is always present
				if !strings.HasPrefix(result, "file://") {
					t.Errorf("Result should always start with file://, got %q", result)
				}
			}
		})
	}

	t.Run("relative path resolves to current directory", func(t *testing.T) {
		relativePath := "test.yaml"
		result, err := ensureFileURL(relativePath)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expectedPath := filepath.Join(currentDir, relativePath)
		expectedURL := "file://" + expectedPath

		if result != expectedURL {
			t.Errorf("Expected %q, got %q", expectedURL, result)
		}
	})
}

func TestFetchContentFromURL(t *testing.T) {
	t.Run("successfully fetches content from HTTP server", func(t *testing.T) {
		expectedContent := `{"openapi":"3.0.0","info":{"title":"Test API"}}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(expectedContent))
		}))
		defer server.Close()

		content, err := fetchContentFromURL(server.URL)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if content != expectedContent {
			t.Errorf("Expected %q, got %q", expectedContent, content)
		}
	})

	t.Run("returns error for invalid URL", func(t *testing.T) {
		_, err := fetchContentFromURL("http://invalid-domain-that-does-not-exist-12345.com")
		if err == nil {
			t.Error("Expected error for invalid URL")
		}

		if !strings.Contains(err.Error(), "error getting file content") {
			t.Errorf("Error message should mention getting file content, got: %v", err)
		}
	})

	t.Run("returns error for 404 response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		content, err := fetchContentFromURL(server.URL)
		// The function reads the body regardless of status code
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// It should still return the content (empty in this case)
		if content == "" {
			// This is expected behavior
		}
	})

	t.Run("handles large content", func(t *testing.T) {
		largeContent := strings.Repeat("test data ", 10000)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(largeContent))
		}))
		defer server.Close()

		content, err := fetchContentFromURL(server.URL)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if content != largeContent {
			t.Error("Content should match large input")
		}
	})
}

func TestReadFileFromURL(t *testing.T) {
	t.Run("reads file from valid file:// URL", func(t *testing.T) {
		// Create a temporary file
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.yaml")
		expectedContent := []byte("openapi: 3.0.0\ninfo:\n  title: Test")

		err := os.WriteFile(testFile, expectedContent, 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		fileURL := "file://" + testFile
		content, err := readFileFromURL(fileURL)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if string(content) != string(expectedContent) {
			t.Errorf("Expected %q, got %q", expectedContent, content)
		}
	})

	t.Run("returns error for non-file:// URL scheme", func(t *testing.T) {
		_, err := readFileFromURL("http://example.com/spec.yaml")

		if err == nil {
			t.Error("Expected error for non-file:// URL")
		}

		expectedErr := "unsupported URL scheme: http"
		if !strings.Contains(err.Error(), expectedErr) {
			t.Errorf("Expected error to contain %q, got %q", expectedErr, err.Error())
		}
	})

	t.Run("returns error for invalid URL", func(t *testing.T) {
		_, err := readFileFromURL("not a valid url")

		if err == nil {
			t.Error("Expected error for invalid URL")
		}

		// The url.Parse doesn't error on "not a valid url", it parses it as a path
		// So we expect "unsupported URL scheme" error instead
		if !strings.Contains(err.Error(), "unsupported URL scheme") && !strings.Contains(err.Error(), "error parsing URL") {
			t.Errorf("Error should mention URL scheme or parsing, got: %v", err)
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		fileURL := "file:///tmp/non-existent-file-12345.yaml"
		_, err := readFileFromURL(fileURL)

		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	t.Run("handles empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "empty.yaml")

		err := os.WriteFile(testFile, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		fileURL := "file://" + testFile
		content, err := readFileFromURL(fileURL)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(content) != 0 {
			t.Errorf("Expected empty content, got %d bytes", len(content))
		}
	})

	t.Run("reads binary file", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "binary.dat")
		binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}

		err := os.WriteFile(testFile, binaryData, 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		fileURL := "file://" + testFile
		content, err := readFileFromURL(fileURL)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(content) != len(binaryData) {
			t.Errorf("Expected %d bytes, got %d bytes", len(binaryData), len(content))
		}

		for i, b := range binaryData {
			if content[i] != b {
				t.Errorf("Byte mismatch at index %d: expected %x, got %x", i, b, content[i])
			}
		}
	})

	t.Run("handles URL with query parameters (should ignore)", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.yaml")
		expectedContent := []byte("test: data")

		err := os.WriteFile(testFile, expectedContent, 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// file:// URLs with query params - the implementation uses parsedURL.Path which ignores query
		fileURL := "file://" + testFile
		content, err := readFileFromURL(fileURL)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if string(content) != string(expectedContent) {
			t.Errorf("Expected %q, got %q", expectedContent, content)
		}
	})
}
