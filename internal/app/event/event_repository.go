package event

import (
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event Event) error
	FindAll() ([]Event, error)
	FindOneByID(eventID string) (*Event, error)
	Update(event Event) error
	Delete(eventID string) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db}
}

func (repo *eventRepository) Create(event Event) error {
	if err := repo.db.Create(&event).Error; err != nil {
		return err
	}
	return nil
}

func (repo *eventRepository) FindAll() ([]Event, error) {
	var events []Event
	if err := repo.db.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (repo *eventRepository) FindOneByID(eventID string) (*Event, error) {
	var event Event
	if err := repo.db.First(&event, "id = ?", eventID).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (repo *eventRepository) Update(event Event) error {
	if err := repo.db.Save(&event).Error; err != nil {
		return err
	}
	return nil
}

func (repo *eventRepository) Delete(eventID string) error {
	if err := repo.db.Delete(&Event{}, "id = ?", eventID).Error; err != nil {
		return err
	}
	return nil
}
