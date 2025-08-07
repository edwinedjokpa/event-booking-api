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
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/services"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment configurations
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Connect to Database
	gormDB, err := db.NewGormDB(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Run Database Migrations
	db.RunMigrations(gormDB)

	// Initialize a single, configured validator instance.
	validator := validator.NewValidator()

	// Initialize the Redis Client
	redisClient, err := redis.NewRedisClient("localhost:6379")
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}

	// Initialize the Session Service
	sessionService := services.NewSessionService(redisClient)

	// Initialize the Otp Service
	otpService := services.NewOTPService(redisClient)

	// Initialize Repositories
	userRepository := user.NewUserRepository(gormDB)
	eventRepository := event.NewEventRepository(gormDB)

	// Initialize Services
	authService := auth.NewAuthService(userRepository, config.JWTSecret, sessionService, otpService)
	eventService := event.NewEventService(eventRepository)

	// Initialize Controllers
	authController := auth.NewAuthController(authService, validator)
	eventController := event.NewEventController(eventService)

	// Use gin.New() to build a custom middleware stack
	router := gin.New()
	router.SetTrustedProxies(nil)

	// Register global middleware here, before any routes
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

	if err := router.Run(":" + config.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
