package event

import (
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller EventController, jwtSecretKey []byte) {
	router.GET("/events", controller.GetAllEvents)
	router.GET("/events/:id", controller.GetEventByID)

	authRouter := router.Group("/events")
	authRouter.Use(middleware.AuthMiddleware(jwtSecretKey))
	{
		authRouter.POST("/", controller.CreateEvent)
		authRouter.PUT("/:id", controller.UpdateEvent)
		authRouter.DELETE("/:id", controller.DeleteEvent)
	}
}
