package webapi

import "github.com/dmawardi/Go-Template/internal/models"

func mergeModuleMaps(map1, map2 map[string]models.ModuleSet) map[string]models.ModuleSet {
	// Create a new map to hold the merged results
	mergedMap := make(map[string]models.ModuleSet)

	// Add all entries from the first map
	for key, value := range map1 {
		mergedMap[key] = value
	}

	// Add all entries from the second map
	for key, value := range map2 {
		mergedMap[key] = value
	}

	return mergedMap
}