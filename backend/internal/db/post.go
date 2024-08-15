package db

import (
	"fmt"
	"time"
)

// DB Schema interface implementation (Used for Admin Panel)
// Mapping of field names to values to allow for dynamic access
func (schemaObject Post) ObtainValue(keyValue string) string {
	// Build foreign key string for easy admin panel digestion
	userForeignKeyString := fmt.Sprintf("[%s]:%s|", "ID", fmt.Sprint(schemaObject.User.ID))
	userForeignKeyString += fmt.Sprintf("[%s]:%s|", "Name", fmt.Sprint(schemaObject.User.Name))
	userForeignKeyString += fmt.Sprintf("[%s]:%s|", "Email", fmt.Sprint(schemaObject.User.Email))
	userForeignKeyString += fmt.Sprintf("[%s]:%s|", "Username", fmt.Sprint(schemaObject.User.Username))
	userForeignKeyString += fmt.Sprintf("[%s]:%s|", "Verified", fmt.Sprint(schemaObject.User.Verified))

	// Map of post fields
	fieldMap := map[string]string{
		"ID":        fmt.Sprint(schemaObject.ID),
		"CreatedAt": schemaObject.CreatedAt.Format(time.RFC3339),
		"UpdatedAt": schemaObject.UpdatedAt.Format(time.RFC3339),
		"Title":     schemaObject.Title,
		"Body":      schemaObject.Body,
		"User":      fmt.Sprintf("%+v", userForeignKeyString),
	}
	// Return value of key
	return fieldMap[keyValue]
}

// Grabs the ID of the schema object as string
func (schemaObject Post) GetID() string {
	return fmt.Sprint(schemaObject.ID)
}
