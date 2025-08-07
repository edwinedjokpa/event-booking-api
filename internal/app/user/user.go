package user

import "time"

type User struct {
	ID        string     `gorm:"primaryKey;not null" json:"id"`
	FirstName string     `gorm:"not null" json:"first_name"`
	LastName  string     `gorm:"not null" json:"last_name"`
	Email     string     `gorm:"unique;not null" json:"email"`
	Password  string     `gorm:"not null"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:NULL" json:"deleted_at,omitempty"`
}
