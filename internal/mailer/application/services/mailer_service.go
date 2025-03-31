package services

import (
	"monolith-domain/internal/mailer/domain"
)

// MailerService defines the email sending operations.
type MailerService struct {
	repo domain.MailerRepository
}

func NewMailerService(repo domain.MailerRepository) *MailerService {
	return &MailerService{
		repo: repo,
	}
}

// SendMail sends a single email.
func (s *MailerService) SendMail(to string, subject string, body string, isHTML bool) error {
	mail := domain.Mail{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  isHTML,
	}
	return s.repo.Send(mail)
}

// SendBulkEmails sends emails to multiple recipients.
func (s *MailerService) SendBulkEmails(recipients []string, subject, body string, isHTML bool) error {
	for _, email := range recipients {
		mail := domain.Mail{
			To:      email,
			Subject: subject,
			Body:    body,
			IsHTML:  isHTML,
		}
		if err := s.repo.Send(mail); err != nil {
			return err
		}
	}
	return nil
}
