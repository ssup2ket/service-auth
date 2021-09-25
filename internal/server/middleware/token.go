package middleware

import (
	"context"
	"fmt"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
)

type ctxKeyUserID int
type ctxKeyUserLoginID int
type ctxKeyUserRole int

const (
	CtxKeyUserID      ctxKeyUserID      = 0
	CtxKeyUserLoginID ctxKeyUserLoginID = 0
	CtxKeyUserRole    ctxKeyUserRole    = 0
)

func SetUserIDToCtx(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, CtxKeyUserID, userID)
}

func SetUserLoginIDToCtx(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, CtxKeyUserLoginID, userID)
}

func SetUserRoleToCtx(ctx context.Context, userRole model.UserRole) context.Context {
	return context.WithValue(ctx, CtxKeyUserRole, userRole)
}

func GetUserIDFromCtx(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(CtxKeyUserID).(string)
	if !ok {
		return "", fmt.Errorf("no user ID in context")
	}
	return userID, nil
}

func GetUserLoginIDFromCtx(ctx context.Context) (string, error) {
	userLoginID, ok := ctx.Value(CtxKeyUserLoginID).(string)
	if !ok {
		return "", fmt.Errorf("no user login ID in context")
	}
	return userLoginID, nil
}

func GetUserRoleFromCtx(ctx context.Context) (model.UserRole, error) {
	userRole, ok := ctx.Value(CtxKeyUserRole).(model.UserRole)
	if !ok {
		return "", fmt.Errorf("no user role in context")
	}
	return userRole, nil
}
