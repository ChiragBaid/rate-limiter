package limiter

import (
    "context"
    "sync"
    "time"
)

// Store defines how token state is persisted (in-memory or redis).
type Store interface {
    // Get returns (tokens, lastRefillTimestampUnixNano, err)
    Get(ctx context.Context, key string) (int64, int64, error)
    // Set updates the token count and last refill timestamp
    Set(ctx context.Context, key string, tokens int64, lastRefill int64) error
}

// TokenBucket config
type TokenBucket struct {
    Rate       float64       // tokens per second
    Capacity   int64         // max tokens
    store      Store
    mu         sync.Mutex
    nowFunc    func() time.Time // for testing
}

// NewTokenBucket creates a new instance
func NewTokenBucket(rate float64, capacity int64, store Store) *TokenBucket {
    return &TokenBucket{
        Rate:     rate,
        Capacity: capacity,
        store:    store,
        nowFunc:  time.Now,
    }
}

// Allow checks if a token can be consumed for key
func (tb *TokenBucket) Allow(ctx context.Context, key string, tokensToConsume int64) (bool, error) {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    tokens, lastRefill, err := tb.store.Get(ctx, key)
    if err != nil {
        return false, err
    }
    if lastRefill == 0 {
        // first time: initialize full bucket
        tokens = tb.Capacity
        lastRefill = tb.nowFunc().UnixNano()
    }

    // refill calculation
    now := tb.nowFunc().UnixNano()
    elapsed := float64(now-lastRefill) / float64(time.Second.Nanoseconds())
    refill := int64(elapsed * tb.Rate)
    if refill > 0 {
        tokens += refill
        if tokens > tb.Capacity {
            tokens = tb.Capacity
        }
        lastRefill = now
    }

    if tokens >= tokensToConsume {
        tokens -= tokensToConsume
        if err := tb.store.Set(ctx, key, tokens, lastRefill); err != nil {
            return false, err
        }
        return true, nil
    }

    // save updated tokens and lastRefill even if not allowed
    if err := tb.store.Set(ctx, key, tokens, lastRefill); err != nil {
        return false, err
    }
    return false, nil
}

