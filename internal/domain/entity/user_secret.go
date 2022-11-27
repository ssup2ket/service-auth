package entity

import (
	"time"

	"gorm.io/gorm"

	"github.com/ssup2ket/service-auth/pkg/entity/uuid"
)

type UserSecret struct {
	ID        uuid.EntityUUID `gorm:"primaryKey;type:binary(16)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	PasswdHash       []byte `gorm:"size:4096"`
	PasswdSalt       []byte `gorm:"size:20"`
	RefreshTokenHash []byte `gorm:"size:4096"`
	RefreshTokenSalt []byte `gorm:"size:20"`
}
