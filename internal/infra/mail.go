package infra

import (
	"net/smtp"
)

// NewSMTPMailer instantiates a new SMTPMailer.
func NewSMTPMailer(address, from string) *SMTPMailer {
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
func (m SMTPMailer) SendEmail(to []string, msg []byte) error {
	return smtp.SendMail(m.address, nil, m.from, to, msg)
}
