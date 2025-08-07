package event

import (
	"errors"
	"fmt"
	"time"

	EventDTO "github.com/edwinedjokpa/event-booking-api/internal/app/event/dto"
	HTTPException "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/httpexception"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/utils"
	"gorm.io/gorm"
)

type EventService interface {
	CreateEvent(userID string, request EventDTO.CreateEventRequest)
	GetAllEvents() []Event
	GetEventByID(eventID string) *Event
	UpdateEvent(userID, eventID string, request EventDTO.UpdateEventRequest)
	DeleteEvent(userID, eventID string)
}

type eventService struct {
	repository EventRepository
}

func NewEventService(repository EventRepository) EventService {
	return &eventService{repository}
}

func (svc *eventService) CreateEvent(userID string, request EventDTO.CreateEventRequest) {
	event := Event{
		ID:          utils.GenerateUUID(),
		Name:        request.Name,
		Description: request.Description,
		Location:    request.Location,
		Date:        request.Date,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := svc.repository.Create(event); err != nil {
		panic(HTTPException.NewBadRequestException("Failed to create event", nil))
	}
}

func (svc *eventService) GetAllEvents() []Event {
	allEvents, err := svc.repository.FindAll()
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to retrieve all events", nil))
	}

	return allEvents
}

func (svc *eventService) GetEventByID(eventID string) *Event {
	event, err := svc.repository.FindOneByID(eventID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}

	if event == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		panic(HTTPException.NewNotFoundException(fmt.Sprintf("Event with ID %s not found", eventID), nil))
	}

	return event
}

func (svc *eventService) UpdateEvent(userID, eventID string, request EventDTO.UpdateEventRequest) {
	existingEvent := svc.GetEventByID(eventID)

	if existingEvent.UserID != userID {
		panic(HTTPException.NewUnauthorizedException("You cannot update an event that was not created by you", nil))
	}

	if request.Name != nil {
		existingEvent.Name = *request.Name
	}

	if request.Description != nil {
		existingEvent.Description = *request.Description
	}

	if request.Location != nil {
		existingEvent.Location = *request.Location
	}

	if request.Date != nil {
		existingEvent.Date = *request.Date
	}

	existingEvent.UpdatedAt = time.Now()

	if err := svc.repository.Update(*existingEvent); err != nil {
		panic(HTTPException.NewBadRequestException("Failed to update event", err.Error()))
	}
}

func (svc *eventService) DeleteEvent(userID, eventID string) {
	event := svc.GetEventByID(eventID)

	if event.UserID != userID {
		panic(HTTPException.NewUnauthorizedException("You cannot delete an event that was not created by you", nil))
	}

	if err := svc.repository.Delete(eventID); err != nil {
		panic(HTTPException.NewBadRequestException("Failed to delete event", err.Error()))
	}
}
