package entities

import (
	"time"
)

type Viewership struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	MovieID   uint      `gorm:"not null" json:"movie_id"`
	UserID    *uint     `gorm:"index;null" json:"user_id"`          // Nullable, optional tracking of user
	IPAddress string    `gorm:"size:45;not null" json:"ip_address"` // To track by IP
	ViewedAt  time.Time `gorm:"autoCreateTime" json:"viewed_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
