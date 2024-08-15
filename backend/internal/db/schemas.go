package db

import (
	"time"

	"gorm.io/gorm"
)

var Models = []interface{}{
	// Core Schemas
	&User{}, // Used for user management
	&Job{},  // Used for job queuing
	&Action{}, // Used for logging actions
	// Additional Schemas
	&Post{},
}

// Base Schemas
type User struct {
	// gorm.Model `json:"-"`
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `swaggertype:"string" json:"created_at,omitempty"`
	UpdatedAt time.Time      `swaggertype:"string" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name,omitempty"`
	Username  string         `json:"username,omitempty"`
	Email     string         `json:"email,omitempty" gorm:"uniqueIndex"`
	Password  string         `json:"-"`
	// Verification
	Verified               *bool     `json:"verified,omitempty" gorm:"default:false"`
	VerificationCode       string    `json:"verification_code,omitempty" gorm:"default:null"`
	VerificationCodeExpiry time.Time `json:"verification_code_expiry,omitempty" gorm:"default:null"`
	// Relationships
	Posts []Post `json:"posts,omitempty" gorm:"foreignKey:UserID"`
}

// Additional Schemas
type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `swaggertype:"string" json:"created_at,omitempty"`
	UpdatedAt time.Time      `swaggertype:"string" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Title     string         `json:"title,omitempty"`
	Body      string         `json:"body,omitempty"`
	UserID    uint           `json:"user_id,omitempty"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
