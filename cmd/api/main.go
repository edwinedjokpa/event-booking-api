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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment configurations
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Connect Database
	gormDB, err := db.NewGormDB(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Run Migrations
	db.RunMigrations(gormDB)

	// Initialize Redis
	redis.InitRedis(config.RedisAddr)

	// Use gin.New() to build a custom middleware stack
	router := gin.New()

	router.SetTrustedProxies(nil)

	// Register global middleware here, before any routes
	router.Use(gin.Logger())
	router.Use(cors.Default())
	router.Use(middleware.RecoveryMiddleware())

	// Initialize Repositories
	userRepository := user.NewUserRepository(gormDB)
	eventRepository := event.NewEventRepository(gormDB)

	// Initialize Services
	authService := auth.NewAuthService(userRepository, config.JWTSecret)
	eventService := event.NewEventService(eventRepository)

	// Initialize Controllers
	authController := auth.NewAuthController(authService)
	eventController := event.NewEventController(eventService)

	api := router.Group("/api")

	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Event Booking API"})
	})

	auth.RegisterRoutes(api, authController)
	event.RegisterRoutes(api, eventController, []byte(config.JWTSecret))

	if err := router.Run(":8000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
