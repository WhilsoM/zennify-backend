package http

import "context"

type ctxKey string

const (
	ctxKeyRequestID ctxKey = "request_id"
	ctxKeyUserID    ctxKey = "user_id"
	ctxKeyUsername  ctxKey = "username"
)

func withRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, requestID)
}

func requestIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyRequestID).(string)
	return v
}

func withUserClaims(ctx context.Context, userID, username string) context.Context {
	ctx = context.WithValue(ctx, ctxKeyUserID, userID)
	ctx = context.WithValue(ctx, ctxKeyUsername, username)
	return ctx
}

func userClaimsFromContext(ctx context.Context) (string, string) {
	userID, _ := ctx.Value(ctxKeyUserID).(string)
	username, _ := ctx.Value(ctxKeyUsername).(string)
	return userID, username
}
