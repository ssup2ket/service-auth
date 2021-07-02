package request

import (
	"fmt"
	"net/mail"
	"regexp"

	gouuid "github.com/satori/go.uuid"
)

func ValidateUserUUID(uuid string) error {
	if _, err := gouuid.FromString(uuid); err != nil {
		return fmt.Errorf("wrong uuid format")
	}

	return nil
}

func ValidateUserCreate(id, passwd, phone, email string) error {
	// ID
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

func ValidateUserUpdate(uuid, passwd, phone, email string) error {
	// UUID
	if _, err := gouuid.FromString(uuid); err != nil {
		return fmt.Errorf("wrong uuid format")
	}

	// Password
	passwdMatched, err := regexp.MatchString("^[a-zA-Z0-9]{8,20}$", passwd)
	if err != nil {
		return fmt.Errorf("wrong password regex")
	}
	if !passwdMatched {
		return fmt.Errorf("wrong password format")
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
