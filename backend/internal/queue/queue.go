package queue

import (
	"sync"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/email"
	"gorm.io/gorm"
)

// Queue represents a job queue backed by a SQL database.
type Queue struct {
	db          *gorm.DB   // Database connection
	mu          sync.Mutex // Mutex for synchronizing access
	cond        *sync.Cond // Condition variable for signaling
	mailService email.Email
}

// Class method for creating a new job queue
// Backed by the given database (uses the job table).
func NewQueue(db *gorm.DB, mailService email.Email) *Queue {
	// Create the queue
	q := &Queue{
		db:          db,
		mailService: mailService,
	}
	// Initialize the mutex and condition variable
	q.cond = sync.NewCond(&q.mu)
	return q
}

// AddJob adds a new job to the queue.
// The jobType is a string that identifies the type of job.
// The payload is a string that contains the job data.
func (q *Queue) AddJob(jobType, payload string) error {
	// Lock the queue
	q.mu.Lock()
	// Unlock the queue when the function returns
	defer q.mu.Unlock()

	// Create a new job
	job := db.Job{
		JobType: jobType,
		Payload: payload,
	}

	// Store the job in the database
	if err := q.db.Create(&job).Error; err != nil {
		return err
	}
	q.cond.Signal() // Signal any waiting workers that a job is available
	return nil
}

// GetJob retrieves the next unprocessed job from the queue.
func (q *Queue) GetJob() (*db.Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var job db.Job
	// Fetch the first unprocessed job from the database
	if err := q.db.Where("processed = ?", false).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// MarkJobAsProcessed marks a job as processed in the database.
func (q *Queue) MarkJobAsProcessed(job *db.Job) error {
	// Lock the queue
	q.mu.Lock()
	defer q.mu.Unlock()

	// Mark the job as processed
	job.Processed = true
	job.Status = "processed"
	// Update the job in the database, returning any error
	return q.db.Save(job).Error
}
