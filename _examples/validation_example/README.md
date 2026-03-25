# Spec Validation Example

This example demonstrates the OpenAPI/Swagger specification validation feature in Scalar Go.

## Features Demonstrated

1. ✅ Enable validation in middleware
2. ✅ Disable validation for performance
3. ✅ Manual validation from file
4. ✅ Manual validation from content
5. ✅ Error handling and reporting

## Prerequisites

- Go 1.22 or higher
- Internet connection (for fetching dependencies)

## Running the Example

1. Navigate to the example directory:
```bash
cd _examples/validation_example
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run main.go
```

4. The server will start on `http://localhost:3000`

## Available Endpoints

### 📚 API Documentation (with validation)
- URL: http://localhost:3000/docs
- Description: Renders API documentation with spec validation enabled
- Validation: ✅ Enabled

### 📚 API Documentation (without validation)
- URL: http://localhost:3000/docs-no-validation
- Description: Renders API documentation without validation for better performance
- Validation: ❌ Disabled

### ✅ Validate Spec File
- URL: http://localhost:3000/validate
- Description: Validates the swagger.yaml file and returns validation result
- Method: GET
- Response Example:
```json
{
  "message": "Spec is valid!"
}
```

### ✅ Validate Spec Content
- URL: http://localhost:3000/validate-content
- Description: Validates inline OpenAPI spec content
- Method: GET
- Response Example:
```json
{
  "message": "Spec content is valid!"
}
```

## Testing Validation

### Valid Spec
The included `swagger.yaml` file is a valid OpenAPI 3.0 specification. Try accessing the endpoints to see successful validation.

### Invalid Spec
To test error handling, you can modify `swagger.yaml`:

1. Remove the `info.title` field:
```yaml
openapi: 3.0.0
info:
  # title: Sample API  <- Remove this
  version: 1.0.0
paths:
  /hello:
    get:
      responses:
        '200':
          description: OK
```

2. Restart the server and access http://localhost:3000/docs

3. You should see a validation error like:
```
spec validation failed: OpenAPI 3.x validation failed: field info.title is required
```

### Testing Different Versions

The validator supports:
- ✅ OpenAPI 3.1.x
- ✅ OpenAPI 3.0.x
- ✅ OpenAPI 2.0 (Swagger)

Try creating different spec files with different versions to test the validator.

## Code Examples

### Enable Validation in Middleware

```go
app.Get("/docs", scalarFiber.Handler(&scalar.Options{
    SpecURL:      "./swagger.yaml",
    ValidateSpec: true, // Enable validation
    CustomOptions: scalar.CustomOptions{
        PageTitle: "API Docs with Validation",
    },
}))
```

### Manual Validation

```go
// Validate from file
err := scalar.ValidateSpecFromFile("./swagger.yaml")
if err != nil {
    log.Printf("Validation failed: %v", err)
}

// Validate from content
err = scalar.ValidateSpec(specContent)
if err != nil {
    log.Printf("Validation failed: %v", err)
}
```

## Performance Notes

- Validation adds ~1-5ms overhead per request
- For production with static specs, validate once during startup
- For development, enable `ValidateSpec: true` for safety
- The validator caches parsed specs internally for better performance

## Learn More

- [OpenAPI Specification](https://swagger.io/specification/)
- [kin-openapi](https://github.com/getkin/kin-openapi) - The validation library used by Scalar Go
- [Scalar Documentation](https://github.com/scalar/scalar)
