package db

import (
	"fmt"
	"time"
)

// DB Schema interface implementation (Used for Admin Panel)
// Mapping of field names to values to allow for dynamic access
func (schemaObject User) ObtainValue(keyValue string) string {
	fieldMap := map[string]string{
		"ID":                     fmt.Sprint(schemaObject.ID),
		"CreatedAt":              schemaObject.CreatedAt.Format(time.RFC3339),
		"UpdatedAt":              schemaObject.UpdatedAt.Format(time.RFC3339),
		"Name":                   schemaObject.Name,
		"Username":               schemaObject.Username,
		"Email":                  schemaObject.Email,
		"Verified":               fmt.Sprint(PointerToStringWithType(schemaObject.Verified, "bool")),
		"VerificationCode":       schemaObject.VerificationCode,
		"VerificationCodeExpiry": schemaObject.VerificationCodeExpiry.Format(time.RFC3339),
	}
	// Return value of key
	return fieldMap[keyValue]
}

// Grabs the ID of the schema object as string
func (schemaObject User) GetID() string {
	return fmt.Sprint(schemaObject.ID)
}
