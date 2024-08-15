package modules

import (
	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	modulecontrollers "github.com/dmawardi/Go-Template/internal/controller/moduleControllers"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	modulerepositories "github.com/dmawardi/Go-Template/internal/repository/module"
	moduleservices "github.com/dmawardi/Go-Template/internal/service/module"
)

// Define setup configurations (to use in setupModules within API setup function)
var ModulesToSetup = []EntityConfig{
	{
		// Used for module name in module map
		Name: "Post",
		// Route name for API
		RouteName: "posts",
		// Uses webapi to create new generic repository, service, controller, and admin controller
		NewRepo:            webapi.NewRepository(modulerepositories.NewPostRepository),
		NewService:         webapi.NewService(moduleservices.NewPostService),
		NewController:      webapi.NewController(modulecontrollers.NewPostController),
		NewAdminController: webapi.NewAdminController(adminpanel.NewAdminPostController),
	},
	// ADD ADDITIONAL BASIC MODULES HERE
}
