package corerepositories

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers/data"
	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)


type ActionRepository interface {
	// Find a list of all users in the Database
	FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[db.Action], error)
	FindById(int) (*db.Action, error)
	Create(action *db.Action) (*db.Action, error)
	Update(int, *db.Action) (*db.Action, error)
	Delete(int) error
	BulkDelete([]int) error
}

type actionRepository struct {
	DB *gorm.DB
}

func NewActionRepository(db *gorm.DB) ActionRepository {
	return &actionRepository{db}
}

// Creates a action in the database
func (r *actionRepository) Create(action *db.Action) (*db.Action, error) {
	// Create above action in database
	result := r.DB.Create(&action)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating action: %w", result.Error)
	}

	return action, nil
}

// Find a list of actions in the database
func (r *actionRepository) FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[db.Action], error) {
	// Build meta data for actions
	metaData, err := data.BuildMetaData(r.DB, db.Action{}, limit, offset, order, conditions)
	if err != nil {
		fmt.Printf("Error building meta data: %s", err)
		return nil, err
	}

	// Query all actions based on the received parameters
	var actions []db.Action
	err = data.QueryAll(r.DB, &actions, limit, offset, order, conditions, []string{"Admin"})
	if err != nil {
		fmt.Printf("Error querying db for list of actions: %s", err)
		return nil, err
	}

	return &models.BasicPaginatedResponse[db.Action]{
		Data: &actions,
		Meta: *metaData,
	}, nil
}

// Find action in database by ID
func (r *actionRepository) FindById(id int) (*db.Action, error) {
	// Create an empty ref object of type action
	action := db.Action{}
	// Check if action exists in db
	result := r.DB.First(&action, id)
	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &action, nil
}

// Delete action in database
func (r *actionRepository) Delete(id int) error {
	// Create an empty ref object of type action
	action := db.Action{}
	// Check if action exists in db
	result := r.DB.Delete(&action, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting action: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Bulk delete actions in database
func (r *actionRepository) BulkDelete(ids []int) error {
	// Delete users with specified IDs
	err := data.BulkDeleteByIds(db.Action{}, ids, r.DB)
	if err != nil {
		fmt.Println("error in deleting actions: ", err)
		return err
	}
	// else
	return nil
}

// Updates action in database
func (r *actionRepository) Update(id int, action *db.Action) (*db.Action, error) {
	// Init
	var err error
	// Find action by id
	found, err := r.FindById(id)
	if err != nil {
		fmt.Println("Action to update not found: ", err)
		return nil, err
	}
	// Set action user id (gorm requires this as it does not automatically set the foreign key)
	if action.AdminID == 0 {
		action.AdminID = action.Admin.ID
	}

	// Update found action
	updateResult := r.DB.Model(&found).Updates(action)
	if updateResult.Error != nil {
		fmt.Println("Action update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Retrieve changed action by id
	updated, err := r.FindById(id)
	if err != nil {
		fmt.Println("Action to update not found: ", err)
		return nil, err
	}
	return updated, nil
}