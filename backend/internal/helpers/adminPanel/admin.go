package adminpanel

import (
	"net/http"
)

// parseFormToMap parses the form data and converts it into a map[string]string
func ParseFormToMap(r *http.Request) (map[string]string, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	formMap := make(map[string]string)
	for key, values := range r.Form { // range over map
		// In form data, key can have multiple values,
		// we'll take the first one only
		formMap[key] = values[0]
	}

	return formMap, nil
}
