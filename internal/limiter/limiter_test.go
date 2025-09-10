package limiter

import (
    "context"
    "testing"
    "time"

    "github.com/youruser/rate-limiter/internal/backend"
)

func TestTokenBucketAllow(t *testing.T) {
    store := backend.NewMemStore()
    tb := NewTokenBucket(1.0, 5, store) // 1 token/sec, cap 5
    ctx := context.Background()
    key := "user:1"

    // consume 1 token - should be allowed
    allowed, err := tb.Allow(ctx, key, 1)
    if err != nil || !allowed {
        t.Fatal("expected allowed")
    }

    // consume 5 tokens - should fail (not enough)
    allowed, _ = tb.Allow(ctx, key, 5)
    if allowed {
        t.Fatal("expected not allowed")
    }

    // simulate time pass: set nowFunc to +10s for refill
    tb.nowFunc = func() time.Time { return time.Now().Add(10 * time.Second) }
    allowed, _ = tb.Allow(ctx, key, 3) // should be allowed after refill
    if !allowed {
        t.Fatal("expected allowed after refill")
    }
}

