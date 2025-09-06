package httpexception

import (
	"net/http"
)

type HTTPException struct {
	StatusCode int
	Message    string
	Errors     interface{}
}

func (e *HTTPException) Error() string {
	return e.Message
}

func (e *HTTPException) ToResponse() map[string]interface{} {
	response := map[string]interface{}{
		"success":    false,
		"statusCode": e.StatusCode,
		"message":    e.Message,
	}

	if e.Errors != nil {
		response["errors"] = e.Errors
	}

	return response
}

func NewBadRequestException(message string, err interface{}) *HTTPException {
	return &HTTPException{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Errors:     err,
	}
}

func NewConflictException(message string, err interface{}) *HTTPException {
	return &HTTPException{
		StatusCode: http.StatusConflict,
		Message:    message,
		Errors:     err,
	}
}

func NewNotFoundException(message string, err interface{}) *HTTPException {
	return &HTTPException{
		StatusCode: http.StatusNotFound,
		Message:    message,
		Errors:     err,
	}
}

func NewUnauthorizedException(message string, err interface{}) *HTTPException {
	return &HTTPException{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		Errors:     err,
	}
}

func NewInternalServerException(err interface{}) *HTTPException {
	return &HTTPException{
		StatusCode: http.StatusInternalServerError,
		Message:    "An unexpected error occurred...",
		Errors:     err,
	}
}
