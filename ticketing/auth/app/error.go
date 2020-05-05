package app

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

var (
	ErrDatabaseConnection = errors.New("failed to connect to database")
	ErrNotFound           = errors.New("not found")
)

type errorResponse struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type ErrBadRequest struct {
	error
}

func (v ErrBadRequest) Responses() (errResponses []errorResponse) {
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
			Message: err.Error(),
		})
	}
	return errResponses
}
