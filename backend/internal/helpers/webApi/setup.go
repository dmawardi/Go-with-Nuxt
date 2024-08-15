package webapi

import (
	"bytes"
	"html/template"

	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

// MODULE STANDARDIZATION
//

// Helper function to create a new repository. Takes a repository creation function and returns a function that takes a database connection and returns an interface
func NewRepository[T any](repoFunc func(*gorm.DB) T) func(*gorm.DB) interface{} {
	return func(db *gorm.DB) interface{} {
		return repoFunc(db)
	}
}

// Helper function to create a new service. Takes a service creation function and returns a function that takes an interface and returns an interface
func NewService[T any, S any](serviceFunc func(T) S) func(interface{}) interface{} {
	return func(repoInterface interface{}) interface{} {
		repo, ok := repoInterface.(T)
		if !ok {
			panic("Incorrect repository type")
		}
		return serviceFunc(repo)
	}
}

// Helper function to create a new controller. Takes a controller creation function and returns a function that takes an interface and returns an interface
func NewController[T any, C any](controllerFunc func(T) C) func(interface{}) interface{} {
	return func(serviceInterface interface{}) interface{} {
		service, ok := serviceInterface.(T)
		if !ok {
			panic("Incorrect service type")
		}
		return controllerFunc(service)
	}
}

// Helper function that takes an admin controller creation function and returns a function that takes an interface and returns a basic admin controller
func NewAdminController[Serv any, AS ActionService](controllerFunc func(Serv, AS) models.BasicAdminController) func(interface{}, ActionService) models.BasicAdminController {
	return func(serviceInterface interface{}, actionserviceInterface ActionService) models.BasicAdminController {
		service, ok := serviceInterface.(Serv)
		if !ok {
			panic("Incorrect service type")
		}
		actionServ, ok := actionserviceInterface.(AS)
		if !ok {
			panic("Incorrect action service type")
		}
		return controllerFunc(service, actionServ)
	}
}

// LoadTemplate parses an HTML template, executes it with the provided data, and returns the result as a string.
func LoadTemplate(templateFilePath string, data interface{}) (string, error) {
	// Parse the template file
	t, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}

	// Build the template with the injected data
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
