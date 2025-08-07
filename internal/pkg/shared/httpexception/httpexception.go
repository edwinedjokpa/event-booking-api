package httpexception

import (
	"net/http"
)

type ApiException struct {
	StatusCode int
	Message    string
	Details    interface{}
}

func (e *ApiException) Error() string {
	return e.Message
}

func (e *ApiException) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"message": e.Message,
		"details": e.Details,
	}
}

func NewBadRequestException(message string, err interface{}) *ApiException {
	return &ApiException{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Details:    err,
	}
}

func NewConflictException(message string, err interface{}) *ApiException {
	return &ApiException{
		StatusCode: http.StatusConflict,
		Message:    message,
		Details:    err,
	}
}

func NewNotFoundException(message string, err interface{}) *ApiException {
	return &ApiException{
		StatusCode: http.StatusNotFound,
		Message:    message,
		Details:    err,
	}
}

func NewUnauthorizedException(message string, err interface{}) *ApiException {
	return &ApiException{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		Details:    err,
	}
}

func NewInternalServerException(err interface{}) *ApiException {
	return &ApiException{
		StatusCode: http.StatusInternalServerError,
		// Message:    "Internal Server Error",
		Message: "An unexpected error occurred...",
		Details: err,
	}
}
