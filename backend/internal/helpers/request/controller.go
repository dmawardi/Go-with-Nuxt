package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dmawardi/Go-Template/internal/models"
)

// Controller helper functions
// Takes struct data and returns as JSON to Response writer
func WriteAsJSON(w http.ResponseWriter, data interface{}) error {
	// Edit content type
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	// Write data as response
	w.Write(jsonData)
	return nil
}

// URL parameter extraction helper functions
// Special handling for search query if found
func addSearchQueryToConditions(r *http.Request, conditionsToExtract map[string]string, currentConditions []models.QueryConditionParameters) []models.QueryConditionParameters {
	// Prepare URL query parameters
	queryParams := r.URL.Query()
	// Prepare search query
	searchQuery := ""
	if searchValue := queryParams.Get("search"); searchValue != "" {
		// Surround search value with % to make it a LIKE query
		searchQuery = "%" + searchValue + "%"
		// Iterate through query list to add the search condition to each as a LIKE query
		for param, conditionType := range conditionsToExtract {
			// If query parameter is string
			if conditionType == "string" {
				// Add query to conditions with param name and make case insensitive
				lowerCaseValue := strings.ToLower(fmt.Sprintf("%v", searchQuery))
				currentConditions = append(currentConditions, models.QueryConditionParameters{Condition: fmt.Sprintf("LOWER(%s) LIKE ?", param), Value: lowerCaseValue})
			}
		}
	}

	return currentConditions
}

// Accepts request and a slice of conditions to extract from the request
// Extracts as a slice of interfaces that are structured as [condition, value]
func ExtractSearchAndConditionParams(r *http.Request, conditionsToExtract map[string]string) ([]models.QueryConditionParameters, error) {
	// Prepare URL query parameters
	queryParams := r.URL.Query()
	// Prepare slice of conditions
	var extractedConditions []models.QueryConditionParameters

	// Special handling for search query if found
	extractedConditions = addSearchQueryToConditions(r, conditionsToExtract, extractedConditions)

	// Iterate through query list
	for param, conditionType := range conditionsToExtract {
		// If query parameter is present and not empty
		if queryValue := queryParams.Get(param); queryValue != "" {
			// Prepare variables
			var condition string
			var value interface{}
			var err error

			// Detecting prefixes for operators
			switch {
			// If greater than ie. age=gt:20
			case strings.HasPrefix(queryValue, "gt:"):
				condition = param + " > ?"
				value, err = parseValue(queryValue[3:], conditionType)
			// If less than ie. age=lt:20
			case strings.HasPrefix(queryValue, "lt:"):
				condition = param + " < ?"
				value, err = parseValue(queryValue[3:], conditionType)
			// If greater than or equal to ie. age=gte:20
			case strings.HasPrefix(queryValue, "gte:"):
				condition = param + " >= ?"
				value, err = parseValue(queryValue[4:], conditionType)
			// If less than or equal to ie. age=lte:20
			case strings.HasPrefix(queryValue, "lte:"):
				condition = param + " <= ?"
				value, err = parseValue(queryValue[4:], conditionType)
			// If default equal condition
			default:
				condition = param + " = ?"
				value, err = parseValue(queryValue, conditionType)
			}

			// If an issue found, return error
			if err != nil {
				return nil, fmt.Errorf("invalid value for %s: %v", param, err)
			}

			// Append to slice
			extractedConditions = append(extractedConditions, models.QueryConditionParameters{Condition: condition, Value: value})
		}
	}

	return extractedConditions, nil
}

// Extracts basic pagination query parameters from request
func ExtractBasicFindAllQueryParams(r *http.Request) (models.BaseFindAllQueryParams, error) {
	// Grab query parameters
	page, err := GrabIntQueryParamOrDefault(r, "page", 1)
	if err != nil {
		return models.BaseFindAllQueryParams{}, err
	}
	limit, err := GrabIntQueryParamOrDefault(r, "limit", 10)
	if err != nil {
		return models.BaseFindAllQueryParams{}, err
	}
	// Grab order
	order := GrabQueryParamOrDefault(r, "order", "id")
	sqlOrder, err := convertToSQLOrderBy(order)
	if err != nil {
		return models.BaseFindAllQueryParams{}, err
	}
	// Calculate offset using pages and limit
	offset := (page - 1) * limit

	return models.BaseFindAllQueryParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
		Order:  sqlOrder,
	}, nil
}

// Converts orderby paramter to SQL ready format
func convertToSQLOrderBy(orderBy string) (string, error) {
	parts := strings.Split(orderBy, "_")
	// Assign first item as column name
	columnName := parts[0]

	if len(parts) == 2 {
		order := parts[1]
		// else, return
		return fmt.Sprintf("%s %s", columnName, strings.ToUpper(order)), nil
	}

	// else, return
	return columnName, nil
}

// Helper function to parse the value based on type
func parseValue(value, conditionType string) (interface{}, error) {
	switch conditionType {
	case "int":
		return strconv.Atoi(value)
	case "string":
		return value, nil
	case "bool":
		return strconv.ParseBool(value)
	default:
		return nil, fmt.Errorf("unknown condition type: %s", conditionType)
	}
}

// Grabs a query parameter from the request, if not present, returns default value
func GrabQueryParamOrDefault(r *http.Request, param string, defaultValue string) string {
	// Grab query parameters
	queryParam := r.URL.Query().Get(param)
	// Check if limit is available, if not, set to default
	if queryParam == "" {
		queryParam = defaultValue
	}
	return queryParam
}

// Grabs an INT type query parameter from the request, if not present, returns default value
func GrabIntQueryParamOrDefault(r *http.Request, param string, defaultValue int) (int, error) {
	// Grab query parameters
	queryParam := r.URL.Query().Get(param)
	// Check if limit is available, if not, set to default
	if queryParam == "" {
		return defaultValue, nil
	}
	// Convert to int
	intQuery, err := strconv.Atoi(queryParam)
	if err != nil {
		return 0, err
	}
	return intQuery, nil
}
