package model

import (
	"time"

	"gorm.io/gorm"
)

type UserSecret struct {
	UUID      string `gorm:"primaryKey;size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	PasswdHash []byte `gorm:"size:4096"`
	PasswdSalt []byte `gorm:"size:20"`
}
