package routes

import (
	"net/http"

	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller/core"
	"github.com/dmawardi/Go-Template/internal/models"
)

// Create new service repository
func BuildRouteState(a *config.AppConfig) {
	app = a
}

type Api interface {
	// Total route builder for API
	Routes() http.Handler
}

// Api that contains all controllers for route creation
type api struct {
	// Basic Controllers
	User   core.UserController
	Policy core.AuthPolicyController
	// Admin Controller
	Admin adminpanel.AdminPanelController
	// Module Controllers
	ModuleMap models.ModuleMap
	// Post   modulecontrollers.PostController

}

func NewApi(
	admin adminpanel.AdminPanelController,
	user core.UserController,
	policy core.AuthPolicyController,
	moduleMap models.ModuleMap) Api {
	return &api{Admin: admin, User: user, Policy: policy, ModuleMap: moduleMap}
}
