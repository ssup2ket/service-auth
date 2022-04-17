package domain

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/config"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/service"
)

type Domain struct {
	// Config
	Configs *config.Configs

	// Service
	User  service.UserService
	Token service.TokenService
}

func New(c *config.Configs) (*Domain, error) {
	// Init domain and config
	domain := Domain{
		Configs: c,
	}

	// Init repo
	txMySQL, primaryMySQL, secondaryMySQL, err := repo.New(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to init repo pkg")
		return nil, fmt.Errorf("failed to init repo pkg")
	}
	outboxRepoPrimaryMysql := repo.NewOutboxRepoImp(primaryMySQL)
	userInfoRepoPrimaryMysql := repo.NewUserInfoRepoImp(primaryMySQL)
	userInfoRepoSecondaryMysql := repo.NewUserInfoRepoImp(secondaryMySQL)
	userSecretRepoPrimaryMysql := repo.NewUserSecretRepoImp(primaryMySQL)
	userSecretRepoSecondaryMysql := repo.NewUserSecretRepoImp(secondaryMySQL)

	// Init services
	userService := service.NewUserServiceImp(txMySQL, outboxRepoPrimaryMysql,
		userInfoRepoPrimaryMysql, userInfoRepoSecondaryMysql, userSecretRepoPrimaryMysql, userSecretRepoSecondaryMysql)
	tokenService := service.NewTokenServiceImp(txMySQL, userInfoRepoSecondaryMysql, userSecretRepoPrimaryMysql, userSecretRepoSecondaryMysql)

	domain.User = userService
	domain.Token = tokenService

	return &domain, nil
}
