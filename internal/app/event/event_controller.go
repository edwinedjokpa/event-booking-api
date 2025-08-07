package event

import (
	"net/http"

	EventDTO "github.com/edwinedjokpa/event-booking-api/internal/app/event/dto"
	APIResponse "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/apiresponse"
	HTTPException "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/httpexception"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type EventController interface {
	CreateEvent(c *gin.Context)
	GetAllEvents(c *gin.Context)
	GetEventByID(c *gin.Context)
	UpdateEvent(c *gin.Context)
	DeleteEvent(c *gin.Context)
}

type eventController struct {
	service   EventService
	validator *validator.Validate
}

func NewEventController(service EventService, validator *validator.Validate) EventController {
	return &eventController{service, validator}
}

func (ctrl *eventController) CreateEvent(ctx *gin.Context) {
	userIDRaw, exists := ctx.Get("userID")
	if !exists {
		exception := HTTPException.NewUnauthorizedException("Unauthorized", nil)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		exception := HTTPException.NewBadRequestException("User ID not found in context", nil)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	var request EventDTO.CreateEventRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := HTTPException.NewBadRequestException("Bad Request Exception", err.Error())
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	if err := ctrl.validator.Struct(request); err != nil {
		exception := utils.FormatValidationErrors(err)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	ctrl.service.CreateEvent(userID, request)
	ctx.JSON(http.StatusCreated, APIResponse.Success("Event created successfully", nil))
}

func (ctrl *eventController) GetAllEvents(ctx *gin.Context) {
	allEvents := ctrl.service.GetAllEvents()
	ctx.JSON(http.StatusOK, APIResponse.Success("All events retrieved successfully", gin.H{"events": allEvents}))
}

func (ctrl *eventController) GetEventByID(ctx *gin.Context) {
	eventID := ctx.Param("id")

	event := ctrl.service.GetEventByID(eventID)
	ctx.JSON(http.StatusOK, APIResponse.Success("Event retrieved successfully", gin.H{"event": event}))
}

func (ctrl *eventController) UpdateEvent(ctx *gin.Context) {
	eventID := ctx.Param("id")

	userIDRaw, exists := ctx.Get("userID")
	if !exists {
		exception := HTTPException.NewUnauthorizedException("Unauthorized", nil)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		exception := HTTPException.NewBadRequestException("User ID not found in context", nil)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	var request EventDTO.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := HTTPException.NewBadRequestException("Bad Request Exception", err.Error())
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	if err := ctrl.validator.Struct(request); err != nil {
		exception := utils.FormatValidationErrors(err)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	ctrl.service.UpdateEvent(userID, eventID, request)
	ctx.JSON(http.StatusOK, APIResponse.Success("Event updated successfully", nil))
}

func (ctrl *eventController) DeleteEvent(ctx *gin.Context) {
	eventID := ctx.Param("id")

	userIDRaw, exists := ctx.Get("userID")
	if !exists {
		exception := HTTPException.NewUnauthorizedException("Unauthorized", nil)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		exception := HTTPException.NewBadRequestException("User ID not found in context", nil)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	ctrl.service.DeleteEvent(userID, eventID)
	ctx.JSON(http.StatusOK, APIResponse.Success("Event deleted successfully", nil))
}
