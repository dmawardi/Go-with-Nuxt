package adminpanel

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dmawardi/Go-Template/internal/models"
)

// FORM DATA AND FUNCTIONS
//

// Used to build form in Go Templates
// Primary form data used in PagerenderData
type FormData struct {
	// Form title
	FormFields []FormField
	// Contains form action and method
	FormDetails FormDetails
}

// Data for each form field
type FormField struct {
	// Label above input
	Label string
	// Labed used for populating values and placeholders
	DbLabel string
	// Used for form submission
	Name string
	// Is this field required?
	Required bool
	// Is this field disabled?
	Disabled bool

	// Current value
	Value string
	// Silhouette
	Placeholder string
	Type        string
	Errors      []ErrorMessage
	Selectors   []FormFieldSelector
}

// Used to store data to render form selectors in Go Templates
type FormFieldSelector struct {
	Value    string
	Label    string
	Selected bool
}

// Display of errors in form
type ErrorMessage string

// Form Details
type FormDetails struct {
	FormAction string
	FormMethod string
}

// Map used to group form selectors for a schema (eg. FormSelector["field_name"])
// Form selectors are used to index the form selectors as functions that return map[string]string
type FormSelectors map[string]func() []FormFieldSelector

// Sets the Errors field in each field of a form
func SetValidationErrorsInForm(form []FormField, validationErrors models.ValidationError) {
	// Iterate through fields
	for i, field := range form {
		// Check if field name is in validation errors
		if errors, ok := validationErrors.Validation_errors[field.Name]; ok {
			// If found, iterate through errors
			for _, err := range errors {
				// Get error message
				errorMessage := ErrorMessage(err)
				// If error contains does not, remove validated text
				if strings.Contains(err, "does not") {
					// Split string using does
					split := strings.Split(err, "does")
					// Update error message with rebuilt string
					errorMessage = ErrorMessage(fmt.Sprintf("Does %s", split[1]))
				}
				// Append error message to field
				form[i].Errors = append(form[i].Errors, errorMessage)
			}
		}
	}
}

// Used to populate FormField values with placeholder values found from previously submitted form
func populateFormValuesWithSubmittedFormMap(form *[]FormField, fieldMap map[string]string) error {
	// Loop through fields in the form and populate placeholders
	for i := range *form {
		// Get pointer to field
		field := &(*form)[i]
		// Only if field is not a password field
		if field.Type != "password" {
			// If the field exists in the map, populate the placeholder
			if val, ok := fieldMap[field.Name]; ok {
				field.Value = val
			} else {
				field.Value = ""
			}
		}
	}
	return nil
}

// Used to populate form field placeholders with data from database (that has been converted to map[string]string)
func populatePlaceholdersWithDBData(form *[]FormField, fieldMap map[string]string) error {
	// Loop through fields and populate placeholders
	for i := range *form {
		// Get pointer to field
		field := &(*form)[i]
		if field.Type == "select" {
			// Update selectors with current value selected
			setDefaultSelected(field.Selectors, fieldMap[field.DbLabel])
			// Else treat as ordinary input
		} else {
			// If the field exists in the map, populate the placeholder
			if val, ok := fieldMap[field.DbLabel]; ok {
				// Populate placeholder as value from field map
				field.Placeholder = val
			} else {
				field.Placeholder = ""
			}
		}
	}

	return nil
}

func populateValuessWithDBData(form *[]FormField, fieldMap map[string]string) error {
	// Loop through fields and populate placeholders
	for i := range *form {
		// Get pointer to field
		field := &(*form)[i]
		if field.Type == "select" {
			// Update selectors with current value selected
			setDefaultSelected(field.Selectors, fieldMap[field.DbLabel])
			// Else treat as ordinary input
		} else {
			// If the field exists in the map, populate the placeholder
			if val, ok := fieldMap[field.DbLabel]; ok {
				// Populate placeholder as value from field map
				field.Value = val
			} else {
				field.Placeholder = ""
			}
		}
	}

	return nil
}

// getValuesUsingFieldMap returns a map of field names to string representations of their values for a given struct
func getValuesUsingFieldMap(entity interface{}) map[string]string {
	fieldMap := make(map[string]string) // Initialize the map to store field names and their string values

	val := reflect.ValueOf(entity) // Get the reflect value of the entity
	typ := reflect.TypeOf(entity)  // Get the reflect type of the entity

	// Iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)       // Get the field value
		fieldType := typ.Field(i)   // Get the field type
		var fieldValue string       // Variable to store the string representation of the field value

		// Handle different field types
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// Handle integer types
			fieldValue = fmt.Sprint(field.Int())
		case reflect.String:
			// Handle string types
			fieldValue = field.String()
		case reflect.Struct:
			// Special handling for struct types, specifically time.Time
			if fieldType.Type == reflect.TypeOf(time.Time{}) {
				fieldValue = field.Interface().(time.Time).Format("January 2, 2006 at 3:04pm")
			}
		default:
			// Default case for other types
			fieldValue = fmt.Sprint(field.Interface())
		}

		// Add the field name and its string value to the map
		fieldMap[fieldType.Name] = fieldValue
	}

	return fieldMap // Return the populated field map
}