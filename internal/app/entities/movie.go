package entities

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Duration    int            `gorm:"not null" json:"duration"` // in minutes
	Artists     string         `gorm:"type:varchar(255)" json:"artists"`
	Genres      string         `gorm:"type:varchar(255)" json:"genres"`
	WatchURL    string         `gorm:"type:varchar(255);not null" json:"watch_url"`
	Views       int            `gorm:"default:0" json:"views"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
