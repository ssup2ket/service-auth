package uuid

import (
	"database/sql/driver"

	uuid "github.com/satori/go.uuid"
)

// UUID type for model
type UUIDModel struct {
	uuid.UUID
}

func (uuid UUIDModel) String() string {
	return uuid.UUID.String()
}

func (uuid UUIDModel) Value() (driver.Value, error) {
	return uuid.UUID[:], nil
}

func NewV4() UUIDModel {
	return UUIDModel{UUID: uuid.NewV4()}
}

func FromStringOrNil(input string) UUIDModel {
	return UUIDModel{UUID: uuid.FromStringOrNil(input)}
}
