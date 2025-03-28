package services

import (
	"monolith-domain/internal/mailer/domain"
	"monolith-domain/internal/mailer/infrastructure"
)

// MailerService defines the email sending operations.
type MailerService struct {
	Mailer *infrastructure.SMTPMailer
}

func NewMailerService(mailer *infrastructure.SMTPMailer) *MailerService {
	return &MailerService{Mailer: mailer}
}

// SendMail sends a single email.
func (m *MailerService) SendMail(toEmail, subject, body string) error {
	if !domain.ValidateEmail(toEmail) {
		return domain.ErrInvalidEmail
	}
	return m.Mailer.Send(toEmail, subject, body)
}

// SendBulkEmails sends emails to multiple recipients.
func (m *MailerService) SendBulkEmails(recipients []string, subject, body string) error {
	for _, email := range recipients {
		if !domain.ValidateEmail(email) {
			return domain.ErrInvalidEmail
		}
		err := m.Mailer.Send(email, subject, body)
		if err != nil {
			return err
		}
	}
	return nil
}
