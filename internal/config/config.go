package config

import (
	"os"
	"strconv"
)

// NewAppConfig loads the application configuration parameters
// and returns an instance of it.
func NewAppConfig() *AppConfig {
	var cfg AppConfig
	cfg.HTTPServer.parseConfig()
	cfg.Mail.parseConfig()
	cfg.Redis.parseConfig()

	return &cfg
}

// AppConfig represents the application configuration params.
type AppConfig struct {
	HTTPServer
	Mail
	Redis
}

// HTTPServer represents the HTTP server configuration params.
type HTTPServer struct {
	// ServerPort is the port where the API server will
	// listen for connections. Defaults to 8080.
	ServerPort int
}

func (s *HTTPServer) parseConfig() {
	var err error
	s.ServerPort, err = strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil || s.ServerPort == 0 {
		s.ServerPort = 8080
	}
}

// Mail represents the mail configuration params.
type Mail struct {
	// MailFrom configures the mail from address of the notification messages.
	MailFrom string
	// SMTPHost is the host for SMTP connection. Defaults to localhost.
	SMTPHost string
	// SMTPPort is the port for SMTP connection. Defaults to 587.
	SMTPPort int
	// SMTPUsername is the username for SMTP authentication.
	SMTPUsername string
	// SMTPPassword is the password for SMTP authentication.
	SMTPPassword string
}

func (m *Mail) parseConfig() {
	m.MailFrom = os.Getenv("MAIL_FROM")
	m.SMTPHost = os.Getenv("SMTP_HOST")
	if m.SMTPHost == "" {
		m.SMTPHost = "localhost"
	}

	var err error
	m.SMTPPort, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil || m.SMTPPort == 0 {
		m.SMTPPort = 587
	}

	m.SMTPUsername = os.Getenv("SMTP_USERNAME")
	m.SMTPPassword = os.Getenv("SMTP_PASSWORD")
}

// Redis represents the Redis cache configuration params.
type Redis struct {
	// RedisHost is the host for Redis connection. Defaults to localhost.
	RedisHost string
	// RedisPort is the port for Redis connection. Defaults to 6379.
	RedisPort int
}

func (r *Redis) parseConfig() {
	r.RedisHost = os.Getenv("REDIS_HOST")
	if r.RedisHost == "" {
		r.RedisHost = "localhost"
	}
	var err error
	r.RedisPort, err = strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil || r.RedisPort == 0 {
		r.RedisPort = 6379
	}
}
