package middlewares

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const userCtxKey = contextKey("userClaims")

func WithClaims(ctx context.Context, claims jwt.MapClaims) context.Context {
	return context.WithValue(ctx, userCtxKey, claims)
}

func FromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(userCtxKey).(jwt.MapClaims)
	if !ok {
		return &Claims{}, false
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return &Claims{}, false
	}

	return &Claims{
		UserID: int64(userID),
	}, ok
}

type Claims struct {
	UserID int64
}
