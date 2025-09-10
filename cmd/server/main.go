package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/youruser/rate-limiter/internal/backend"
    "github.com/youruser/rate-limiter/internal/limiter"
    "github.com/youruser/rate-limiter/internal/middleware"
)

func main() {
    // Read config from env
    redisAddr := os.Getenv("REDIS_ADDR") // optional
    var store limiter.Store
    if redisAddr != "" {
        rs := backend.NewRedisStore(redisAddr)
        store = rs
        fmt.Println("Using Redis store:", redisAddr)
    } else {
        ms := backend.NewMemStore()
        store = ms
        fmt.Println("Using In-Memory store")
    }

    // Example: 10 tokens/sec, burst capacity 20
    tb := limiter.NewTokenBucket(10.0, 20, store)

    mw := middleware.NewRateLimitMiddleware(middleware.Config{
        TokenBucket: tb,
        KeyFunc:     middleware.KeyByHeader,
        TokensPerRequest: 1,
    })

    mux := http.NewServeMux()
    mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "hello world")
    })

    handler := mw(mux)
    srv := &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }

    log.Println("listening on :8080")
    if err := srv.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

