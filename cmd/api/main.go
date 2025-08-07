package main

import (
	"log"

	"github.com/edwinedjokpa/event-booking-api/internal/config"
)

func main() {
	// Load environment configurations
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Call the setup function to get the configured router
	router, err := SetupApp(config)
	if err != nil {
		log.Fatalf("Error setting up application: %v", err)
	}

	// Start the server
	if err := router.Run(":" + config.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
