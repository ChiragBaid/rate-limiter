package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/ChiragBaid/rate-limiter/internal/backend"
)

func TestTokenBucketAllow(t *testing.T) {
	store := backend.NewMemStore()
	tb := NewTokenBucket(1.0, 5, store) // 1 token/sec, cap 5
	ctx := context.Background()
	key := "user:1"

	allowed, err := tb.Allow(ctx, key, 1)
	if err != nil || !allowed {
		t.Fatal("expected allowed")
	}

	allowed, _ = tb.Allow(ctx, key, 5)
	if allowed {
		t.Fatal("expected not allowed")
	}

	// simulate refill by advancing nowFunc
	tb.nowFunc = func() time.Time { return time.Now().Add(10 * time.Second) }
	allowed, _ = tb.Allow(ctx, key, 3)
	if !allowed {
		t.Fatal("expected allowed after refill")
	}
}
