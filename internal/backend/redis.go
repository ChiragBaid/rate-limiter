package backend

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStore{rdb: rdb}
}

func (s *RedisStore) Get(ctx context.Context, key string) (int64, int64, error) {
	res, err := s.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	if len(res) == 0 {
		return 0, 0, nil
	}
	tokens, _ := strconv.ParseInt(res["tokens"], 10, 64)
	last, _ := strconv.ParseInt(res["lastRefill"], 10, 64)
	return tokens, last, nil
}

func (s *RedisStore) Set(ctx context.Context, key string, tokens int64, lastRefill int64) error {
	vals := map[string]interface{}{
		"tokens":     strconv.FormatInt(tokens, 10),
		"lastRefill": strconv.FormatInt(lastRefill, 10),
	}
	return s.rdb.HSet(ctx, key, vals).Err()
}
