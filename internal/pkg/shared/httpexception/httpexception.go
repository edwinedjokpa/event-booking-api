package httpexception

import (
	"net/http"
)

type HTTPException struct {
	StatusCode int
	Message    string
	Details    interface{}
}

func (e *HTTPException) Error() string {
	return e.Message
}

func (e *HTTPException) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"message": e.Message,
		"details": e.Details,
	}
}

func NewBadRequestException(message string, err interface{}) *HTTPException {
	return &HTTPException{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Details:    err,
	}
}

func NewConflictException(message string, err interface{}) *HTTPException {
	return &HTTPException{
		StatusCode: http.StatusConflict,
		Message:    message,
		Details:    err,
	}
}

func NewNotFoundException(message string, err interface{}) *HTTPException {
	return &HTTPException{
		StatusCode: http.StatusNotFound,
		Message:    message,
		Details:    err,
	}
}

func NewUnauthorizedException(message string, err interface{}) *HTTPException {
	return &HTTPException{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		Details:    err,
	}
}

func NewInternalServerException(err interface{}) *HTTPException {
	return &HTTPException{
		StatusCode: http.StatusInternalServerError,
		// Message:    "Internal Server Error",
		Message: "An unexpected error occurred...",
		Details: err,
	}
}
