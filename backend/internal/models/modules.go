package models

import "gorm.io/gorm"

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
	NewAdminController func(interface{}, interface{}) interface{}
	// PolicySet is used to setup the different policies for the module
	// The policy set will set the policy for the non-admin CRUD portion of the API
	PolicySet ModulePolicySet
}

// ModulePolicySet is used to store the different policies for the module
type ModulePolicySet = []RolePolicySet

// RolePolicySet is used to store the different policies for the module
// eg. policySet["admin"] = []string{"create", "read", "update", "delete}
type RolePolicySet = map[string][]string


// Basic Paginated Response
type BasicPaginatedResponse[dbSchema any] struct {
	Data *[]dbSchema	`json:"data"`
	Meta SchemaMetaData 		`json:"meta"`
}
