package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type Email interface {
	SendEmail(recipient, subject, body string) error
}

// Email struct
type email struct {
	Auth        smtp.Auth
	SmtpAddress string
	FromAddress string
	SendMail    func(from, recipient, body string) error
}

func NewSMTPEmail() Email {
	return &email{
		Auth:        smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST")),
		SmtpAddress: os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT"),
		FromAddress: os.Getenv("SMTP_USERNAME"),
	}
}

// Sends email using SMTP
func (e *email) SendEmail(recipient, subject, body string) error {
	// Set MIME and other headers
	headers := make(map[string]string)
	headers["From"] = e.FromAddress
	headers["To"] = recipient
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Build the headers
	header := ""
	for k, v := range headers {
		header += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// The msg parameter should be an RFC 822-style email with headers first,
	// a blank line, and then the message body (Go Template).
	msg := []byte(
		header + "\r\n" +
			body + "\r\n")

	// This sends the email with a plain auth setup
	err := smtp.SendMail(e.SmtpAddress, e.Auth, e.FromAddress, []string{recipient}, msg)
	if err != nil {
		return fmt.Errorf("smtp.SendMail() failed with: %s", err)
	}
	return nil
}
