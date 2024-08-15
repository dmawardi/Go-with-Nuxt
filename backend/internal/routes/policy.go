package routes

import (
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/controller/core"
	"github.com/go-chi/chi/v5"
)

// Adds Authorization routes to a Chi mux router
func AddAuthRBACApiRoutes(router *chi.Mux, policy core.AuthPolicyController) *chi.Mux {
	// Public routes
	router.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(auth.AuthenticateJWT)

			// @tag.name Private routes
			// @tag.description Protected routes
			// AUTH
			//
			// Policies
			mux.Get("/api/auth", policy.FindAll)
			mux.Get("/api/auth/{policy-slug}", policy.FindByResource)
			mux.Post("/api/auth", policy.Create)
			mux.Put("/api/auth", policy.Update)
			mux.Delete("/api/auth", policy.Delete)
			// Roles
			mux.Get("/api/auth/roles", policy.FindAllRoles)
			mux.Put("/api/auth/roles", policy.AssignUserRole)
			mux.Post("/api/auth/roles", policy.CreateRole)
			// Inheritance
			mux.Get("/api/auth/inheritance", policy.FindAllRoleInheritance)
			mux.Post("/api/auth/inheritance", policy.CreateInheritance)
			mux.Delete("/api/auth/inheritance", policy.DeleteInheritance)
		})

	})
	return router
}
