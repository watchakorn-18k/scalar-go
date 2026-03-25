package scalar

import (
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	t.Run("sets default CDN when empty", func(t *testing.T) {
		options := Options{
			SpecURL: "test.yaml",
		}

		result := DefaultOptions(options)

		if result.CDN != DefaultCDN {
			t.Errorf("Expected CDN to be %q, got %q", DefaultCDN, result.CDN)
		}
	})

	t.Run("preserves custom CDN", func(t *testing.T) {
		customCDN := "https://custom-cdn.example.com/scalar"
		options := Options{
			CDN:     customCDN,
			SpecURL: "test.yaml",
		}

		result := DefaultOptions(options)

		if result.CDN != customCDN {
			t.Errorf("Expected CDN to be %q, got %q", customCDN, result.CDN)
		}
	})

	t.Run("sets default Layout when empty", func(t *testing.T) {
		options := Options{
			SpecURL: "test.yaml",
		}

		result := DefaultOptions(options)

		if result.Layout != LayoutModern {
			t.Errorf("Expected Layout to be %q, got %q", LayoutModern, result.Layout)
		}
	})

	t.Run("preserves custom Layout", func(t *testing.T) {
		options := Options{
			Layout:  LayoutClassic,
			SpecURL: "test.yaml",
		}

		result := DefaultOptions(options)

		if result.Layout != LayoutClassic {
			t.Errorf("Expected Layout to be %q, got %q", LayoutClassic, result.Layout)
		}
	})

	t.Run("preserves all other options", func(t *testing.T) {
		options := Options{
			Theme:              ThemeMoon,
			SpecURL:            "https://example.com/spec.yaml",
			Proxy:              "https://proxy.example.com",
			IsEditable:         true,
			ShowSidebar:        true,
			HideModels:         true,
			HideDownloadButton: true,
			DarkMode:           true,
			SearchHotKey:       "Ctrl+K",
			MetaData:           "test-metadata",
			HiddenClients:      []string{"client1", "client2"},
			CustomCss:          ".custom { color: red; }",
			Authentication:     "bearer token",
			PathRouting:        "/api",
			BaseServerURL:      "https://api.example.com",
			WithDefaultFonts:   true,
			CustomOptions: CustomOptions{
				PageTitle: "Custom API Docs",
			},
		}

		result := DefaultOptions(options)

		// Verify all fields are preserved
		if result.Theme != options.Theme {
			t.Errorf("Theme not preserved: expected %q, got %q", options.Theme, result.Theme)
		}
		if result.SpecURL != options.SpecURL {
			t.Errorf("SpecURL not preserved: expected %q, got %q", options.SpecURL, result.SpecURL)
		}
		if result.Proxy != options.Proxy {
			t.Errorf("Proxy not preserved: expected %q, got %q", options.Proxy, result.Proxy)
		}
		if result.IsEditable != options.IsEditable {
			t.Error("IsEditable not preserved")
		}
		if result.ShowSidebar != options.ShowSidebar {
			t.Error("ShowSidebar not preserved")
		}
		if result.HideModels != options.HideModels {
			t.Error("HideModels not preserved")
		}
		if result.HideDownloadButton != options.HideDownloadButton {
			t.Error("HideDownloadButton not preserved")
		}
		if result.DarkMode != options.DarkMode {
			t.Error("DarkMode not preserved")
		}
		if result.SearchHotKey != options.SearchHotKey {
			t.Errorf("SearchHotKey not preserved: expected %q, got %q", options.SearchHotKey, result.SearchHotKey)
		}
		if result.MetaData != options.MetaData {
			t.Errorf("MetaData not preserved: expected %q, got %q", options.MetaData, result.MetaData)
		}
		if len(result.HiddenClients) != len(options.HiddenClients) {
			t.Errorf("HiddenClients length mismatch: expected %d, got %d", len(options.HiddenClients), len(result.HiddenClients))
		}
		if result.CustomCss != options.CustomCss {
			t.Errorf("CustomCss not preserved: expected %q, got %q", options.CustomCss, result.CustomCss)
		}
		if result.Authentication != options.Authentication {
			t.Errorf("Authentication not preserved: expected %q, got %q", options.Authentication, result.Authentication)
		}
		if result.PathRouting != options.PathRouting {
			t.Errorf("PathRouting not preserved: expected %q, got %q", options.PathRouting, result.PathRouting)
		}
		if result.BaseServerURL != options.BaseServerURL {
			t.Errorf("BaseServerURL not preserved: expected %q, got %q", options.BaseServerURL, result.BaseServerURL)
		}
		if result.WithDefaultFonts != options.WithDefaultFonts {
			t.Error("WithDefaultFonts not preserved")
		}
		if result.CustomOptions.PageTitle != options.CustomOptions.PageTitle {
			t.Errorf("PageTitle not preserved: expected %q, got %q", options.CustomOptions.PageTitle, result.CustomOptions.PageTitle)
		}
	})

	t.Run("returns pointer to Options", func(t *testing.T) {
		options := Options{
			SpecURL: "test.yaml",
		}

		result := DefaultOptions(options)

		if result == nil {
			t.Error("Expected non-nil pointer")
		}

		// Verify it's a pointer by modifying the result
		result.DarkMode = true
		if options.DarkMode == true {
			t.Error("Result should be a copy, not a reference to the original")
		}
	})
}

