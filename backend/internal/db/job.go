package db

import (
	"time"

	"gorm.io/gorm"
)

// Job (used for async jobs)
type Job struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `swaggertype:"string" json:"created_at,omitempty"`
	UpdatedAt time.Time      `swaggertype:"string" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	JobType   string         // Type of job (e.g., "email", "otherType")
	// Status
	Status string `json:"status,omitempty" gorm:"default:'pending'"`
	// Payload
	Payload   string `json:"payload,omitempty"`
	Processed bool
}

// Used prior to job creation
func (job *Job) BeforeCreate(tx *gorm.DB) (err error) {
	job.CreatedAt = time.Now()
	return
}
