package seed

import (
	"github.com/dmawardi/Go-Template/internal/seed/factory"
	"github.com/dmawardi/Go-Template/internal/seed/schemas"
	"gorm.io/gorm"
)

func factoryRegister(db *gorm.DB) []schemas.FactoryRegistration {
	// Define the items to seed
	return []schemas.FactoryRegistration{
		// Add the list of factories here, with the name of the factory and the factory itself
		// This will be accessible in the factoryMap eg. factoryMap["Name"]
		// UserFactory
		{
			Factory: factory.NewUserFactory(db),
			Name:    "User",
		},
	}
}
