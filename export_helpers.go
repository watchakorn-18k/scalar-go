package scalar

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// ExportToJSON is a helper function to convert any data to JSON string
func ExportToJSON(data interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// ReadSpecContent reads spec content from URL or file path
func ReadSpecContent(specURL string) (string, error) {
	if specURL == "" {
		return "", fmt.Errorf("specURL is empty")
	}

	// Check if it's a URL
	if strings.HasPrefix(specURL, "http://") || strings.HasPrefix(specURL, "https://") {
		content, err := fetchContentFromURL(specURL)
		if err != nil {
			return "", err
		}
		return content, nil
	}

	// Handle local file
	cleanPath := filepath.Clean(specURL)
	parts := strings.Split(cleanPath, string(filepath.Separator))
	specPath := filepath.Join(parts...)
	absPath, err := filepath.Abs(specPath)
	if err != nil {
		return "", err
	}

	urlPath, err := ensureFileURL(filepath.ToSlash(absPath))
	if err != nil {
		return "", err
	}

	content, err := readFileFromURL(urlPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
