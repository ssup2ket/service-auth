package uuid

import (
	"database/sql/driver"

	gouuid "github.com/satori/go.uuid"
)

// UUID type for model
type ModelUUID struct {
	gouuid.UUID
}

func (uuid ModelUUID) String() string {
	return uuid.UUID.String()
}

func (uuid ModelUUID) Value() (driver.Value, error) {
	return uuid.UUID[:], nil
}

func NewV4() ModelUUID {
	return ModelUUID{UUID: gouuid.NewV4()}
}

func FromStringOrNil(input string) ModelUUID {
	return ModelUUID{UUID: gouuid.FromStringOrNil(input)}
}
