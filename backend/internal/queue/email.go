package queue

import (
	"encoding/json"
)

// EmailJobPayload defines the structure of the email job payload
type EmailJobPayload struct {
	Recipient string
	Subject   string
	Body      string
}

// ProcessEmailJob processes an email job
func (q *Queue) ProcessEmailJob(payload string) error {
	var emailPayload EmailJobPayload
	// Unmarshal the payload into the email payload struct
	if err := json.Unmarshal([]byte(payload), &emailPayload); err != nil {
		return err
	}

	// Use the mail service to send the email
	return q.mailService.SendEmail(emailPayload.Recipient, emailPayload.Subject, emailPayload.Body)
}