func TestThemeIdConstants(t *testing.T) {
	themes := []ThemeId{
		ThemeDefault,
		ThemeAlternate,
		ThemeMoon,
		ThemePurple,
		ThemeSolarized,
		ThemeBluePlanet,
		ThemeDeepSpace,
		ThemeSaturn,
		ThemeKepler,
		ThemeMars,
		ThemeNone,
		ThemeNil,
	}

	// Verify all theme constants are defined
	expectedValues := map[ThemeId]string{
		ThemeDefault:    "default",
		ThemeAlternate:  "alternate",
		ThemeMoon:       "moon",
		ThemePurple:     "purple",
		ThemeSolarized:  "solarized",
		ThemeBluePlanet: "bluePlanet",
		ThemeDeepSpace:  "deepSpace",
		ThemeSaturn:     "saturn",
		ThemeKepler:     "kepler",
		ThemeMars:       "mars",
		ThemeNone:       "none",
		ThemeNil:        "",
	}

	for _, theme := range themes {
		expectedValue, exists := expectedValues[theme]
		if !exists {
			t.Errorf("Theme %q not found in expected values", theme)
			continue
		}

		if string(theme) != expectedValue {
			t.Errorf("Theme constant mismatch: expected %q, got %q", expectedValue, theme)
		}
	}
}

func TestReferenceLayoutTypeConstants(t *testing.T) {
	layouts := []struct {
		layout   ReferenceLayoutType
		expected string
	}{
		{LayoutModern, "modern"},
		{LayoutClassic, "classic"},
	}

	for _, test := range layouts {
		if string(test.layout) != test.expected {
			t.Errorf("Layout constant mismatch: expected %q, got %q", test.expected, test.layout)
		}
	}
}

func TestCustomThemeCSS(t *testing.T) {
	t.Run("CustomThemeCSS constant is defined", func(t *testing.T) {
		if CustomThemeCSS == "" {
			t.Error("CustomThemeCSS should not be empty")
		}
	})

	t.Run("CustomThemeCSS contains required CSS variables", func(t *testing.T) {
		requiredVars := []string{
			"--scalar-color-1",
			"--scalar-color-2",
			"--scalar-color-3",
			"--scalar-background-1",
			"--scalar-background-2",
			"--scalar-background-3",
			".light-mode",
			".dark-mode",
		}

		for _, cssVar := range requiredVars {
			if !contains(CustomThemeCSS, cssVar) {
				t.Errorf("CustomThemeCSS should contain %q", cssVar)
			}
		}
	})
}

func TestDefaultCDN(t *testing.T) {
	expectedCDN := "https://cdn.jsdelivr.net/npm/@scalar/api-reference"

	if DefaultCDN != expectedCDN {
		t.Errorf("DefaultCDN mismatch: expected %q, got %q", expectedCDN, DefaultCDN)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
