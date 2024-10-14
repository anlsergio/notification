package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"notification/internal/controller/dto"
	"notification/internal/controller/middleware"
	"notification/internal/domain"
	"notification/internal/repository"
	"notification/internal/service"
)

// NewNotification creates a new Notification controller instance.
func NewNotification(svc service.NotificationSender) *Notification {
	return &Notification{svc}
}

// Notification is the notification controller.
// It defines routes and handlers for the notification resources.
type Notification struct {
	svc service.NotificationSender
}

// SetRouter returns the router r with all the necessary routes for the
// Notification controller setup.
func (n Notification) SetRouter(r *mux.Router) {
	r.HandleFunc("/send", middleware.Logger(middleware.SetJSONContent(n.send))).
		Methods(http.MethodPost)
}

// @Summary Send a notification message
// @Description Sends a notification message
// @Tags notification
// @Accept json
// @Produce json
// @Param notification body dto.Notification true "Notification object to be sent"
// @Success 200
// @Failure 400 {object} string "Bad Request"
// @Failure 409 {object} string "Conflict"
// @Failure 429 {object} string "Too Many Requests"
// @Failure 500 {object} string "Internal Server Error"
// @Header 429 {string} Retry-After "3600"
// @Router /send [post]
func (n Notification) send(w http.ResponseWriter, r *http.Request) {
	var notificationDTO dto.Notification

	if err := json.NewDecoder(r.Body).Decode(&notificationDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := notificationDTO.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notificationType, err := domain.ToNotificationType(notificationDTO.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notification := domain.Notification{
		CorrelationID: notificationDTO.CorrelationID,
		Type:          notificationType,
		Message:       notificationDTO.Message,
	}

	retryAfter, err := n.svc.Send(r.Context(), notificationDTO.UserID, notification)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidUserID):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, service.ErrRateLimitExceeded):
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(retryAfter.Seconds())))
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		case errors.Is(err, service.ErrIdempotencyViolation):
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
