package controller

import (
	"encoding/json"
	"errors"
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
	r.HandleFunc("/send", middleware.SetJSONContent(n.send)).
		Methods(http.MethodPost)
}

// @Summary Send a notification message
// @Description Sends a notification message
// @Tags notification
// @Accept json
// @Produce json
// @Param notification body dto.Notification true "Notification object to be sent"
// @Success 200
// @Failure 400 {object} string "Error message"
// @Failure 500 {object} string "Error message"
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

	err = n.svc.Send(r.Context(), notificationDTO.UserID, notificationDTO.Message, notificationType)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidUserID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
