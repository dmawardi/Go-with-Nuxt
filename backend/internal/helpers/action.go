package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// GenerateChangeLog compares two objects of the same type and returns a JSON string representing the differences.
func GenerateChangeLog(changeType string, changeLog ChangeLogInput) (string, error) {
	// Get struct value of the old object (in case it's a pointer)
	oldObj, err := GetStructValue(changeLog.OldObj) 
	if err != nil {
		return "", err
	}
	newObj, err := GetStructValue(changeLog.NewObj)
	if err != nil {
		return "", err
	}
	// Use reflection to get the values of the old and new objects
	oldVal := reflect.ValueOf(oldObj)
	newVal := reflect.ValueOf(newObj)
	
	// Initialize a map to hold the change log
	changeLogMap := make(map[string]map[string]interface{})

	// Get the type of the old object (if provided) or new object (for creation)
	objType := oldVal.Type()
	if changeType == "create" {
		objType = newVal.Type()
	}

	// Handle based on the change type
	switch changeType {
	case "create":
		// Iterate through the fields of the new object and log them as new values
		for i := 0; i < newVal.NumField(); i++ {
			// Grab field name and field value
			fieldName := objType.Field(i).Name
			fieldValue := newVal.Field(i).Interface()

			// Skip the time fields (CreatedAt, UpdatedAt, DeletedAt)
			if fieldName != "CreatedAt" && fieldName != "UpdatedAt" && fieldName != "DeletedAt" {
				// Record the new value for the field
				changeLogMap[fieldName] = map[string]interface{}{
					"old": nil,
					"new": fieldValue,
				}
			}
		}
	case "update":
		// Check if the types of the two objects are the same
		if oldVal.Type() != newVal.Type() {
			return "", fmt.Errorf("type mismatch: %v vs %v", oldVal.Type(), newVal.Type())
		}
		// For updates, compare the old and new objects and log changes
		for i := 0; i < oldVal.NumField(); i++ {
			fieldName := objType.Field(i).Name
			oldFieldValue := oldVal.Field(i).Interface()
			newFieldValue := newVal.Field(i).Interface()

			// If the field values are different, log the change
			if !reflect.DeepEqual(oldFieldValue, newFieldValue) {

			// Skip the time fields (CreatedAt, UpdatedAt, DeletedAt)
			if fieldName != "CreatedAt" && fieldName != "UpdatedAt" && fieldName != "DeletedAt" {
				// Record the old and new values for the field
				changeLogMap[fieldName] = map[string]interface{}{
					"old": oldFieldValue,
					"new": newFieldValue,
				}
			}
			}
		}
	case "delete":
		// For deletion, log all fields from the old object as old
		for i := 0; i < oldVal.NumField(); i++ {
			fieldName := objType.Field(i).Name
			fieldValue := oldVal.Field(i).Interface()
			// Skip the time fields (CreatedAt, UpdatedAt, DeletedAt)
			if fieldName != "CreatedAt" && fieldName != "UpdatedAt" && fieldName != "DeletedAt" {
				// Record the old value for the field
				changeLogMap[fieldName] = map[string]interface{}{
					"old": fieldValue,
					"new": nil,
				}
			}
		}
	default:
		return "", fmt.Errorf("invalid change type: %s", changeType)
	}

	// Marshal the change log map to a JSON string
	changeLogJSON, err := json.Marshal(changeLogMap)
	if err != nil {
		return "", err
	}

	return string(changeLogJSON), nil
}
// GetStructValue returns the struct value if the input is a struct or a pointer to a struct.
// Returns an error if the input is neither a struct nor a pointer to a struct.
func GetStructValue(input interface{}) (interface{}, error) {
	// First, check if the input is already a struct
	val := reflect.ValueOf(input)

	if val.Kind() == reflect.Struct {
		// Input is already a struct, return it directly
		return input, nil
	}

	// If the value is a pointer, dereference it
	if val.Kind() == reflect.Ptr {
		// Dereference the pointer to get the underlying value
		val = val.Elem()

		// Check if the dereferenced value is a struct
		if val.Kind() == reflect.Struct {
			// Return the struct as an interface{}
			return val.Interface(), nil
		}
	}

	// If it's neither a struct nor a pointer to a struct, return an error
	return nil, errors.New("input is neither a struct nor a pointer to a struct")
}
// GenerateChangeDescription generates a human-readable description of the changes
func GenerateChangeDescription(changeLogJSON string, entityType string, actionType string, entityID uint) (string, error) {
	// If the entity was deleted, return a deletion description
	if actionType == "delete" {
		return fmt.Sprintf("Deleted the %s with ID: %v", entityType, entityID), nil
	}
	if actionType == "create" {
		return fmt.Sprintf("Created a new %s with ID: %v", entityType, entityID), nil
	}

	// Initialize a map to hold the change log data
	var changeLogMap map[string]map[string]interface{}
	
	// Unmarshal the JSON change log into the map
	err := json.Unmarshal([]byte(changeLogJSON), &changeLogMap)
	if err != nil {
		return "", err
	}

	// If the change log is empty, return a message indicating no changes were detected
	if len(changeLogMap) == 0 {
		return fmt.Sprintf("No changes detected for the %s.", entityType), nil
	}

	// Prepare a slice to hold the field names that were changed
	var changes []string
	for field := range changeLogMap {
		changes = append(changes, fmt.Sprintf("%s", strings.ToLower(field)))
	}

	// If there is only one change, specify that field in the description
	if len(changes) == 1 {
		return fmt.Sprintf("Updated the %s %s.", entityType, changes[0]), nil
	}

	// If there are multiple changes, join them with "and" and specify in the description
	return fmt.Sprintf("Updated the %s %s.", entityType, strings.Join(changes, " and ")), nil
}
// ChangeLogInput represents the input data for generating a change log
// It contains the old and new objects to compare
type ChangeLogInput struct {
	OldObj interface{}
	NewObj interface{}
}