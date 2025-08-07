package utils

import (
	"fmt"
	"strings"

	HttpException "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/httpexception"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func FormatValidationErrors(err error) *HttpException.ApiException {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string

		for _, fieldErr := range validationErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' failed validation: %s", fieldErr.Field(), fieldErr.Tag()))
		}

		return HttpException.NewBadRequestException("Validation failed", strings.Join(errorMessages, ", "))
	}

	return HttpException.NewInternalServerException(err)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
