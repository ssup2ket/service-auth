package model

import (
	modeluuid "github.com/ssup2ket/ssup2ket-auth-service/pkg/model/uuid"
)

type Outbox struct {
	ID            modeluuid.ModelUUID `gorm:"primaryKey;type:binary(16)"`
	AggregateType string              `gorm:"size:255"`
	AggregateID   string              `gorm:"size:255"`
	Type          string              `gorm:"size:255"`
	Payload       string              `gorm:"size:255"`
}
