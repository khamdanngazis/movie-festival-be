package entities

import (
	"time"
)

type Vote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_movie" json:"user_id"`
	MovieID   uint      `gorm:"not null;uniqueIndex:idx_user_movie" json:"movie_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
