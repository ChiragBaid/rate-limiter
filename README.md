# rate-limiter

A minimal pluggable rate-limiter service in Go demonstrating token-bucket algorithm with pluggable storage backends (in-memory, Redis). Inspired by various open-source limiters.

## Features
- Token-bucket algorithm
- In-memory store for local testing
- Redis-backed store for distributed deployments
- HTTP middleware (API-key or IP-based keys)
- Simple benchmarks and tests

## Quick start
```bash
# In-memory
go run ./cmd/server

# With redis
REDIS_ADDR=localhost:6379 go run ./cmd/server
