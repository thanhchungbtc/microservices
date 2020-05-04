package app

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type errorResponse struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func newErrorResponse(msg string) errorResponse {
	return errorResponse{Message: msg}
}

type Error interface {
	StatusCode() int
	Error() string
	Json() []errorResponse
}

type ErrDatabaseConnection struct {
	error
}

func (e ErrDatabaseConnection) StatusCode() int {
	return 400
}

func (e ErrDatabaseConnection) Json() []errorResponse {
	return []errorResponse{newErrorResponse(e.Error())}
}

// ErrNotFound
type ErrNotFound struct {
	error
}

func (ErrNotFound) StatusCode() int {
	return 404
}

func (e ErrNotFound) Json() []errorResponse {
	return []errorResponse{newErrorResponse(e.Error())}
}

// ErrValidation
type ErrValidation struct {
	error
}

func (ErrValidation) StatusCode() int {
	return 400
}
func (v ErrValidation) Json() (errResponses []errorResponse) {
	switch err := v.error.(type) {
	case validator.ValidationErrors:
		for _, e := range err {
			errResponses = append(errResponses, errorResponse{
				Message: fmt.Sprintf("error validation for field %s with tag %s", e.Field(), e.Tag()),
				Field:   e.Field(),
			})
		}
	default:
		errResponses = append(errResponses, errorResponse{
			Message: "Invalid data",
		})
	}
	return errResponses
}
