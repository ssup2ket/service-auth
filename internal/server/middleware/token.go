package middleware

import (
	"context"
	"fmt"
)

type ctxKeyUserID int
type ctxKeyUserLoginID int

const (
	CtxKeyUserID      ctxKeyUserID      = 0
	CtxKeyUserLoginID ctxKeyUserLoginID = 0
)

func SetUserIDToCtx(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, CtxKeyUserID, userID)
}

func SetUserLoginIDToCtx(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, CtxKeyUserLoginID, userID)
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
