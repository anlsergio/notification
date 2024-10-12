package controller_test

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"notification/internal/controller"
	"notification/internal/domain"
	"notification/mocks"
	"strings"
	"testing"
)

func TestNotification(t *testing.T) {
	t.Run("notification sending", func(t *testing.T) {
		t.Run("operation is successful", func(t *testing.T) {
			svc := mocks.NewNotificationSender(t)
			svc.
				On("Send", mock.Anything, "abc-123", "Hey there!", domain.Marketing).
				Return(nil)

			notificationController := controller.NewNotification(svc)

			r := mux.NewRouter()
			notificationController.SetRouter(r)

			requestBody := `
{
	"user_id": "abc-123",
	"type": "marketing",
	"message": "Hey there!"
}
`

			req := httptest.NewRequest(http.MethodPost, "/send", strings.NewReader(requestBody))
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("HTTP status is OK", func(t *testing.T) {
				assert.Equal(t, http.StatusOK, rr.Code)
			})
		})

		t.Run("service errors out", func(t *testing.T) {
			svc := mocks.NewNotificationSender(t)
			svc.
				On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(errors.New("oops"))

			notificationController := controller.NewNotification(svc)

			r := mux.NewRouter()
			notificationController.SetRouter(r)

			requestBody := `
{
	"user_id": "abc-123",
	"type": "marketing",
	"message": "Hey there!"
}
`

			req := httptest.NewRequest(http.MethodPost, "/send", strings.NewReader(requestBody))
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("HTTP status is Internal Server Error", func(t *testing.T) {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
			})
		})

		t.Run("fail to parse request body", func(t *testing.T) {
			svc := mocks.NewNotificationSender(t)
			svc.
				On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(nil).
				Maybe()

			notificationController := controller.NewNotification(svc)

			r := mux.NewRouter()
			notificationController.SetRouter(r)

			req := httptest.NewRequest(http.MethodPost, "/send", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("HTTP status is Internal Server Error", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			})

			t.Run("service send isn't called", func(t *testing.T) {
				svc.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
		})

		t.Run("invalid notification type", func(t *testing.T) {
			svc := mocks.NewNotificationSender(t)
			svc.
				On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(nil).
				Maybe()

			notificationController := controller.NewNotification(svc)

			r := mux.NewRouter()
			notificationController.SetRouter(r)

			requestBody := `
{
	"user_id": "abc-123",
	"type": "invalid",
	"message": "Hey there!"
}
`

			req := httptest.NewRequest(http.MethodPost, "/send", strings.NewReader(requestBody))
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("HTTP status is Internal Server Error", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			})

			t.Run("service send isn't called", func(t *testing.T) {
				svc.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
		})
	})
}
