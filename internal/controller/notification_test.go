package controller_test

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"notification/internal/controller"
	"notification/internal/domain"
	"notification/internal/repository"
	"notification/internal/service"
	"notification/mocks"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNotification(t *testing.T) {
	t.Run("notification sending", func(t *testing.T) {
		t.Run("operation is successful", func(t *testing.T) {
			svc := mocks.NewNotificationSender(t)
			svc.
				On("Send", mock.Anything, "abc-123", "Hey there!", domain.Marketing).
				Return(time.Duration(0), nil)

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
				Return(time.Duration(0), errors.New("oops"))

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
				Return(time.Duration(0), nil).
				Maybe()

			notificationController := controller.NewNotification(svc)

			r := mux.NewRouter()
			notificationController.SetRouter(r)

			req := httptest.NewRequest(http.MethodPost, "/send", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("HTTP status is Bad Request", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			})

			t.Run("service send isn't called", func(t *testing.T) {
				svc.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
		})

		t.Run("fail to pass schema validation", func(t *testing.T) {
			svc := mocks.NewNotificationSender(t)
			svc.
				On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(time.Duration(0), nil).
				Maybe()

			notificationController := controller.NewNotification(svc)

			r := mux.NewRouter()
			notificationController.SetRouter(r)

			requestBody := `
{
	"user_id": "",
	"type": "",
	"message": ""
}
`

			req := httptest.NewRequest(http.MethodPost, "/send", strings.NewReader(requestBody))
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("HTTP status is Bad Request", func(t *testing.T) {
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
				Return(time.Duration(0), nil).
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

			t.Run("HTTP status is Bad Request", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			})

			t.Run("service send isn't called", func(t *testing.T) {
				svc.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
		})

		t.Run("invalid user ID", func(t *testing.T) {
			svc := mocks.NewNotificationSender(t)
			svc.
				On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(time.Duration(0), fmt.Errorf("oops: %w", repository.ErrInvalidUserID))

			notificationController := controller.NewNotification(svc)

			r := mux.NewRouter()
			notificationController.SetRouter(r)

			requestBody := `
{
	"user_id": "abc-123",
	"type": "status",
	"message": "Hey there!"
}
`

			req := httptest.NewRequest(http.MethodPost, "/send", strings.NewReader(requestBody))
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("HTTP status is Bad Request", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			})
		})

		t.Run("rate limit exceeded", func(t *testing.T) {
			retryAfter := time.Minute

			svc := mocks.NewNotificationSender(t)
			svc.
				On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(retryAfter, fmt.Errorf("oops: %w", service.ErrRateLimitExceeded))

			notificationController := controller.NewNotification(svc)

			r := mux.NewRouter()
			notificationController.SetRouter(r)

			requestBody := `
{
	"user_id": "abc-123",
	"type": "status",
	"message": "Hey there!"
}
`

			req := httptest.NewRequest(http.MethodPost, "/send", strings.NewReader(requestBody))
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("HTTP status is Too Many Requests", func(t *testing.T) {
				assert.Equal(t, http.StatusTooManyRequests, rr.Code)
			})
			t.Run("Retry-After header is informed", func(t *testing.T) {
				retryAfterHeader := rr.Header().Get("Retry-After")
				require.NotEmpty(t, retryAfterHeader)
				assert.Equal(t, strconv.Itoa(int(retryAfter.Seconds())), retryAfterHeader)
			})
		})
	})
}
