package infrastructure

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

// SMTPMailer handles email sending via SMTP.
type SMTPMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// NewSMTPMailer initializes an SMTPMailer using environment variables.
func NewSMTPMailer() (*SMTPMailer, error) {
	godotenv.Load()

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" || username == "" || password == "" || from == "" {
		return nil, fmt.Errorf("missing required SMTP environment variables")
	}

	return &SMTPMailer{host, port, username, password, from}, nil
}

// Send sends an email using SMTP.
func (s *SMTPMailer) Send(toEmail, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", s.from, toEmail, subject, body))

	err := smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
