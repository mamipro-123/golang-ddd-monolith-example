package domain

type Mail struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

type MailerRepository interface {
	Send(mail Mail) error
}

type MailerService interface {
	SendMail(to string, subject string, body string, isHTML bool) error
} 