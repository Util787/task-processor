package models

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
)

// ErrValidation is an abstraction that should be used only to get right status code in handlers
//
// Any validation error should contain this abstraction
var ErrValidation = errors.New("validation error")
