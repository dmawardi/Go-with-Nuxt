package service

import (
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/models"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
	moduleservices "github.com/dmawardi/Go-Template/internal/service/module"
)

// Repository used by handler package
var app *config.AppConfig

// Create new service repository
func SetAppConfig(a *config.AppConfig) {
	// Set app state in core services
	coreservices.SetAppConfig(a)
	// Set app state in module services
	moduleservices.SetAppConfig(a)

	app = a
}

type BasicModuleService[dbSchema, create, update any] interface {
	FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[dbSchema], error)
	FindById(int) (*dbSchema, error)
	Create(entity *create) (*dbSchema, error)
	Update(int, *update) (*dbSchema, error)
	Delete(int) error
	BulkDelete([]int) error
}