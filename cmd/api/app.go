package main

import (
	"log"
	"net/http"

	"github.com/edwinedjokpa/event-booking-api/internal/app/auth"
	"github.com/edwinedjokpa/event-booking-api/internal/app/event"
	"github.com/edwinedjokpa/event-booking-api/internal/app/user"
	"github.com/edwinedjokpa/event-booking-api/internal/config"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/db"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/middleware"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/redis"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/service/otp"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/service/session"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/validator"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupApp initializes all application components and returns a configured Gin router.
func SetupApp(config *config.Config) (*gin.Engine, error) {
	// Connect to Database
	gormDB, err := db.NewGormDB(config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Run Database Migrations
	db.RunMigrations(gormDB)

	// Initialize a single, configured validator instance.
	appValidator := validator.NewValidator()

	// Initialize the Redis Client
	redisClient, err := redis.NewRedisClient(config.RedisAddr)
	if err != nil {
		return nil, err
	}

	// Initialize the Session Service
	sessionService := session.NewSessionService(redisClient)

	// Initialize the Otp Service
	otpService := otp.NewOTPService(redisClient)

	// Initialize Repositories
	userRepository := user.NewUserRepository(gormDB)
	eventRepository := event.NewEventRepository(gormDB)

	// Initialize Services
	authService := auth.NewAuthService(userRepository, config.JWTSecret, sessionService, otpService)
	eventService := event.NewEventService(eventRepository)

	// Initialize Controllers
	authController := auth.NewAuthController(authService, appValidator)
	eventController := event.NewEventController(eventService, appValidator)

	// Create a new Gin router with custom middleware stack.
	router := gin.New()

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	// Register global middleware
	router.Use(gin.Logger())
	router.Use(cors.Default())
	router.Use(middleware.RecoveryMiddleware())

	api := router.Group("/api")

	api.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to the Event Booking API")
	})

	// Register Router
	auth.RegisterRoutes(api, authController)
	event.RegisterRoutes(api, eventController, []byte(config.JWTSecret))

	return router, nil
}
