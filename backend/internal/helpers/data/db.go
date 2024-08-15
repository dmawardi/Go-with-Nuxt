package data

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

// Query all records in a table based on conditions (results are populated into a slice of the schema type)
func QueryAll(dbClient *gorm.DB, dbSchema interface{}, limit, offset int, order string, conditions []models.QueryConditionParameters, preloads []string) error {
	// Build base query for query schema
	query := dbClient.Model(dbSchema)

	// Add parameters into query as needed
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if order != "" {
		// Add order to query
		query = query.Order(order)
	}

	query = AddWhereConditionsToQuery(query, conditions)

	// Iterate through preloads to preload foreign key fields
	for _, fieldToPreload := range preloads {
		query = query.Preload(fieldToPreload)
	}

	// Query database
	if err := query.Find(dbSchema).Error; err != nil {
		fmt.Printf("Error in query all: %v\n", err)
		return err
	}

	return nil
}

// Iterates through conditions and adds them to the query
// Initial is added as where followed by subsequent conditions as or
func AddWhereConditionsToQuery(query *gorm.DB, conditions []models.QueryConditionParameters) *gorm.DB {
	// Iterate through conditions (stop at second last element)
	// Increment by 2 to account for condition and value
	if len(conditions) > 0 {
		for i, cond := range conditions {
			// For the first condition, use Where
			if i == 0 {
				query = query.Where(cond.Condition, cond.Value)
			} else {
				// For subsequent conditions, use Or
				query = query.Or(cond.Condition, cond.Value)
			}
		}
	}

	return query
}

// Build Meta data for QueryAll requests
func BuildMetaData(dbClient *gorm.DB, dbSchema interface{}, limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.SchemaMetaData, error) {
	// Fetch metadata from database
	var totalCount *int64

	// Count the total number of records
	totalCount, err := CountBasedOnConditions(dbSchema, conditions, dbClient)
	if err != nil {
		return nil, err
	}
	// Find the total number of pages from total count and limit
	totalPages := int(*totalCount) / limit
	if int(*totalCount)%limit != 0 {
		totalPages += 1
	}
	// Calculate current page
	currentPage := offset/limit + 1
	// Calculate next page
	var nextPage *int // Using a pointer to represent the absence of a next page
	if currentPage < totalPages {
		next := currentPage + 1
		nextPage = &next
	}
	// Calculate previous page
	var prevPage *int // Using a pointer to represent the absence of a previous page
	if currentPage > 1 {
		prev := currentPage - 1
		prevPage = &prev
	}
	// Build metadata object
	metaData := models.NewSchemaMetaData(*totalCount, limit, totalPages, currentPage, nextPage, prevPage)

	// Return meta data
	return &metaData, nil
}

// Count using conditions
func CountBasedOnConditions(databaseSchema interface{}, conditions []models.QueryConditionParameters, dbClient *gorm.DB) (*int64, error) {
	// Fetch metadata from database
	var totalCount int64

	// Count the total number of records
	query := dbClient.Model(databaseSchema)

	if len(conditions) > 0 {
		for i, cond := range conditions {
			// For the first condition, use Where
			if i == 0 {
				query = query.Where(cond.Condition, cond.Value)
			} else {
				// For subsequent conditions, use Or
				query = query.Or(cond.Condition, cond.Value)
			}
		}
	}

	// Execute query
	countResult := query.Count(&totalCount)
	if countResult.Error != nil {
		return nil, countResult.Error
	}
	return &totalCount, nil
}

// Bulk deletes all records within a table based on ids
func BulkDeleteByIds(databaseSchema interface{}, ids []int, dbClient *gorm.DB) error {
	// Start a transaction (to avoid partial deletion)
	err := dbClient.Transaction(func(tx *gorm.DB) error {
		// In the transaction, delete users with specified IDs
		if err := tx.Where("id IN ?", ids).Delete(&databaseSchema).Error; err != nil {
			return err // Return any error to rollback the transaction
		}

		return nil // Return nil to commit the transaction
	})

	// Check if the transaction was successful
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Only allows the insertion of unique record into the database
func InsertUniqueRecord(db *gorm.DB, items []interface{}) error {
	for _, item := range items {
		// Attempt to find the existing item or create a new one if not found.
		// This uses all non-zero fields of the item to check for an existing record.
		result := db.FirstOrCreate(item, item)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			fmt.Println("Item already exists, skipping...")
		} else {
			fmt.Printf("Seed item created: %v\n", item)
		}
	}
	return nil
}
