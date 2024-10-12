package dto_test

import (
	"github.com/stretchr/testify/assert"
	"notification/internal/controller/dto"
	"testing"
)

func TestNotification_Validate(t *testing.T) {
	tests := []struct {
		name         string
		notification dto.Notification
		wantErr      error
	}{
		{
			name: "valid",
			notification: dto.Notification{
				UserID:  "123-abc",
				Type:    "marketing",
				Message: "Hey there!",
			},
			wantErr: nil,
		},
		{
			name: "missing user ID",
			notification: dto.Notification{
				Type:    "marketing",
				Message: "Hey there!",
			},
			wantErr: dto.ErrFailedValidation,
		},
		{
			name: "missing type",
			notification: dto.Notification{
				UserID:  "123-abc",
				Message: "Hey there!",
			},
			wantErr: dto.ErrFailedValidation,
		},
		{
			name: "missing message",
			notification: dto.Notification{
				UserID: "123-abc",
				Type:   "marketing",
			},
			wantErr: dto.ErrFailedValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.notification.Validate()
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
