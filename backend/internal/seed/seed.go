package seed

import (
	"github.com/dmawardi/Go-Template/internal/seed/schemas"
)

// Seed creates dummy data and runs at server boot
// This is useful for testing purposes, as it allows you to test the application with a large amount of data without having to manually input it.
func seed(seed map[string]schemas.BasicFactory) {
	// This is where you would place seed commands to run at the start of the application.
	// Using the Name given in the factory registration, you can access the factory from the map
	//
	//
	// seed["User"].Factory(10)
	// fmt.Printf("Database seeded\n")
}
