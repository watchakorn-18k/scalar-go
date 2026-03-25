# Branding Example

This example demonstrates how to customize your Scalar API documentation with custom logos and favicons.

## Features

- Custom logo and favicon from URL
- Custom logo and favicon from base64 data URI
- Logo only customization
- Favicon only customization
- Default (no branding)

## Running the Example

```bash
go run main.go
```

Then visit:
- http://localhost:3000/docs/url - URL-based assets
- http://localhost:3000/docs/base64 - Base64-encoded assets
- http://localhost:3000/docs/logo-only - Custom logo only
- http://localhost:3000/docs/favicon-only - Custom favicon only
- http://localhost:3000/docs/default - Default branding

## Supported Formats

### Logo
- PNG, JPG, SVG, GIF (via URL)
- Base64 data URI
- Remote URLs

### Favicon
- .ico, .png, .svg (via URL)
- Base64 data URI
- Remote URLs
- SVG emoji (e.g., `data:image/svg+xml,<svg>...</svg>`)

## Usage Examples

### Using URLs
```go
app.Use("/docs", scalarfiber.Handler(&scalar.Options{
    SpecURL: "./swagger.yaml",
    CustomOptions: scalar.CustomOptions{
        LogoURL:    "/assets/logo.png",
        FaviconURL: "/assets/favicon.ico",
    },
}))
```

### Using Base64
```go
app.Use("/docs", scalarfiber.Handler(&scalar.Options{
    SpecURL: "./swagger.yaml",
    CustomOptions: scalar.CustomOptions{
        LogoURL:    "data:image/svg+xml;base64,PHN2Zy...",
        FaviconURL: "data:image/png;base64,iVBORw0KG...",
    },
}))
```

### Using SVG Emoji Favicon
```go
app.Use("/docs", scalarfiber.Handler(&scalar.Options{
    SpecURL: "./swagger.yaml",
    CustomOptions: scalar.CustomOptions{
        FaviconURL: "data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>🚀</text></svg>",
    },
}))
```
