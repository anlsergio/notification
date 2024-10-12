package domain_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"notification/internal/domain"
	"testing"
)

func TestToNotificationType(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    domain.NotificationType
		wantErr error
	}{
		{
			"conversion is successful",
			"status",
			domain.Status,
			nil,
		},
		{
			"invalid type",
			"invalid",
			0,
			domain.ErrInvalidNotificationType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ToNotificationType(tt.in)
			require.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
