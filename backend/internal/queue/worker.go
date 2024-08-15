package queue

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// Worker processes jobs from the queue.
func (q *Queue) Worker() {
	for {
		// Get the next job
		job, err := q.GetJob()
		if err != nil {
			// If there's another error aside from "record not found", log it
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("Worker: Error getting job: %v\n", err)
			}
			// Wait for a signal that a job is available
			time.Sleep(5 * time.Second)
			continue
		}
		// Process the job using the Process function with the payload
		if err := q.ProcessJob(job.JobType, job.Payload); err != nil {
			log.Printf("Worker: Error processing job: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		// Mark the job as processed
		if err := q.MarkJobAsProcessed(job); err != nil {
			log.Printf("Worker: Error marking job as processed: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

func (q *Queue) ProcessJob(jobType, payload string) error {
	switch jobType {
	case "email":
		return q.ProcessEmailJob(payload)
	default:
		return errors.New("unknown job type")
	}
}
