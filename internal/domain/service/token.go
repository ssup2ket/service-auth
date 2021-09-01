package service

import (
	"context"
)

type TokenService interface {
	CreateToken(ctx context.Context, id, passwd string) (string, error)
}

type TokenServiceImp struct {
}

func NewTokenServiceImp() *TokenServiceImp {
	return &TokenServiceImp{}
}

func (t *TokenServiceImp) CreateToken(ctx context.Context, id, password string) (string, error) {
	return "", nil
}
