package event

import (
	"time"
)

type Event struct {
	ID          string     `gorm:"primaryKey;not null" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `gorm:"not null" json:"description"`
	Location    string     `gorm:"not null" json:"location"`
	Date        time.Time  `gorm:"not null" json:"date"`
	UserID      string     `gorm:"not null" json:"user_id"`
	CreatedAt   time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"default:NULL" json:"deleted_at,omitempty"`
}
