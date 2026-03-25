package scalar

import (
	"strings"
	"testing"
)

func TestCustomLogoAndFavicon(t *testing.T) {
	tests := []struct {
		name             string
		options          *Options
		expectFavicon    bool
		expectLogo       bool
		faviconURL       string
		logoURL          string
	}{
		{
			name: "With both logo and favicon URLs",
			options: &Options{
				SpecContent: map[string]interface{}{
					"openapi": "3.0.0",
					"info": map[string]interface{}{
						"title":   "Test API",
						"version": "1.0.0",
					},
				},
				CustomOptions: CustomOptions{
					PageTitle:  "Custom API",
					LogoURL:    "https://example.com/logo.png",
					FaviconURL: "https://example.com/favicon.ico",
				},
			},
			expectFavicon: true,
			expectLogo:    true,
			faviconURL:    "https://example.com/favicon.ico",
			logoURL:       "https://example.com/logo.png",
		},
		{
			name: "With base64 logo and favicon",
			options: &Options{
				SpecContent: map[string]interface{}{
					"openapi": "3.0.0",
					"info": map[string]interface{}{
						"title":   "Test API",
						"version": "1.0.0",
					},
				},
				CustomOptions: CustomOptions{
					LogoURL:    "data:image/png;base64,iVBORw0KGgoAAAANS...",
					FaviconURL: "data:image/x-icon;base64,AAABAAEAEBAAAAEAIAAoBAAAFgAAA...",
				},
			},
			expectFavicon: true,
			expectLogo:    true,
			faviconURL:    "data:image/x-icon;base64,AAABAAEAEBAAAAEAIAAoBAAAFgAAA...",
			logoURL:       "data:image/png;base64,iVBORw0KGgoAAAANS...",
		},
		{
			name: "With only favicon",
			options: &Options{
				SpecContent: map[string]interface{}{
					"openapi": "3.0.0",
					"info": map[string]interface{}{
						"title":   "Test API",
						"version": "1.0.0",
					},
				},
				CustomOptions: CustomOptions{
					FaviconURL: "/favicon.ico",
				},
			},
			expectFavicon: true,
			expectLogo:    false,
			faviconURL:    "/favicon.ico",
		},
		{
			name: "With only logo",
			options: &Options{
				SpecContent: map[string]interface{}{
					"openapi": "3.0.0",
					"info": map[string]interface{}{
						"title":   "Test API",
						"version": "1.0.0",
					},
				},
				CustomOptions: CustomOptions{
					LogoURL: "/logo.svg",
				},
			},
			expectFavicon: false,
			expectLogo:    true,
			logoURL:       "/logo.svg",
		},
		{
			name: "Without logo and favicon",
			options: &Options{
				SpecContent: map[string]interface{}{
					"openapi": "3.0.0",
					"info": map[string]interface{}{
						"title":   "Test API",
						"version": "1.0.0",
					},
				},
				CustomOptions: CustomOptions{
					PageTitle: "API Docs",
				},
			},
			expectFavicon: false,
			expectLogo:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := ApiReferenceHTML(tt.options)
			if err != nil {
				t.Fatalf("ApiReferenceHTML() error = %v", err)
			}

			// Check favicon in HTML
			if tt.expectFavicon {
				expectedFaviconTag := `<link rel="icon" href="` + tt.faviconURL + `" />`
				if !strings.Contains(html, expectedFaviconTag) {
					t.Errorf("Expected favicon tag not found in HTML.\nExpected: %s\nGot HTML snippet: %s",
						expectedFaviconTag, extractHeadSection(html))
				}
			} else {
				if strings.Contains(html, `<link rel="icon"`) {
					t.Error("Unexpected favicon tag found in HTML")
				}
			}

			// Check logo in configuration
			if tt.expectLogo {
				// Logo is passed through JSON configuration
				if !strings.Contains(html, tt.logoURL) {
					t.Errorf("Expected logo URL not found in HTML configuration: %s", tt.logoURL)
				}
			}
		})
	}
}

func TestFaviconTypes(t *testing.T) {
	tests := []struct {
		name       string
		faviconURL string
		valid      bool
	}{
		{"PNG favicon", "/favicon.png", true},
		{"ICO favicon", "/favicon.ico", true},
		{"SVG favicon", "/favicon.svg", true},
		{"Base64 PNG", "data:image/png;base64,iVBORw0KGgoAAAANS...", true},
		{"Base64 ICO", "data:image/x-icon;base64,AAA...", true},
		{"Remote URL", "https://example.com/favicon.png", true},
		{"Relative path", "../assets/favicon.ico", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &Options{
				SpecContent: map[string]interface{}{
					"openapi": "3.0.0",
					"info": map[string]interface{}{
						"title":   "Test",
						"version": "1.0.0",
					},
				},
				CustomOptions: CustomOptions{
					FaviconURL: tt.faviconURL,
				},
			}

			html, err := ApiReferenceHTML(options)
			if err != nil {
				t.Fatalf("ApiReferenceHTML() error = %v", err)
			}

			if tt.valid {
				expectedTag := `<link rel="icon" href="` + tt.faviconURL + `" />`
				if !strings.Contains(html, expectedTag) {
					t.Errorf("Expected favicon tag not found: %s", expectedTag)
				}
			}
		})
	}
}

// Helper function to extract head section from HTML for debugging
func extractHeadSection(html string) string {
	start := strings.Index(html, "<head>")
	end := strings.Index(html, "</head>")
	if start == -1 || end == -1 {
		return ""
	}
	return html[start : end+7]
}
