package moduleservices

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

var app *config.AppConfig

func SetAppConfig(appConfig *config.AppConfig) {
	app = appConfig
}


// BasicModuleService is an interface for basic service CRUD operations
type BasicModuleService[dbSchema, create, update any] interface {
	FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[dbSchema], error)
	FindById(int) (*dbSchema, error)
	Create(entity *create) (*dbSchema, error)
	Update(int, *update) (*dbSchema, error)
	Delete(int) error
	BulkDelete([]int) error
}
// A generic struct for basic service
type BasicServiceStruct[dbSchema, createDTO, updateDTO any] struct {
	Repo repository.BasicModuleRepository[dbSchema]
	schemaName string
	// Mapping functions
	mapCreateToDbSchema func(*createDTO) *dbSchema
	mapUpdateToDbSchema func(*updateDTO) *dbSchema
}
// Returns a new basic service
func newBasicModuleService[dbSchema, createDTO, updateDTO any](repo repository.BasicModuleRepository[dbSchema]) BasicModuleService[dbSchema, createDTO, updateDTO] {
	return &BasicServiceStruct[dbSchema, createDTO, updateDTO]{
		Repo: repo,
	}
}

// Receiver Functions
// 
// 
// Creates a new entity in database
func (s *BasicServiceStruct[dbSchema, createDTO, updateDTO]) Create(create *createDTO) (*dbSchema, error) {
	// Maps incoming DTO to db schema
	toCreate := s.mapCreateToDbSchema(create)

	// Create above post in database
	created, err := s.Repo.Create(toCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creation of %s: %w", s.schemaName, err)
	}

	return created, nil
}
// Find all entities in database
func (s *BasicServiceStruct[dbSchema, createDTO, updateDTO]) FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[dbSchema], error) {
	entities, err := s.Repo.FindAll(limit, offset, order, conditions)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
// Find entity by id
func (s *BasicServiceStruct[dbSchema, createDTO, updateDTO]) FindById(id int) (*dbSchema, error) {
	// Search cache
	// Define a key with a naming convention
	cacheKey := fmt.Sprintf("%s:%d", s.schemaName, id)
	// Check if entity is in cache
	cachedPost, found := app.Cache.Load(cacheKey)
	if found {
		// If found, return cached post
		return cachedPost.(*dbSchema), nil
	}

	// Find entity by id
	entity, err := s.Repo.FindById(id)
	// If error detected
	if err != nil {
		return nil, err
	}

	// Store entity in cache
	app.Cache.Store(cacheKey, entity)

	return entity, nil
}
// Delete entity in database
func (s *BasicServiceStruct[dbSchema, createDTO, updateDTO]) Delete(id int) error {
	err := s.Repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Printf("error in deleting %s: %s", s.schemaName, err)
		return err
	}
	// else
	cacheKey := fmt.Sprintf("%s:%d", s.schemaName, id)
	app.Cache.Delete(cacheKey)
	return nil
}
// Deletes multiple entities in database
func (s *BasicServiceStruct[dbSchema, createDTO, updateDTO]) BulkDelete(ids []int) error {
	err := s.Repo.BulkDelete(ids)
	// If error detected
	if err != nil {
		fmt.Printf("error in bulk deleting %s: %s", s.schemaName, err)
		return err
	}
	// else
	for _, id := range ids {
		cacheKey := fmt.Sprintf("%s:%d", s.schemaName, id)
		app.Cache.Delete(cacheKey)
	}
	return nil
}
// Updates entity in database
func (s *BasicServiceStruct[dbSchema, createDTO, updateDTO]) Update(id int, update *updateDTO) (*dbSchema, error) {
	// Create entity type from incoming DTO
	toUpdate := s.mapUpdateToDbSchema(update)

	// Update using Repo
	updated, err := s.Repo.Update(id, toUpdate)
	if err != nil {
		return nil, err
	}

	// Store updated entity in cache
	cacheKey := fmt.Sprintf("%s:%d", s.schemaName, id)
	app.Cache.Store(cacheKey, updated)

	return updated, nil
}

