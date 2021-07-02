package service

import (
	"fmt"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo"
)

// Config
var config *ServiceConfigs

type ServiceConfigs struct{}

// Init
func Init(c *ServiceConfigs) {
	config = c
}

// Error
var (
	// Common
	ErrServerErr error = fmt.Errorf("server error")

	// Repository
	ErrRepoNotFound    error = fmt.Errorf("repo resource not found")
	ErrRepoConflict    error = fmt.Errorf("repo conflict")
	ErrRepoServerError error = fmt.Errorf("repo server error")
)

func getReturnErr(err error) error {
	switch err {
	case repo.ErrNotFound:
		return ErrRepoNotFound
	case repo.ErrConflict:
		return ErrRepoConflict
	case ErrRepoServerError:
		return ErrRepoServerError
	}
	return ErrServerErr
}
