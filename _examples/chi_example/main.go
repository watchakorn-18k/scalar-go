package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/watchakorn-18k/scalar-go"
	scalarchi "github.com/watchakorn-18k/scalar-go/middleware/chi"
)

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Scalar API Documentation
	r.Get("/docs", scalarchi.Handler(&scalar.Options{
		SpecURL: "./swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Chi API Documentation",
		},
		DarkMode: true,
	}))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Hello from Chi!",
			})
		})
	})

	log.Println("Server running on http://localhost:3000")
	log.Println("API Docs available at http://localhost:3000/docs")
	log.Fatal(http.ListenAndServe(":3000", r))
}
