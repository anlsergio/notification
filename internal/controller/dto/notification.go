package dto

import (
	"errors"
)

// Notification is the Data Transfer Object for HTTP handler operations.
type Notification struct {
	// CorrelationID is the ID used for correlating the notification message with different microservices,
	// including the sender microservice, which means it's treated as a universally unique identifier.
	CorrelationID string `json:"correlationId"`
	// UserID is the ID corresponding to the user the notification is meant to be sent to.
	UserID string `json:"userId"`
	// Type is the notification type.
	Type string `json:"type"`
	// Message is the message content of the notification.
	Message string `json:"message"`
}

// Validate returns an error ErrFailedValidation if Notification
// doesn't pass schema validation.
func (n Notification) Validate() error {
	var err error

	if n.CorrelationID == "" {
		err = errors.Join(ErrFailedValidation, errors.New("missing correlation ID"))
	}

	if n.UserID == "" {
		err = errors.Join(err, ErrFailedValidation, errors.New("user id is empty"))
	}

	if n.Type == "" {
		err = errors.Join(err, ErrFailedValidation, errors.New("type is empty"))
	}

	if n.Message == "" {
		err = errors.Join(err, ErrFailedValidation, errors.New("message is empty"))
	}

	return err
}
