package routes

import (
	"fmt"

	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/go-chi/chi/v5"
)

// Adds a basic Admin CRUD route set to a Chi mux router
func AddAdminRouteSet(router *chi.Mux, protected bool, urlExtension string, controller models.BasicAdminController) *chi.Mux {
	// Reassign for consistency
	r := router
	r.Group(func(mux chi.Router) {
		// Set to use JWT authentication if protected
		if protected {
			mux.Use(auth.AuthenticateJWT)
		}
		// Read (all users)
		mux.Get(fmt.Sprintf("/admin/%s", urlExtension), controller.FindAll)
		// Create (GET form / POST form)
		mux.Get(fmt.Sprintf("/admin/%s/create", urlExtension), controller.Create)
		mux.Post(fmt.Sprintf("/admin/%s/create", urlExtension), controller.Create)
		mux.Get(fmt.Sprintf("/admin/%s/create/success", urlExtension), controller.CreateSuccess)
		// Delete
		mux.Get(fmt.Sprintf("/admin/%s/delete/{id}", urlExtension), controller.Delete)
		mux.Post(fmt.Sprintf("/admin/%s/delete/{id}", urlExtension), controller.Delete)
		mux.Get(fmt.Sprintf("/admin/%s/delete/success", urlExtension), controller.DeleteSuccess)
		// Bulk delete (from table)
		mux.Delete(fmt.Sprintf("/admin/%s/bulk-delete", urlExtension), controller.BulkDelete)

		// Edit/Update (GET data in form / POST form)
		mux.Get(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.Edit)
		mux.Post(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.Edit)
		mux.Get(fmt.Sprintf("/admin/%s/edit/success", urlExtension), controller.EditSuccess)
	})
	return router
}

func AddAdminActionRouteSet(router *chi.Mux, protected bool, urlExtension string, controller adminpanel.AdminActionController) *chi.Mux {
	// Reassign for consistency
	r := router
	r.Group(func(mux chi.Router) {
		// Set to use JWT authentication if protected
		if protected {
			mux.Use(auth.AuthenticateJWT)
		}
		// Read All
		mux.Get(fmt.Sprintf("/admin/%s", urlExtension), controller.FindAll)

		// View One
		mux.Get(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.View)
	})
	return router
}

// Adds routess for editing and creating admin auth policies for the admin panel
func AddAdminPolicySet(router *chi.Mux, protected bool, urlExtension string, controller adminpanel.AdminAuthPolicyController) *chi.Mux {
	// Reassign for consistency
	r := router
	r.Group(func(mux chi.Router) {
		// Set to use JWT authentication if protected
		if protected {
			mux.Use(auth.AuthenticateJWT)
		}
		// Read (all users)
		mux.Get(fmt.Sprintf("/admin/%s", urlExtension), controller.FindAll)
		mux.Get(fmt.Sprintf("/admin/%s/roles", urlExtension), controller.FindAllRoles)
		mux.Get(fmt.Sprintf("/admin/%s/inheritance", urlExtension), controller.FindAllRoleInheritance)
		// Create Policy (GET form / POST form)
		mux.Get(fmt.Sprintf("/admin/%s/create", urlExtension), controller.Create)
		mux.Post(fmt.Sprintf("/admin/%s/create", urlExtension), controller.Create)
		mux.Get(fmt.Sprintf("/admin/%s/create/success", urlExtension), controller.CreateSuccess)
		// Create Role
		mux.Get(fmt.Sprintf("/admin/%s/create-role", urlExtension), controller.CreateRole)
		mux.Post(fmt.Sprintf("/admin/%s/create-role", urlExtension), controller.CreateRole)
		mux.Get(fmt.Sprintf("/admin/%s/create-role/success", urlExtension), controller.CreateRoleSuccess)
		// Create inheritance
		mux.Get(fmt.Sprintf("/admin/%s/create-inheritance", urlExtension), controller.CreateInheritance)
		mux.Post(fmt.Sprintf("/admin/%s/create-inheritance", urlExtension), controller.CreateInheritance)
		mux.Get(fmt.Sprintf("/admin/%s/create-inheritance/success", urlExtension), controller.CreateInheritanceSuccess)
		// // Delete Inheritance
		mux.Get(fmt.Sprintf("/admin/%s/delete-inheritance/{inherit-slug}", urlExtension), controller.DeleteInheritance)
		mux.Post(fmt.Sprintf("/admin/%s/delete-inheritance/{inherit-slug}", urlExtension), controller.DeleteInheritance)
		mux.Get(fmt.Sprintf("/admin/%s/delete-inheritance/success", urlExtension), controller.DeleteInheritanceSuccess)

		// Edit/Update (GET data in form / POST form)
		mux.Get(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.Edit)
		mux.Post(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.Edit)
		mux.Delete(fmt.Sprintf("/admin/%s/{id}", urlExtension), controller.Edit)
	})
	return router
}

// Function to adds the base admin panel routes to an existing Chi mux router (eg. login, logout, home)
func AddBasicAdminRoutes(router *chi.Mux, controller adminpanel.AdminCoreController) *chi.Mux {
	// Public routes
	router.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes
		mux.Get("/admin", controller.AdminRedirectBasedOnLoginStatus)
		mux.Get("/admin/login", controller.Login)
		mux.Post("/admin/login", controller.Login)

		// admin logout
		mux.Get("/admin/logout", controller.Logout)

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(auth.AuthenticateJWT)

			// @tag.name Private routes
			// @tag.description Protected routes
			// admin home
			mux.Get("/admin/home", controller.Home)
			// Change password
			mux.Get("/admin/change-password", controller.ChangePassword)
			mux.Post("/admin/change-password", controller.ChangePassword)

			mux.Get("/admin/change-password-success", controller.ChangePasswordSuccess)

		})

	})

	return router
}
