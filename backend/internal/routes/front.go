package routes

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/go-chi/chi/v5"
)

// Serve Front end Vue.js SPA
func ServeFrontEnd(router *chi.Mux, protected bool) *chi.Mux {
	// Reassign for consistency
	r := router

	return reverseProxy(r, protected, 3000)
}

func reverseProxy(router *chi.Mux, protected bool, portNumber int) *chi.Mux {
	// Parse the target URL (ensure it parses correctly)
	targetURL := fmt.Sprintf("http://localhost:%d", portNumber)
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}
	// Create a reverse proxy that will forward requests to the Nuxt.js server
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Ensure request rewriting is handled safely
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		// Call the original director
		originalDirector(req)
		// Ensure the URL is correct
		req.Host = target.Host
	}

	router.Group(func(mux chi.Router) {
		// Set to use JWT authentication if protected
		if protected {
			mux.Use(auth.AuthenticateJWT)
		}
		// Handle all requests by proxying to the Nuxt.js server
		mux.Handle("/*", proxy)
	})

	return router
}

// serveStaticAssets handles serving static assets like images, CSS, and JavaScript files
// func serveStaticAssets(router *chi.Mux) *chi.Mux {
// 	// Serve the assets from the Nuxt.js static directory
// 	staticHandler := http.StripPrefix("/_nuxt/", http.FileServer(http.Dir("../frontend/.nuxt/dist/client")))

// 	// In development mode, let Vite serve the assets dynamically
// 	router.Handle("/_nuxt/*", staticHandler)

// 	// If you have other static assets in a different directory, handle them similarly
// 	assetsHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("../frontend/static")))
// 	router.Handle("/static/*", assetsHandler)

// 	return router
// }
