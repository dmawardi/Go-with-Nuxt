package helpers_test

import (
	"reflect"
	"testing"

	"github.com/dmawardi/Go-Template/internal/helpers/request"
	"github.com/dmawardi/Go-Template/internal/helpers/utility"
	"github.com/dmawardi/Go-Template/internal/models"
)

func TestGoValidateStruct(t *testing.T) {
	var testTable = []struct {
		name             string
		structToValidate interface{}
		expected         interface{}
		isErr            bool
	}{
		{"create-user-valid-data", &models.CreateUser{
			Username: "Hellow",
			Name:     "Psygog",
			Password: "pass-sword",
			Email:    "Regal@gmail.com",
		}, &models.ValidationError{}, false},
		{"create-user-missing-data", &models.CreateUser{
			Username: "Hellow",
			Password: "pass-sword",
			Email:    "Regal@gmail.com",
		}, &models.ValidationError{
			Validation_errors: map[string][]string{
				"name": []string{"non zero value required"},
			},
		}, true},
		{"create-user-invalid-data", &models.CreateUser{
			Username: "sway",
			Name:     "Psygo",
			Password: "pas",
			Email:    "Regal",
		}, &models.ValidationError{
			Validation_errors: map[string][]string{
				"username": []string{"sway does not validate as length(6|25)"},
				"name":     []string{"Psygo does not validate as length(6|80)"},
				"email":    []string{"Regal does not validate as email"},
				"password": []string{"pas does not validate as length(6|30)"},
			},
		}, true},
		{"create-user-some-invalid", &models.CreateUser{
			Username: "Hel",
			Name:     "Psybaba",
			Password: "pass-sword",
			Email:    "Regal@gmail.com",
		}, &models.ValidationError{
			Validation_errors: map[string][]string{
				"username": []string{"Hel does not validate as length(6|25)"},
			},
		}, true},
		{"update-user-valid-partial-data", &models.UpdateUser{
			Username: "Hellow",
			Name:     "Psygow",
			Email:    "Regal@gmail.com",
		}, &models.ValidationError{}, false},
		{"update-user-invalid-partial-data", &models.UpdateUser{
			Username: "sway",
			Name:     "Psygo",
			Password: "pas",
		}, &models.ValidationError{
			Validation_errors: map[string][]string{
				"username": []string{"sway does not validate as length(6|25)"},
				"name":     []string{"Psygo does not validate as length(6|80)"},
				"password": []string{"pas does not validate as length(6|30)"},
			},
		}, true},
	}
	// for test struct in tests array
	for _, tt := range testTable {
		pass, validation := request.GoValidateStruct(tt.structToValidate)
		// Check that errors happen when they're supposed to
		if tt.isErr {
			if pass == true {
				t.Errorf("expected an error but didn't get one: %s", tt.name)
			}
			// Check that error values are as expected
			// if validation != tt.expected {
			// 	t.Errorf("Error: %s value received: %v\n not as expected: %v\n", tt.name, tt.expected, validation)
			// }

			if !reflect.DeepEqual(validation, tt.expected) {
				t.Errorf("Error: %s value received: %v\n not as expected: %v\n", tt.name, tt.expected, validation)
			}
		} else {
			if pass == false {
				t.Errorf("did not expect an error in %s but got one: %v\n", tt.name, validation)
			}
		}

	}
}

// TestGenerateRandomString tests the GenerateRandomString function.
func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{"zero length", 0, false},
		{"positive length", 10, false},
		{"large length", 1000, false},
		// Can't really force an error in this case since rand.Read rarely fails,
		// but you could hypothetically test that case as well.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utility.GenerateRandomString(tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.length {
				t.Errorf("GenerateRandomString() got length = %v, want %v", len(got), tt.length)
			}
			for _, r := range got {
				if !isLetterOrDigit(r) {
					t.Errorf("GenerateRandomString() got a non-letter/digit character: %v", r)
				}
			}
		})
	}
}

// (helper function for generate random string) isLetterOrDigit checks if a rune is a letter or a digit.
func isLetterOrDigit(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9')
}
