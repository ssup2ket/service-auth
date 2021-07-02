package model

import (
	"time"

	"gorm.io/gorm"
)

type UserInfo struct {
	UUID      string `gorm:"primaryKey;size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ID    string `gorm:"unique;size:20"` // Unique key
	Phone string `gorm:"size:13"`
	Email string `gorm:"size:40"`
}
