package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/pkg/uuid"
)

type UserInfo struct {
	UUID      uuid.UUIDModel `gorm:"primaryKey;type:binary(16)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ID    string `gorm:"unique;size:20"` // Unique key
	Phone string `gorm:"size:13"`
	Email string `gorm:"size:40"`
}
