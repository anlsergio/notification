package infra

import (
	"fmt"
	"net/smtp"
)

// NewSMTPMailer instantiates a new SMTPMailer.
func NewSMTPMailer(address, from string, opts ...SMTPMailerOption) *SMTPMailer {
	mailer := SMTPMailer{
		address: address,
		from:    from,
	}

	for _, opt := range opts {
		opt(&mailer)
	}

	return &mailer
}

// SMTPMailer defines the SMTP Mailer implementation.
type SMTPMailer struct {
	address string
	from    string
	auth    smtp.Auth
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

	var auth smtp.Auth
	if m.auth != nil {
		auth = m.auth
	}

	return smtp.SendMail(m.address, auth, m.from, []string{to}, composedMsg)
}

// SMTPMailerOption defines the optional params for SMTPMailer.
type SMTPMailerOption func(*SMTPMailer)

// WithAuth optionally adds authentication capabilities to the mail sending mechanism.
// It's basically a wrapper for smtp.PlainAuth so refer to its documentation as reference on
// how to configure.
func WithAuth(identity, username, password, host string) SMTPMailerOption {
	return func(mailer *SMTPMailer) {
		mailer.auth = smtp.PlainAuth(identity, username, password, host)
	}
}
