package middlewares

import (
	"context"
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/service/limiter"
)

const shareLimiterCtxKey = contextKey("limiterCtxKey")

func ShareLimiterMiddleware(rl *limiter.RateLimiter) func(http.Handler) http.Handler {
	return LimiterMiddleware(rl, shareLimiterCtxKey)
}

func ShareLimiterFromContext(ctx context.Context) (*limiter.RateLimiter, bool) {
	rl, ok := ctx.Value(shareLimiterCtxKey).(*limiter.RateLimiter)
	if !ok {
		return nil, false
	}
	return rl, ok
}

func LimiterMiddleware(rl *limiter.RateLimiter, limiterCtxKey contextKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := context.WithValue(r.Context(), limiterCtxKey, rl)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
