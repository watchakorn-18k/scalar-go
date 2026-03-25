package chi

import (
	"net/http"

	"github.com/watchakorn-18k/scalar-go"
)

// Handler creates a Chi/net/http handler for Scalar API documentation
//
// Example usage:
//
//	r := chi.NewRouter()
//	r.Get("/docs", chi.Handler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//		CustomOptions: scalar.CustomOptions{
//			PageTitle: "My API Documentation",
//		},
//	}))
//
// With UI authentication:
//
//	r.Get("/docs", chi.Handler(&scalar.Options{
//		SpecURL: "./swagger.yaml",
//		UIUsername: "admin",
//		UIPassword: "secret",
//	}))
func Handler(options *scalar.Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if UI authentication is enabled
		if options.IsUIAuthEnabled() {
			auth := r.Header.Get("Authorization")
			username, password, ok := scalar.ParseBasicAuth(auth)

			if !ok || !options.ValidateUICredentials(username, password) {
				w.Header().Set("WWW-Authenticate", `Basic realm="Scalar API Documentation"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		htmlContent, err := scalar.ApiReferenceHTML(options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlContent))
	}
}
