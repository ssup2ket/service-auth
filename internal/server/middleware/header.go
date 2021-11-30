package middleware

import (
	"context"
	"fmt"
)

type ctxKeyRequestID int

const (
	CtxKeyRequestID ctxKeyRequestID = 0

	HeaderRequestID = "X-Request-ID"
)

func SetRequestIDToCtx(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, CtxKeyRequestID, requestID)
}

func GetRequestIDFromCtx(ctx context.Context) (string, error) {
	requestID, ok := ctx.Value(CtxKeyRequestID).(string)
	if !ok {
		return "", fmt.Errorf("no request ID in context")
	}
	return requestID, nil
}
