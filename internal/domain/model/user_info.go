package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/pkg/uuidmodel"
)

type UserInfo struct {
	ID        uuidmodel.UUIDModel `gorm:"primaryKey;type:binary(16)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	LoginID string `gorm:"unique;size:20"` // Unique key
	Phone   string `gorm:"size:13"`
	Email   string `gorm:"size:40"`
}
