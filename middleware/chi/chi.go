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
func Handler(options *scalar.Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
