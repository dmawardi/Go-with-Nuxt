package core

// Used to init the query params for easy extraction in controller
// Returns: map[string]string{"age": "int", "name": "string", "active": "bool"}
func ActionConditionQueryParams() map[string]string {
	return map[string]string{
		"action_type":    "string",
		"entity_type":     "string",
		"description": "string",	
	}
}