package infrastructure

import (
	"fmt"
	"net/smtp"
	"monolith-domain/internal/mailer/domain"
	"go.uber.org/zap"
	"monolith-domain/pkg/observability"
)

// SMTPMailer handles email sending via SMTP.
type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
	from     string
	secure   bool
	logger   *zap.Logger
}

// NewSMTPMailer initializes a new SMTPMailer with the provided configuration.
func NewSMTPMailer(host string, port int, username, password, from string, secure bool) (*SMTPMailer, error) {
	return &SMTPMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		secure:   secure,
		logger:   observability.GetLogger(),
	}, nil
}

// Send sends an email using SMTP with the provided parameters.
// It logs the attempt and any errors that occur during the process.
func (m *SMTPMailer) Send(mail domain.Mail) error {
	m.logger.Info("Attempting to send email",
		zap.String("to", mail.To),
		zap.String("subject", mail.Subject),
		zap.String("from", m.from),
		zap.String("host", m.host),
		zap.Int("port", m.port),
	)

	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	headers := make(map[string]string)
	headers["From"] = m.from
	headers["To"] = mail.To
	headers["Subject"] = mail.Subject

	if mail.IsHTML {
		headers["MIME-Version"] = "1.0"
		headers["Content-Type"] = "text/html; charset=UTF-8"
	} else {
		headers["Content-Type"] = "text/plain; charset=UTF-8"
	}

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + mail.Body

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	err := smtp.SendMail(addr, auth, m.from, []string{mail.To}, []byte(message))
	if err != nil {
		m.logger.Error("Failed to send email",
			zap.String("to", mail.To),
			zap.String("subject", mail.Subject),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	m.logger.Info("Email sent successfully",
		zap.String("to", mail.To),
		zap.String("subject", mail.Subject),
	)

	return nil
}
