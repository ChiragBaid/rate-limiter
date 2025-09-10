package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ChiragBaid/rate-limiter/internal/limiter"
)

type Config struct {
	TokenBucket      *limiter.TokenBucket
	KeyFunc          func(r *http.Request) string
	TokensPerRequest int64
}

func NewRateLimitMiddleware(cfg Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := cfg.KeyFunc(r)
			allowed, err := cfg.TokenBucket.Allow(context.Background(), key, cfg.TokensPerRequest)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			if !allowed {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func KeyByIP(r *http.Request) string {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	return "ip:" + ip
}

func KeyByHeader(r *http.Request) string {
	v := r.Header.Get("X-API-Key")
	if v == "" {
		return KeyByIP(r)
	}
	return "api:" + v
}
