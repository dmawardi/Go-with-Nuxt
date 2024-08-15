package repository

import (
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/models"
	corerepositories "github.com/dmawardi/Go-Template/internal/repository/core"
)

var app *config.AppConfig

func SetAppConfig(appConfig *config.AppConfig) {
	app = appConfig
	// Set app config in repository
	corerepositories.SetAppConfig(app)
}

type BasicModuleRepository[dbSchema any] interface {
	FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[dbSchema], error)
	FindById(int) (*dbSchema, error)
	Create(entity *dbSchema) (*dbSchema, error)
	Update(int, *dbSchema) (*dbSchema, error)
	Delete(int) error
	BulkDelete([]int) error
}