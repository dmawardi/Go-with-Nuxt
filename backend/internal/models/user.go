package models

import (
	"fmt"
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type LoginResponse struct {
	Token string `json:"token"`
}

type ChangePassword struct {
	CurrentPassword    string `json:"current_password" valid:"length(6|30),required"`
	NewPassword        string `json:"new_password" valid:"length(6|30),required"`
	ConfirmNewPassword string `json:"confirm_new_password" valid:"length(6|30),required"`
}

// Users
// Create User structure for Data transfer.
type CreateUser struct {
	Username string `json:"username" valid:"length(6|25),required"`
	Password string `json:"password" valid:"length(6|30),required"`
	Name     string `json:"name" valid:"length(6|80),required"`
	Email    string `json:"email" valid:"email,required"`
	Verified bool   `json:"verified,omitempty"`
	Role     string `json:"role,omitempty" valid:""`
}

// Update User structure for Data transfer.
type UpdateUser struct {
	Username string `json:"username,omitempty" valid:"length(6|25)"`
	Password string `json:"password,omitempty" valid:"length(6|30)"`
	Name     string `json:"name,omitempty" valid:"length(6|80)"`
	Email    string `json:"email,omitempty" valid:"email"`
	Verified bool   `json:"verified,omitempty"`
	Role     string `json:"role,omitempty" valid:""`
}

type ResetPasswordAndEmailVerification struct {
	Email string `json:"email" valid:"email,required"`
}

type PaginatedUsers struct {
	Data *[]db.User     `json:"data"`
	Meta SchemaMetaData `json:"meta"`
}

type PaginatedUsersWithRole struct {
	Data *[]UserWithRole `json:"data"`
	Meta SchemaMetaData  `json:"meta"`
}

// Receiver functions for UserWithRole
func (schemaObject UserWithRole) GetID() string {
	return fmt.Sprint(schemaObject.ID)
}

type UserWithRole struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `swaggertype:"string" json:"created_at,omitempty"`
	UpdatedAt time.Time      `swaggertype:"string" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name,omitempty"`
	Username  string         `json:"username,omitempty"`
	Email     string         `json:"email,omitempty"`
	Password  string         `json:"-"`
	Role      string         `json:"role,omitempty"`
	// Verification
	Verified               *bool     `json:"verified,omitempty" gorm:"default:false"`
	VerificationCode       string    `json:"verification_code,omitempty" gorm:"default:null"`
	VerificationCodeExpiry time.Time `json:"verification_code_expiry,omitempty" gorm:"default:null"`
}

func (schemaObject UserWithRole) ObtainValue(keyValue string) string {
	fieldMap := map[string]string{
		"ID":                     fmt.Sprint(schemaObject.ID),
		"CreatedAt":              schemaObject.CreatedAt.Format(time.RFC3339),
		"UpdatedAt":              schemaObject.UpdatedAt.Format(time.RFC3339),
		"Name":                   schemaObject.Name,
		"Username":               schemaObject.Username,
		"Email":                  schemaObject.Email,
		"Verified":               fmt.Sprint(db.PointerToStringWithType(schemaObject.Verified, "bool")),
		"VerificationCode":       schemaObject.VerificationCode,
		"VerificationCodeExpiry": schemaObject.VerificationCodeExpiry.Format(time.RFC3339),
		"Role":                   schemaObject.Role,
	}
	// Return value of key
	return fieldMap[keyValue]
}
