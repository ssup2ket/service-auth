package model

import (
	"time"

	modeluuid "github.com/ssup2ket/ssup2ket-auth-service/pkg/model/uuid"
)

type Outbox struct {
	ID        modeluuid.ModelUUID `gorm:"primaryKey;type:binary(16)"`
	CreatedAt time.Time

	AggregateType string `gorm:"column:aggregatetype;size:255"`
	AggregateID   string `gorm:"column:aggregateid;size:255"`
	EventType     string `gorm:"column:eventtype;size:255"`
	Payload       string `gorm:"size:255"`
	SpanContext   string `gorm:"column:spancontext;size:255"`
}
