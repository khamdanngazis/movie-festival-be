package entities

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleAdmin string = "admin"
	RoleUser  string = "user"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"type:varchar(50);default:'user'" json:"role"` // admin or user
	LoggedIn  bool           `gorm:"type:boolean;default:false" json:"logged_in"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
