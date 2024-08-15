package modules

import (
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

// Used to setup modules aside from core modules (User, Policy)
// Returns module map that contains structs with modules (controller, service, & repo) using module name as key
func SetupModules(modulesToSetup []EntityConfig, client *gorm.DB, actionService webapi.ActionService) models.ModuleMap {
	// Init
	moduleMap := make(map[string]models.ModuleSet)

	for _, module := range modulesToSetup {
		// Create repo, service, and controller using the client
		repo := module.NewRepo(client)
		service := module.NewService(repo)
		controller := module.NewController(service)

		// Assign constructor function to newAdminController
		newAdminController := module.NewAdminController
		// If admin controller constructor is not nil, add it to the module map
		if newAdminController != nil {

			// Create admin controller using the service
			adminController := newAdminController(service, actionService)

			// Add module set including admin controller to the map
			moduleMap[module.Name] = models.ModuleSet{
				RouteName:       module.RouteName,
				Repo:            repo,
				Service:         service,
				Controller:      controller,
				AdminController: adminController,
			}
		} else {
			// Add module set without admin controller to the map
			moduleMap[module.Name] = models.ModuleSet{
				RouteName:       module.RouteName,
				Repo:            repo,
				Service:         service,
				Controller:      controller,
				AdminController: nil,
			}
		}
	}
	return moduleMap
}

// Used in API setup to standardize the array of setup configurations
type EntityConfig struct {
	// Name is for the module map key
	Name string
	// RouteName is for the API route
	RouteName string
	// NewRepo, NewService, NewController, and NewAdminController are used to create the different parts of the module
	NewRepo            func(*gorm.DB) interface{}
	NewService         func(interface{}) interface{}
	NewController      func(interface{}) interface{}
	NewAdminController func(interface{}, webapi.ActionService) models.BasicAdminController
	// PolicySet is used to setup the different policies for the module
	// The policy set will set the policy for the non-admin CRUD portion of the API
	PolicySet ModulePolicySet
}

// ModulePolicySet is used to store the different policies for the module
type ModulePolicySet = []RolePolicySet

// RolePolicySet is used to store the different policies for the module
// eg. policySet["admin"] = []string{"create", "read", "update", "delete}
type RolePolicySet = map[string][]string
