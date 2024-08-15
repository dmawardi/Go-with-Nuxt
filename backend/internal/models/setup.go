package models

// Wrappers used for API setup function to ensure different repositories are standardized
// type RepoFactory func(*gorm.DB) interface{}
// type ServiceFactory func(interface{}) interface{}
// type ControllerFactory func(interface{}) interface{}

// Module set is used within SetupBasicModules to store the different modules
type ModuleSet = struct {
	RouteName       string
	Repo            interface{}
	Service         interface{}
	Controller      interface{}
	AdminController interface{}
}

// ModuleMap is used to store the different modules in a map for dynamic usage
// Outputted by SetupModules() in webApi/setup.go
type ModuleMap = map[string]ModuleSet
