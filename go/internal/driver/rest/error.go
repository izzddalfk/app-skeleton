package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/izzdalfk/app-skeleton/go/internal/core/service"
)

// ADJUST: feel free to adjust this struct following your api specification
type RestError struct {
	StatusCode int
	Err        string
	Message    string
}

func (e *RestError) Error() string {
	return fmt.Sprintf("%d: %v - %v", e.StatusCode, e.Message, e.Err)
}

func (e *RestError) Is(target error) bool {
	var restErr *RestError
	if !errors.As(target, &restErr) {
		return false
	}

	return *e == *restErr
}

func NewBadRequestError(msg string) *RestError {
	return &RestError{
		StatusCode: http.StatusBadRequest,
		Err:        "ERR_BAD_REQUEST",
		Message:    msg,
	}
}

func NewInvalidAccessTokenError() *RestError {
	return &RestError{
		StatusCode: http.StatusUnauthorized,
		Err:        "ERR_INVALID_ACCESS_TOKEN",
		Message:    "invalid access token",
	}
}

func NewForbiddenError(msg string) *RestError {
	return &RestError{
		StatusCode: http.StatusForbidden,
		Err:        "ERR_FORBIDDEN_ACCESS",
		Message:    msg,
	}

}

func NewNotFoundError(msg string) *RestError {
	return &RestError{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_NOT_FOUND",
		Message:    msg,
	}
}

func NewInternalServerError(msg string) *RestError {
	return &RestError{
		StatusCode: http.StatusInternalServerError,
		Err:        "ERR_INTERNAL_ERROR",
		Message:    msg,
	}
}

func readErr(err error) *RestError {
	// ADJUST: treat specific errors here to be returned as proper error response
	switch {
	case errors.Is(err, service.ErrWrongEntity):
		return NewBadRequestError(err.Error())
	}

	return NewInternalServerError(err.Error())
}
