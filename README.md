## Scalar Go

<img src="_examples/image.png" align="center">


## Overview 🌐

The Scalar package serves as a provider for the [Scalar](https://github.com/scalar/scalar) project. It offers a comprehensive suite of functions designed for generating API references in HTML format, specializing in JSON data handling and web presentation customization. This includes functions to serialize options into JSON, manage HTML escaping, and dynamically handle different types of specification content.

## Installation 📦

```bash
go get -u github.com/watchakorn-18k/scalar-go@latest
```

## Features 🚀

### JSON Serialization and HTML Escaping

- **safeJSONConfiguration** 🔒: Serializes configuration options into JSON format and escapes HTML characters to prevent XSS attacks. This ensures that the JSON can be safely embedded within HTML documents.

### Specification Content Handling

- **specContentHandler** 📝: Dynamically handles different types of specification content. It can process content as a function returning a map, a direct map, or a plain string, converting the specification content into JSON format suitable for web use.

### HTML Generation

- **ApiReferenceHTML** 📄: Generates a complete HTML document for API reference. It allows for extensive customization, including themes, layouts, and CDN configuration. It handles both direct specification content and content fetched from a URL, providing error handling for missing specifications.

## Customization Options ⚙️

The package allows extensive customization of the generated API reference through the `Options` struct, supporting:

- **CDN**: URL of the CDN to load additional scripts or styles.
- **Theme**: Customizable themes for styling the API reference.
- **Layout**: Choice between modern and classic layout designs.
- **SpecURL**: URL from which the specification content can be fetched.
- **Dark Mode**: Option to enable a dark theme for the API reference.

### Themes and Styles 🎨

- Provides a default set of CSS styles for both light and dark themes.
- Allows custom CSS injections to tailor the appearance to specific branding or aesthetic requirements.

## Error Handling 🛠️

- Robust error handling for scenarios where necessary parameters like `SpecURL` or `SpecContent` are missing.
- Errors during the fetching of specification content from URLs are properly managed and reported.

## Usage 📚

### Framework Middleware 🎯

Scalar Go provides ready-to-use middleware for popular Go web frameworks:

#### Fiber

```go
package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/watchakorn-18k/scalar-go"
	scalarfiber "github.com/watchakorn-18k/scalar-go/middleware/fiber"
)

func main() {
	app := fiber.New()

	app.Use("/docs", scalarfiber.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API Documentation",
		},
		DarkMode: true,
	}))

	log.Fatal(app.Listen(":3000"))
}
```

#### Gin

```go
package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/watchakorn-18k/scalar-go"
	scalargin "github.com/watchakorn-18k/scalar-go/middleware/gin"
)

func main() {
	r := gin.Default()

	r.GET("/docs", scalargin.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API Documentation",
		},
		DarkMode: true,
	}))

	log.Fatal(r.Run(":3000"))
}
```

#### Echo

```go
package main

import (
	"log"
	"github.com/labstack/echo/v4"
	"github.com/watchakorn-18k/scalar-go"
	scalarecho "github.com/watchakorn-18k/scalar-go/middleware/echo"
)

func main() {
	e := echo.New()

	e.GET("/docs", scalarecho.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API Documentation",
		},
		DarkMode: true,
	}))

	log.Fatal(e.Start(":3000"))
}
```

#### Chi

```go
package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/watchakorn-18k/scalar-go"
	scalarchi "github.com/watchakorn-18k/scalar-go/middleware/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/docs", scalarchi.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API Documentation",
		},
		DarkMode: true,
	}))

	log.Fatal(http.ListenAndServe(":3000", r))
}
```

### Manual Usage (without middleware)

You can also use the core `ApiReferenceHTML` function directly:

```go
package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/watchakorn-18k/scalar-go"
)

func main() {
	app := fiber.New()

	app.Use("/api/docs", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.yaml",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
			},
			DarkMode: true,
		})

		if err != nil {
			return err
		}
		c.Type("html")
		return c.SendString(htmlContent)
	})

	log.Fatal(app.Listen(":3000"))
}
```
## License

This project is forked from [MarceloPetrucio/go-scalar-api-reference](https://github.com/MarceloPetrucio/go-scalar-api-reference) and is licensed under the MIT License.