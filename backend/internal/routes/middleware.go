package routes

import (
	"net/http"
)

// Middleware that adds CORS headers to every response
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers to allow cross-origin requests.
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with allowed origins.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests (OPTIONS method).
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler.
		next.ServeHTTP(w, r)
	})
}
