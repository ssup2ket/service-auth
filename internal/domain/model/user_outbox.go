package model

import (
	modeluuid "github.com/ssup2ket/ssup2ket-auth-service/pkg/model/uuid"
)

type UserOutbox struct {
	ID            modeluuid.ModelUUID `gorm:"primaryKey;type:binary(16)"`
	AggregateType string              `gorm:"size:255"`
	AggregateID   string              `gorm:"size:255"`
	Type          string              `gorm:"size:255"`
	Payload       string              `gorm:"size:255"`
}

type UserOutboxPayload struct {
	ID      string `json:"id"`
	LoginID string `json:"loginId"`
	Role    string `json:"role"`
}
