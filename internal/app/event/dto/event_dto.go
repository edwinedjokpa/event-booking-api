package dto

import (
	"time"
)

type CreateEventRequest struct {
	Name        string    `json:"name" validate:"required,min=3"`
	Description string    `json:"description" validate:"required,min=10"`
	Location    string    `json:"location" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}

type UpdateEventRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Location    *string    `json:"location"`
	Date        *time.Time `json:"date"`
}
