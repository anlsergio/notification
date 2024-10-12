package dto

import "errors"

var (
	// ErrFailedValidation is used when the validation check has failed.
	ErrFailedValidation = errors.New("failed validation")
)
