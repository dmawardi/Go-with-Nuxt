package routes

import (
	"net/http"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/go-chi/chi/v5"
)

// Serve Front end Vue.js SPA
func ServeFrontEnd(router *chi.Mux, protected bool) *chi.Mux {
	// Reassign for consistency
	r := router

	// Build fileserver using static directory
	fileServer := http.FileServer(http.Dir("../frontend/dist"))
	// Handle all calls to /static/* by stripping prefix and sending to file server
	r.Handle("/front/*", http.StripPrefix("/front/", fileServer))

	r.Group(func(mux chi.Router) {
		// Set to use JWT authentication if protected
		if protected {
			mux.Use(auth.AuthenticateJWT)
		}
		// Read All
		mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../frontend/dist/index.html")
		})
	})

	return router
}
