package model

import (
	"time"

	"gorm.io/gorm"

	modeluuid "github.com/ssup2ket/ssup2ket-auth-service/pkg/model/uuid"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

func IsValidUserRole(role string) bool {
	if role == string(UserRoleAdmin) {
		return true
	} else if role == string(UserRoleUser) {
		return true
	}
	return false
}

type UserInfo struct {
	ID        modeluuid.ModelUUID `gorm:"primaryKey;type:binary(16)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	LoginID string   `gorm:"unique;size:20"` // Unique key
	Role    UserRole `gorm:"size:20"`
	Phone   string   `gorm:"size:13"`
	Email   string   `gorm:"size:40"`
}
