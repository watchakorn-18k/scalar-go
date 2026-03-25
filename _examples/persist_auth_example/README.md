# 🔑 Auth Persistence Example

This example demonstrates the **authentication persistence** feature in Scalar Go. When enabled, bearer tokens and other authentication credentials are automatically saved to browser localStorage and restored when the page is refreshed.

## ✨ Features

- 💾 **Auto-save tokens** - Bearer tokens saved to localStorage
- 🔄 **Auto-restore** - Tokens restored on page refresh
- ⏰ **Auto-expire** - Saved tokens expire after 24 hours
- 🗑️ **Manual clear** - Clear tokens with `scalarClearAuth()`
- 🔒 **Secure** - Only stores tokens in browser localStorage (client-side only)

## 🚀 Quick Start

### 1. Run the server

```bash
cd _examples/persist_auth_example
go mod tidy
go run main.go
```

### 2. Open the docs

Visit: http://localhost:3000/docs

### 3. Try it out!

1. Click on the **GET /api/protected** endpoint
2. Click "Auth" or add Authorization header
3. Enter bearer token: `test-token-123`
4. Send the request (you should get success)
5. **Refresh the page** 🔄
6. Your token is still there! ✨

## 📊 Comparison

### ✅ With Persistence (`/docs`)

```go
app.Get("/docs", scalarFiber.Handler(&scalar.Options{
    SpecURL:     "./swagger.yaml",
    PersistAuth: true, // 🔑 Token saved!
}))
```

- ✅ Token saved to localStorage
- ✅ Token restored on refresh
- ✅ Expires after 24 hours
- ✅ Can be manually cleared

### ❌ Without Persistence (`/docs-no-persist`)

```go
app.Get("/docs-no-persist", scalarFiber.Handler(&scalar.Options{
    SpecURL:     "./swagger.yaml",
    PersistAuth: false, // ❌ Token lost on refresh
}))
```

- ❌ Token lost on refresh
- ❌ Must re-enter token every time
- ❌ No persistence

## 🧪 Testing

### Test the protected endpoint

```bash
# Without token (should fail)
curl http://localhost:3000/api/protected

# With token (should succeed)
curl -H "Authorization: Bearer test-token-123" \
     http://localhost:3000/api/protected
```

### Test the public endpoint

```bash
# No auth required
curl http://localhost:3000/api/public
```

## 💡 How It Works

1. **JavaScript Injection**: When `PersistAuth: true`, Scalar Go injects a script that:
   - Monitors auth input fields using MutationObserver
   - Saves values to localStorage on change
   - Restores values on page load
   - Auto-expires after 24 hours

2. **localStorage Keys**:
   - `scalar_auth_tokens` - Stored auth data
   - `scalar_auth_expiry` - Expiration timestamp

3. **Security**:
   - Data only stored client-side in browser
   - Never sent to server
   - Automatically expires after 24 hours
   - User can clear anytime with `scalarClearAuth()`

## 🔧 Configuration

### Basic Usage

```go
scalar.Options{
    PersistAuth: true, // Enable auth persistence
}
```

### Combined with Other Features

```go
scalar.Options{
    SpecURL:      "./swagger.yaml",
    PersistAuth:  true,           // Save auth tokens
    ValidateSpec: true,            // Validate spec
    UIUsername:   "admin",         // Protect UI
    UIPassword:   "secret",
    DarkMode:     true,
    CustomOptions: scalar.CustomOptions{
        PageTitle:  "My API",
        FaviconURL: "/favicon.ico",
    },
}
```

## 🎯 Use Cases

### ✅ When to Enable

- **Development environments** - Don't lose tokens during testing
- **Staging environments** - Easier testing workflow
- **Demo/sandbox APIs** - Better user experience
- **Internal tools** - Convenience for team members

### ⚠️ When to Consider Disabling

- **Shared computers** - Security concern if multiple users
- **Public documentation** - If you don't want browsers to cache tokens
- **High-security apps** - If you have strict security requirements

## 🔒 Security Notes

1. **Client-side only** - Tokens are stored in browser localStorage, not sent to server
2. **HTTPS recommended** - Use HTTPS in production to prevent token interception
3. **Auto-expiry** - Tokens automatically expire after 24 hours
4. **Manual clear** - Users can clear tokens with `scalarClearAuth()`
5. **Private browsing** - Won't work in incognito/private browsing mode

## 🐛 Debugging

### Check saved tokens

Open browser console:

```javascript
// View saved auth data
localStorage.getItem('scalar_auth_tokens')

// View expiry time
localStorage.getItem('scalar_auth_expiry')

// Clear saved tokens
scalarClearAuth()
```

### Common Issues

**Q: Tokens not persisting?**
- Check browser console for errors
- Make sure `PersistAuth: true` is set
- Verify localStorage is not disabled
- Check if private browsing mode

**Q: Tokens expired?**
- Tokens auto-expire after 24 hours
- Clear and re-enter token

**Q: Want to change expiry time?**
- Modify `DEFAULT_EXPIRY_HOURS` in `persist_auth.go`

## 📚 Related Features

- [Spec Validation](../validation_example/) - Validate OpenAPI specs
- [UI Authentication](../ui_auth_example/) - Protect docs with password
- [Authentication Config](../auth_example/) - Configure API auth methods

## 📖 Documentation

Full documentation: https://github.com/watchakorn-18k/scalar-go

---

Made with ❤️ using [Scalar Go](https://github.com/watchakorn-18k/scalar-go)
