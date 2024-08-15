package seed

import (
	"github.com/dmawardi/Go-Template/internal/seed/schemas"
	"gorm.io/gorm"
)

// Boot is run at server boot to seed the database
func Boot(db *gorm.DB) error {
	// Prepare the factory map based off the list of factories in config.go
	factoryMap := PrepareSeedFactoryMap(db)

	// Seed the database with the seed function in seed.go
	seed(factoryMap)

	return nil
}

// Builds the factory map and prepares it for use by user
func PrepareSeedFactoryMap(db *gorm.DB) map[string]schemas.BasicFactory {
	// Build the list of factories
	factoryList := factoryRegister(db)

	// Build into a map of factories for easier usage
	factoryMap := buildFactoryMap(factoryList)

	return factoryMap
}

// Takes a list of FactoryRegistrations and builds a map of the factories for easy usage
func buildFactoryMap(factoryList []schemas.FactoryRegistration) map[string]schemas.BasicFactory {
	// Create a map to store the factories
	factoryMap := make(map[string]schemas.BasicFactory)

	// Loop through the list of factory registrations
	for _, record := range factoryList {
		// Add the factory to the map using the name as the key
		factoryMap[record.Name] = record.Factory
	}

	return factoryMap
}
