package infra

import (
	"fmt"
	"net/smtp"
)

// NewSMTPMailer instantiates a new SMTPMailer.
func NewSMTPMailer(address, from string) *SMTPMailer {
	// TODO: add auth capabilities.
	//auth := smtp.PlainAuth("", "mailtrap.foo@gmail.com", "your_password", "live.smtp.mailtrap.io")
	return &SMTPMailer{
		address: address,
		from:    from,
	}
}

// SMTPMailer defines the SMTP Mailer implementation.
type SMTPMailer struct {
	address string
	from    string
}

// SendEmail sends the email message through SMTP integration.
func (m SMTPMailer) SendEmail(to string, subject string, msg string) error {
	composedMsg := []byte(fmt.Sprintf(
		`
To: %s
Subject: %s

%s
`, to, subject, msg,
	))
	return smtp.SendMail(m.address, nil, m.from, []string{to}, composedMsg)
}
