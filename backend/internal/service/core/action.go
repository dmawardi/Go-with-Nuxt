package coreservices

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	"github.com/dmawardi/Go-Template/internal/models"
	corerepositories "github.com/dmawardi/Go-Template/internal/repository/core"
)



type actionService struct {
	repo corerepositories.ActionRepository
}

func NewActionService(repo corerepositories.ActionRepository) webapi.ActionService {
	return &actionService{repo: repo}
}
// Record action in database
func (s *actionService) RecordAction(r *http.Request, schemaName string, schemaID uint, recordAction *models.RecordedAction, changeObjects helpers.ChangeLogInput) error {
	// Generate change log
	changes, err := helpers.GenerateChangeLog(recordAction.ActionType, changeObjects)
	if err != nil {
		fmt.Println("Error generating change log: ", err)
		return err
	}
	// Generate change description
	changeDescription, err := helpers.GenerateChangeDescription(changes, schemaName, recordAction.ActionType, schemaID)
	if err != nil {
		fmt.Println("Error generating change description: ", err)
		return err
	}
	// Validate and parse token to obtain adminID
	admin, err := auth.ValidateAndParseToken(r)
	if err != nil {
		fmt.Println("Error validating token: ", err)
		return err
	}
	// Convert adminID to int
	intAdminID, err := strconv.ParseInt(admin.UserID, 10, 64)
	if err != nil {
		fmt.Println("Error converting adminID to uint: ", err)
		return err
	}

	// Build action using record action
	action := &models.CreateAction{
		ActionType:  recordAction.ActionType,
		EntityType:  recordAction.EntityType,
		EntityID:    recordAction.EntityID,
		Changes:     changes,
		Description: changeDescription,
		IPAddress:   r.RemoteAddr,
		AdminID:     uint(intAdminID),
	}

	_, err = s.Create(action)
	if err != nil {
		fmt.Println("Error recording action: ", err)
		return err
	}
	return nil
}
func (s *actionService) RecordBulkDelete(r *http.Request, schemaName, pluralSchemaName string,  schemaIDs []int, recordAction *models.RecordedAction) error {
	// Validate and parse token to obtain adminID
	admin, err := auth.ValidateAndParseToken(r)
	if err != nil {
		fmt.Println("Error validating token: ", err)
		return err
	}
	// Convert adminID to int
	intAdminID, err := strconv.ParseInt(admin.UserID, 10, 64)
	if err != nil {
		fmt.Println("Error converting adminID to uint: ", err)
		return err
	}

	// Build action using record action
	action := &models.CreateAction{
		ActionType:  recordAction.ActionType,
		EntityType:  schemaName,
		EntityID:    fmt.Sprint(schemaIDs),
		Changes:     "{}",
		Description: fmt.Sprintf("Bulk Deleted %d %s: %s", len(schemaIDs), pluralSchemaName, fmt.Sprint(schemaIDs)),
		IPAddress:   r.RemoteAddr,
		AdminID:     uint(intAdminID),
	}

	_, err = s.Create(action)
	if err != nil {
		fmt.Println("Error recording action: ", err)
		return err
	}
	return nil
}
// Creates a action in the database
func (s *actionService) Create(action *models.CreateAction) (*db.Action, error) {
	// Map incoming DTO to db schema
	toCreate := db.Action{
		ActionType:  action.ActionType,
		EntityType:  action.EntityType,
		EntityID:    action.EntityID,
		Changes:     action.Changes,
		Description: action.Description,
		AdminID:     action.AdminID,
		IPAddress:   action.IPAddress,
	}

	// Create above action in database
	created, err := s.repo.Create(&toCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating action: %w", err)
	}

	return created, nil
}
// Find a list of actions in the database
func (s *actionService) FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[db.Action], error) {
	actions, err := s.repo.FindAll(limit, offset, order, conditions)
	if err != nil {
		return nil, err
	}
	return actions, nil
}
// Find action in database by ID
func (s *actionService) FindById(id int) (*db.Action, error) {
	// Search cache
	// Define a key with a naming convention
	cacheKey := fmt.Sprintf("action:%d", id)
	// Check if action is in cache
	cachedPost, found := app.Cache.Load(cacheKey)
	if found {
		// If found, return cached action
		return cachedPost.(*db.Action), nil
	}

	// Find action by id
	action, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		return nil, err
	}

	// Store action in cache
	app.Cache.Store(cacheKey, action)

	return action, nil
}
// Delete action in database
func (s *actionService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting action: ", err)
		return err
	}
	// else
	cacheKey := fmt.Sprintf("action:%d", id)
	app.Cache.Delete(cacheKey)
	return nil
}
// Deletes multiple actions in database
func (s *actionService) BulkDelete(ids []int) error {
	err := s.repo.BulkDelete(ids)
	// If error detected
	if err != nil {
		fmt.Println("error in bulk deleting actions: ", err)
		return err
	}
	// else
	for _, id := range ids {
		cacheKey := fmt.Sprintf("action:%d", id)
		app.Cache.Delete(cacheKey)
	}
	return nil
}
// Updates action in database
func (s *actionService) Update(id int, action *models.UpdateAction) (*db.Action, error) {
	// Create action type from incoming DTO
	toUpdate := &db.Action{
		ActionType:  action.ActionType,
		EntityType:  action.EntityType,
		EntityID:    action.EntityID,
		Changes:     action.Changes,
		Description: action.Description,
	}

	// Update using repo
	updated, err := s.repo.Update(id, toUpdate)
	if err != nil {
		return nil, err
	}

	// Store updated action in cache
	cacheKey := fmt.Sprintf("action:%d", id)
	app.Cache.Store(cacheKey, updated)

	return updated, nil
}
