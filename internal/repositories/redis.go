package repositories

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisInternal is a struct that contains a Redis client and a mutex
type RedisInternal struct {
	Redis *redis.Client
	mu    sync.Mutex
}

// NewRedisInternal is a function that returns a new RedisInternal struct
func NewRedisInternal() (*RedisInternal, error) {
	addr := getRedisAddress()

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &RedisInternal{
		Redis: rdb,
		mu:    sync.Mutex{},
	}, nil
}

func getRedisAddress() string {
	if os.Getenv("ENVIROMENT_EXEC") == "prod" {
		return "redis:6379"
	}
	return "localhost:6379"
}

// Get is a function that returns the value of a key
func (r *RedisInternal) Get(ctx context.Context, key string) *redis.StringCmd {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Redis.Get(ctx, key)
}

// Set is a function that sets a key value pair
func (r *RedisInternal) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Redis.Set(ctx, key, value, expiration)
}

// Expire is a function that sets a key expiration time
func (r *RedisInternal) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Redis.Expire(ctx, key, expiration)
}

// FlushAll is a function that flushes all keys
func (r *RedisInternal) FlushAll(ctx context.Context) *redis.StatusCmd {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Redis.FlushAll(ctx)
}

// TTL is a function that returns the time to live of a key
func (r *RedisInternal) TTL(ctx context.Context, key string) *redis.DurationCmd {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Redis.TTL(ctx, key)
}

// Incr is a function that increments a key
func (r *RedisInternal) Incr(ctx context.Context, key string) *redis.IntCmd {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Redis.Incr(ctx, key)
}
