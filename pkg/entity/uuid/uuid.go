package uuid

import (
	"database/sql/driver"

	gouuid "github.com/satori/go.uuid"
)

// UUID type for entity
type EntityUUID struct {
	gouuid.UUID
}

func (uuid EntityUUID) String() string {
	return uuid.UUID.String()
}

func (uuid EntityUUID) Value() (driver.Value, error) {
	return uuid.UUID[:], nil
}

func NewV4() EntityUUID {
	return EntityUUID{UUID: gouuid.NewV4()}
}

func FromStringOrNil(input string) EntityUUID {
	return EntityUUID{UUID: gouuid.FromStringOrNil(input)}
}
