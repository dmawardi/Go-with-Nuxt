package data

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// Converts list of strings to list of ints
func ConvertStringSliceToIntSlice(stringSlice []string) ([]int, error) {
	intSlice := make([]int, 0, len(stringSlice)) // Create a slice of ints with the same length

	for _, str := range stringSlice {
		num, err := strconv.Atoi(str) // Convert string to int
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, num) // Append the converted int to the slice
	}
	return intSlice, nil
}

// Convert json string to map[string]string
func JsonToMap(jsonStr string) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Convert struct string to map. Struct string will be in format: "[key1]:value1|[key2]:value2"
func StringToMap(input string) (map[string]string, error) {
	result := make(map[string]string)

	// Split the input string by "|"
	parts := strings.Split(input, "|")

	for _, part := range parts {
		// Check if the part has "[]" to identify a key name
		keyValueSlice := strings.Split(part, ":")

		// If a key value pair is found
		if len(keyValueSlice) == 2 {
			// Grab the first item in slice as key, and remove the "[" and "]" characters
			key := keyValueSlice[0]
			// // Grab the second item in slice as value
			value := keyValueSlice[1]
			// Add key value pair to result map
			result[key[1:len(key)-1]] = value
		}

	}

	return result, nil
}

// Checks if a slice contains a type of targetType. Returns boolean
func CheckSliceType(slice interface{}, targetType reflect.Type) bool {
	// Get the type of the slice argument
	sliceType := reflect.TypeOf(slice)

	// Check if the slice argument is actually a slice
	if sliceType.Kind() != reflect.Slice {
		return false
	}

	// Compare element type of the slice with the target type
	elemType := sliceType.Elem()
	return elemType == targetType
}
