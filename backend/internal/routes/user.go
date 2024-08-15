package routes

import (
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/controller/core"
	chi "github.com/go-chi/chi/v5"
)

// Adds User routes to a Chi mux router (includes login, forgot password, etc)
func AddUserApiRoutes(router *chi.Mux, user core.UserController) *chi.Mux {
	// Public routes
	router.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes
		// Login
		mux.Post("/api/users/login", user.Login)
		// Forgot password
		mux.Post("/api/users/forgot-password", user.ResetPassword)
		// Verify Email
		mux.Get("/api/users/verify-email/{token}", user.EmailVerification)

		// Create new user
		mux.Post("/api/users", user.Create)

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(auth.AuthenticateJWT)

			// @tag.name Private routes
			// @tag.description Protected routes
			// users
			mux.Get("/api/users", user.FindAll)
			mux.Get("/api/users/{id}", user.Find)
			mux.Put("/api/users/{id}", user.Update)
			mux.Delete("/api/users/{id}", user.Delete)

			// My profile
			mux.Get("/api/me", user.GetMyUserDetails)
			mux.Post("/api/me", controller.HealthCheck)
			mux.Put("/api/me", user.UpdateMyProfile)

			// Email verification
			mux.Post("/api/users/send-verification-email", user.ResendVerificationEmail)

		})

	})
	return router
}
