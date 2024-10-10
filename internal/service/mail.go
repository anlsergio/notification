package service

// Mailer is the abstraction layer of the external email service integration itself.
type Mailer interface {
	// SendEmail sends the email message through the appropriate external service integration.
	SendEmail(to string, subject string, msg string) error
}
