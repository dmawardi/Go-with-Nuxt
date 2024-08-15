package routes

import (
	"fmt"
	"net/http"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

var app *config.AppConfig

// Overall Routes builder for server
func (a api) Routes() http.Handler {
	// Create new router
	mux := chi.NewRouter()
	// Use built in Chi middleware
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(corsMiddleware)

	// Add user and group API routes
	mux = AddUserApiRoutes(mux, a.User)
	mux = AddAuthRBACApiRoutes(mux, a.Policy)

	// Add basic admin panel routes (home, login, etc)
	mux = AddBasicAdminRoutes(mux, a.Admin.Base)
	// Add admin user routes
	mux = AddAdminRouteSet(mux, true, "users", a.Admin.User)
	// Add admin policy routes
	mux = AddAdminPolicySet(mux, true, "policy", a.Admin.Auth)
	// Add admin action routes
	mux = AddAdminActionRouteSet(mux, true, "actions", a.Admin.Action)

	// Other schemas
	for _, module := range a.ModuleMap {
		// Add basic CRUD API routes
		mux = AddBasicCrudApiRoutes(mux, module.RouteName, module.Controller.(models.BasicController))
		// Add admin panel schema route sets
		mux = AddAdminRouteSet(mux, false, module.RouteName, module.AdminController.(models.BasicAdminController))
	}

	// Serve API Swagger docs at built URL from config state
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/static/docs/swagger.json", app.BaseURL)), //The url pointing to API definition
	))
	fmt.Printf("Serving Swagger docs at http://%s/swagger/index.html\n", app.BaseURL)

	// Build fileserver using static directory
	fileServer := http.FileServer(http.Dir("./static"))
	// Handle all calls to /static/* by stripping prefix and sending to file server
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
