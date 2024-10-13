package config_test

import (
	"github.com/stretchr/testify/assert"
	"notification/internal/config"
	"os"
	"testing"
)

func TestNewAppConfig(t *testing.T) {
	t.Run("server port is populated", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "8081")
		defer os.Unsetenv("SERVER_PORT")

		cfg := config.NewAppConfig()

		assert.Equal(t, 8081, cfg.ServerPort)
	})
	t.Run("server port defaults to 8080", func(t *testing.T) {
		cfg := config.NewAppConfig()
		assert.Equal(t, 8080, cfg.ServerPort)
	})
}
