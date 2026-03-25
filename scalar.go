package scalar

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

func safeJSONConfiguration(options *Options) string {
	// Serializes the options to JSON
	jsonData, _ := json.Marshal(options)
	// Escapes double quotes into HTML entities
	escapedJSON := strings.ReplaceAll(string(jsonData), `"`, `&quot;`)
	return escapedJSON
}

func specContentHandler(specContent interface{}) string {
	switch spec := specContent.(type) {
	case func() map[string]interface{}:
		// If specContent is a function, it calls the function and serializes the return
		result := spec()
		jsonData, _ := json.Marshal(result)
		return string(jsonData)
	case map[string]interface{}:
		// If specContent is a map, it serializes it directly
		jsonData, _ := json.Marshal(spec)
		return string(jsonData)
	case string:
		// If it is a string, it returns directly
		return spec
	default:
		// Otherwise, returns empty
		return ""
	}
}

func ApiReferenceHTML(optionsInput *Options) (string, error) {
	options := DefaultOptions(*optionsInput)

	if options.SpecURL == "" && options.SpecContent == nil {
		return "", fmt.Errorf("specURL or specContent must be provided")
	}

	var specContentStr string

	if options.SpecContent == nil && options.SpecURL != "" {

		if strings.HasPrefix(options.SpecURL, "http") {
			content, err := fetchContentFromURL(options.SpecURL)
			if err != nil {
				return "", err
			}
			specContentStr = content
			options.SpecContent = content
		} else {
			cleanPath := filepath.Clean(optionsInput.SpecURL)
			absPath, err := filepath.Abs(cleanPath)
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

			specContentStr = string(content)
			options.SpecContent = specContentStr
		}
	} else if options.SpecContent != nil {
		// Convert SpecContent to string for validation
		switch spec := options.SpecContent.(type) {
		case string:
			specContentStr = spec
		case func() map[string]interface{}:
			result := spec()
			jsonData, _ := json.Marshal(result)
			specContentStr = string(jsonData)
		case map[string]interface{}:
			jsonData, _ := json.Marshal(spec)
			specContentStr = string(jsonData)
		}
	}

	// Validate spec if validation is enabled
	if options.ValidateSpec && specContentStr != "" {
		if err := ValidateSpec(specContentStr); err != nil {
			return "", fmt.Errorf("spec validation failed: %w", err)
		}
	}

	dataConfig := safeJSONConfiguration(options)
	specContentHTML := specContentHandler(options.SpecContent)

	var pageTitle string

	if options.CustomOptions.PageTitle != "" {
		pageTitle = options.CustomOptions.PageTitle
	} else {
		pageTitle = "Scalar API Reference"
	}

	customThemeCss := CustomThemeCSS

	if options.Theme != "" {
		customThemeCss = ""
	}

	// Build favicon link tag if FaviconURL is provided
	var faviconTag string
	if options.CustomOptions.FaviconURL != "" {
		faviconTag = fmt.Sprintf(`<link rel="icon" href="%s" />`, options.CustomOptions.FaviconURL)
	}

	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
      <head>
        <title>%s</title>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />%s
        <style>%s</style>
      </head>
      <body>
        <script id="api-reference" type="application/json" data-configuration="%s">%s</script>
        <script src="%s"></script>
      </body>
    </html>
  `, pageTitle, faviconTag, customThemeCss, dataConfig, specContentHTML, options.CDN), nil
}
