package limiter

import (
	"context"
	"sync"
	"time"
)

// Store interface persists token state
type Store interface {
	// Get returns tokens, lastRefillUnixNano, error
	Get(ctx context.Context, key string) (int64, int64, error)
	// Set sets tokens and lastRefillUnixNano
	Set(ctx context.Context, key string, tokens int64, lastRefill int64) error
}

// TokenBucket implements a token-bucket algorithm with a pluggable Store.
type TokenBucket struct {
	Rate     float64       // tokens per second
	Capacity int64         // max tokens
	store    Store
	mu       sync.Mutex
	nowFunc  func() time.Time
}

// NewTokenBucket creates a token bucket instance.
func NewTokenBucket(rate float64, capacity int64, store Store) *TokenBucket {
	return &TokenBucket{
		Rate:     rate,
		Capacity: capacity,
		store:    store,
		nowFunc:  time.Now,
	}
}

// Allow tries to consume tokensToConsume tokens for key. Returns true if allowed.
func (tb *TokenBucket) Allow(ctx context.Context, key string, tokensToConsume int64) (bool, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tokens, lastRefill, err := tb.store.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if lastRefill == 0 {
		// initialize full bucket on first use
		tokens = tb.Capacity
		lastRefill = tb.nowFunc().UnixNano()
	}

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

	// save current state (no change in tokens other than refill)
	if err := tb.store.Set(ctx, key, tokens, lastRefill); err != nil {
		return false, err
	}
	return false, nil
}
