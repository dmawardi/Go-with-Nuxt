package adminpanel

import (
	"html/template"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/models"
)

// init state for db access
var app *config.AppConfig

// Links for admin header
var header = HeaderSection{
	HomeUrl:           "#",
	ViewSiteUrl:       "#",
	LogOutUrl:         "#",
	ChangePasswordUrl: "#",
}

// Function called in main.go to connect app state to current file
func SetStateInAdminPanel(a *config.AppConfig) {
	app = a
	// Set header urls after setting state
	header.HomeUrl = template.URL("http://" + app.BaseURL + "/admin/home")
	header.ChangePasswordUrl = template.URL("http://" + app.BaseURL + "/admin/change-password")
	header.LogOutUrl = template.URL("http://" + app.BaseURL + "/admin/logout")
	header.ViewSiteUrl = template.URL("http://" + app.BaseURL + "/swagger/index.html")
}

// Admin controller (used in API)
type AdminPanelController struct {
	// Basic modules
	Base AdminCoreController
	User AdminUserController
	Auth AdminAuthPolicyController
	Action AdminActionController
	// Additional modules contained in module map
	ModuleMap models.ModuleMap
}

// Constructor
func NewAdminPanelController(
							base AdminCoreController, 
							users AdminUserController, 
							authPolicies AdminAuthPolicyController, 
							action AdminActionController,
							moduleMap models.ModuleMap) AdminPanelController {
	return AdminPanelController{base, users, authPolicies, action, moduleMap}
}



