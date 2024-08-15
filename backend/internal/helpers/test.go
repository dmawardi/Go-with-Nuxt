package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Setup database connection
func SetupTestDatabase() *gorm.DB {
	// Open a new, temporary database for testing
	dbClient, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		fmt.Printf("failed to open database: %v", err)
	}

	// Migrate the database schema
	for _, table := range db.Models {
		if err := dbClient.AutoMigrate(table); err != nil {
			fmt.Printf("failed to migrate database schema: %v", err)
		}
	}

	return dbClient
}

// CompareObjects compares the specified fields of two interface{} objects.
// It uses reflection to dynamically compare the field values of both objects.
func CompareObjects(actualObject interface{}, expectedObject interface{}, t *testing.T, fieldsToCheck []string) {
	// Convert both objects to reflect.Value to facilitate comparison.
	actualValue := reflect.ValueOf(actualObject)
	if actualValue.Kind() == reflect.Ptr {
		actualValue = actualValue.Elem()
	}
	expectedValue := reflect.ValueOf(expectedObject)
	if expectedValue.Kind() == reflect.Ptr {
		expectedValue = expectedValue.Elem()
	}

	// Iterate over the specified fields to compare their values.
	for _, field := range fieldsToCheck {
		actualFieldValue := actualValue.FieldByName(field)
		expectedFieldValue := expectedValue.FieldByName(field)

		// Check if both fields are valid.
		if !actualFieldValue.IsValid() {
			t.Errorf("actual object does not have field %s", field)
			continue
		}
		if !expectedFieldValue.IsValid() {
			t.Errorf("expected object does not have field %s", field)
			continue
		}

		// Compare the actual and expected field values.
		if !reflect.DeepEqual(actualFieldValue.Interface(), expectedFieldValue.Interface()) {
			t.Errorf("field %s does not match: expected %v, got %v", field, expectedFieldValue.Interface(), actualFieldValue.Interface())
		}
	}
}

// UpdateModelFields updates the fields of a GORM model based on a map[string]string.
// The model parameter is expected to be a pointer to a struct that's a GORM model.
// The updates parameter is a map where keys are field names and values are new values for those fields, as strings.
func UpdateModelFields(model interface{}, updates map[string]string) error {
	// Ensure the model is a pointer to a struct.
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr || modelValue.Elem().Kind() != reflect.Struct {
		return errors.New("model must be a pointer to a struct")
	}

	// Get the underlying struct value.
	structValue := modelValue.Elem()

	// Iterate through the updates map to update struct fields.
	for field, newValue := range updates {
		// Find the struct field.
		structField := structValue.FieldByName(field)
		if !structField.IsValid() {
			return fmt.Errorf("no such field: %s in model", field)
		}

		// Ensure the field can be set.
		if !structField.CanSet() {
			return fmt.Errorf("cannot set field: %s", field)
		}

		// Convert and set the field value based on its kind.
		switch structField.Kind() {
		case reflect.String:
			structField.SetString(newValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(newValue, 10, 64)
			if err != nil {
				return fmt.Errorf("cannot convert %s to int for field: %s", newValue, field)
			}
			structField.SetInt(intVal)
		case reflect.Float32, reflect.Float64:
			floatVal, err := strconv.ParseFloat(newValue, 64)
			if err != nil {
				return fmt.Errorf("cannot convert %s to float for field: %s", newValue, field)
			}
			structField.SetFloat(floatVal)
		// Add more cases here for other types as needed.
		default:
			return fmt.Errorf("unsupported field type: %s for field: %s", structField.Type(), field)
		}
	}

	return nil
}

// Update a struct field dynamically
func UpdateStructField(structPtr interface{}, fieldName string, fieldValue interface{}) error {
	value := reflect.ValueOf(structPtr)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("invalid struct pointer")
	}

	structValue := value.Elem()
	if !structValue.CanSet() {
		return fmt.Errorf("cannot set struct field value")
	}

	field := structValue.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("invalid struct field name")
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set struct field value")
	}

	fieldValueRef := reflect.ValueOf(fieldValue)
	if !fieldValueRef.Type().AssignableTo(field.Type()) {
		return fmt.Errorf("field value type mismatch")
	}

	field.Set(fieldValueRef)
	return nil
}

// Test helper function: Hashes password and generates a new user in the database
func HashPassAndGenerateUserInDb(user *db.User, client *gorm.DB, t *testing.T) (*db.User, error) {
	// Hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		t.Fatalf("Couldn't hash password")
	}
	user.Password = string(hashedPass)

	// Create user
	createResult := client.Create(user)
	if createResult.Error != nil {
		t.Fatalf("Couldn't create user: %v", user.Email)
	}

	return user, nil
}

// Mocking
// Mock Email
type EmailMock struct {
}

func (e *EmailMock) SendEmail(recipient, subject, body string) error {
	return nil
}

// A helper function to build an API request that starts with url of '/api/'
func BuildApiRequest(method string, urlSuffix string, body io.Reader, authHeaderRequired bool, token string) (request *http.Request, err error) {
	req, err := http.NewRequest(method, fmt.Sprintf("/api/%v", urlSuffix), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	// If authorization header required
	if authHeaderRequired {
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", token))
	}
	return req, nil
}

// Build a struct object to a type of bytes.reader to fulfill io.reader interface
func BuildReqBody(data interface{}) *bytes.Reader {
	// Marshal to JSON
	marshalled, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Failed to marshal JSON")
	}
	// Make into reader
	readerReqBody := bytes.NewReader(marshalled)
	return readerReqBody
}
