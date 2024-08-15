package webapi

import (
	"net/http"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
)

type ActionService interface {
	// Record action in database
	RecordAction(r *http.Request, schemaName string, schemaID uint, recordAction *models.RecordedAction, changeObjects helpers.ChangeLogInput) error
	RecordBulkDelete(r *http.Request, schemaName, pluralSchemaName string, schemaIDs []int, recordAction *models.RecordedAction) error
	// CRUD operations
	FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[db.Action], error)
	FindById(int) (*db.Action, error)
	Create(action *models.CreateAction) (*db.Action, error)
	Update(int, *models.UpdateAction) (*db.Action, error)
	Delete(int) error
	BulkDelete([]int) error
}