package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/pkg/uuidmodel"
)

type UserSecret struct {
	ID        uuidmodel.UUIDModel `gorm:"primaryKey;type:binary(16)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	PasswdHash []byte `gorm:"size:4096"`
	PasswdSalt []byte `gorm:"size:20"`
}
