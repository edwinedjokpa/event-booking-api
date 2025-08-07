package db

import (
	"fmt"
	"log"
	"reflect"

	"github.com/edwinedjokpa/event-booking-api/internal/app/event"
	"github.com/edwinedjokpa/event-booking-api/internal/app/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}
	return db, nil
}

func RunMigrations(db *gorm.DB) {
	modelsToMigrate := []any{
		&user.User{},
		&event.Event{},
	}

	fmt.Println("Checking all registered models:")
	for _, model := range modelsToMigrate {
		fmt.Println(" - ", reflect.TypeOf(model))
	}

	if err := db.AutoMigrate(&user.User{}, &event.Event{}); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	log.Println("Database migrations executed successfully!")
}
