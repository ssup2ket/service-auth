package request

import (
	"fmt"
	"net/mail"
	"regexp"

	gouuid "github.com/satori/go.uuid"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/entity"
)

func ValidateUserUUID(uuid string) error {
	if _, err := gouuid.FromString(uuid); err != nil {
		return fmt.Errorf("wrong uuid format")
	}

	return nil
}

func ValidateUserCreate(id, passwd, role, phone, email string) error {
	// Login ID
	idMatched, err := regexp.MatchString("^[a-zA-Z0-9]{8,20}$", id)
	if err != nil {
		return fmt.Errorf("wrong id regex")
	}
	if !idMatched {
		return fmt.Errorf("wrong id format")
	}

	// Password
	passwdMatched, err := regexp.MatchString("^[a-zA-Z0-9]{8,20}$", passwd)
	if err != nil {
		return fmt.Errorf("wrong password regex")
	}
	if !passwdMatched {
		return fmt.Errorf("wrong password format")
	}

	// Role
	if !entity.IsValidUserRole(role) {
		return fmt.Errorf("wrong role")
	}

	// Phone
	phoneMatched, err := regexp.MatchString("^[0-9]{3}[-]+[0-9]{4}[-]+[0-9]{4}$", phone)
	if err != nil {
		return fmt.Errorf("wrong phone regex")
	}
	if !phoneMatched {
		return fmt.Errorf("wrong phone format")
	}

	// Email
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("wrong email format")
	}

	return nil
}

func ValidateUserUpdate(uuid, passwd, role, phone, email string) error {
	// UUID
	if uuid != "" {
		if _, err := gouuid.FromString(uuid); err != nil {
			return fmt.Errorf("wrong uuid format")
		}
	}

	// Password
	if passwd != "" {
		passwdMatched, err := regexp.MatchString("^[a-zA-Z0-9]{8,20}$", passwd)
		if err != nil {
			return fmt.Errorf("wrong password regex")
		}
		if !passwdMatched {
			return fmt.Errorf("wrong password format")
		}
	}

	// Role
	if role != "" {
		if !entity.IsValidUserRole(role) {
			return fmt.Errorf("wrong role")
		}
	}

	// Phone
	if phone != "" {
		phoneMatched, err := regexp.MatchString("^[0-9]{3}[-]+[0-9]{4}[-]+[0-9]{4}$", phone)
		if err != nil {
			return fmt.Errorf("wrong phone regex")
		}
		if !phoneMatched {
			return fmt.Errorf("wrong phone format")
		}
	}

	// Email
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			return fmt.Errorf("wrong email format")
		}
	}

	return nil
}
