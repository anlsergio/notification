package config

import (
	"os"
	"strconv"
)

// NewAppConfig loads the application configuration parameters
// and returns an instance of it.
func NewAppConfig() *AppConfig {
	var cfg AppConfig

	var err error
	cfg.ServerPort, err = strconv.Atoi(os.Getenv("SERVE_PORT"))
	if err != nil || cfg.ServerPort == 0 {
		cfg.ServerPort = 8080
	}

	cfg.MailFrom = os.Getenv("MAIL_FROM")
	cfg.SMTPHost = os.Getenv("SMTP_HOST")
	if cfg.SMTPHost == "" {
		cfg.SMTPHost = "localhost"
	}

	cfg.SMTPPort, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil || cfg.SMTPPort == 0 {
		cfg.SMTPPort = 587
	}

	return &cfg
}

// AppConfig represents the application configuration params.
type AppConfig struct {
	// ServerPort is the port where the API server will
	// listen for connections. Defaults to 8080.
	ServerPort int
	Mail
}

// Mail represents the mail configuration params.
type Mail struct {
	// MailFrom configures the mail from address of the notification messages.
	MailFrom string
	// SMTPHost is the host for SMTP connection. Defaults to localhost.
	SMTPHost string
	// SMTPPort is the port for SMTP connection. Defaults to 587.
	SMTPPort int
}
