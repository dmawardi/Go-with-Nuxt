package request

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dmawardi/Go-Template/internal/models"
)

// Validation
// Takes in a list of errors from Go Validator and formats into JSON ready struct
func CreateStructFromValidationErrorString(errs []error) *models.ValidationError {
	// Prepare validation model for appending errors
	validation := &models.ValidationError{
		Validation_errors: make(map[string][]string),
	}
	// Loop through slice of errors
	for _, e := range errs {
		// Grab error strong
		errorString := e.Error()
		// Split by colon
		errorArray := strings.Split(errorString, ": ")

		// Prepare err message array for map preparation
		var errMessageArray []string
		// Append the error to the array
		errMessageArray = append(errMessageArray, errorArray[1])

		// Assign to validation struct
		validation.Validation_errors[errorArray[0]] = errMessageArray
	}

	return validation
}

// Uses a DTO struct's key value "valid" config to assess whether it's valid
// then returns a struct ready for JSON marshal
func GoValidateStruct(objectToValidate interface{}) (bool, *models.ValidationError) {
	// Validate the incoming DTO
	_, err := govalidator.ValidateStruct(objectToValidate)

	// if no error found
	if err != nil {
		// Prepare slice of errors
		errs := err.(govalidator.Errors).Errors()

		// Grabs the error slice and creates a front-end ready validation error
		validationResponse := CreateStructFromValidationErrorString(errs)
		// return failure on validation and validation response
		return false, validationResponse
	}
	// Return pass on validation and empty validation response
	return true, &models.ValidationError{}
}
