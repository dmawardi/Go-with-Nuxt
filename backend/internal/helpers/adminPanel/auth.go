package adminpanel

import (
	"github.com/dmawardi/Go-Template/internal/helpers/utility"
	"github.com/dmawardi/Go-Template/internal/models"
)

// Policy search helpers
// Searches a list of policies for a given resource based on search term
func SearchPoliciesByResource(maps []models.PolicyRuleCombinedActions, searchTerm string) []models.PolicyRuleCombinedActions {
	var result []models.PolicyRuleCombinedActions

	// Iterate through map of policies
	for _, m := range maps {

		// If success and resource contains search term
		if utility.ContainsString(m.Resource, searchTerm) {
			result = append(result, m)
		}
	}

	return result
}

// Searches a list of maps for a given key based on search term
func SearchMapKeysFor(maps []map[string]string, mapKeysToSearch []string, searchTerm string) []map[string]string {
	var result []map[string]string
	// Init to record if already added to results
	addedToResult := false

	// Iterate through map of policies
	for _, m := range maps {
		// Reset added to result
		addedToResult = false
		// Iterate through list of keys to search for term
		for _, keyToSearch := range mapKeysToSearch {
			// Grab value
			value, ok := m[keyToSearch]
			// If success, and the record hasn't been added already and value contains search term
			if ok && utility.ContainsString(value, searchTerm) && !addedToResult {
				// Append
				result = append(result, m)
				// Set added to true
				addedToResult = true
			}
		}
	}

	return result
}

// Searches a list of policies for a given resource based on search term
func SearchPoliciesForExactResouceMatch(maps []models.PolicyRuleCombinedActions, searchTerm string) []models.PolicyRuleCombinedActions {
	var result []models.PolicyRuleCombinedActions

	// Iterate through map of policies
	for _, m := range maps {
		// If success and resource contains search term
		if m.Resource == searchTerm {
			result = append(result, m)
		}
	}

	return result
}

func ApplyNamingConventionToRoleInheritanceRecord(inherit *models.GRecord) {
	inherit.Role = "role:" + inherit.Role
	inherit.InheritsFrom = "role:" + inherit.InheritsFrom
}
