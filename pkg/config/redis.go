package config

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

// newClientRedis is a function that returns a new Redis client
func (cfg *Config) newClientRedis() error {

	addr := "localhost:6379"

	enviroment := os.Getenv("ENVIROMENT_EXEC")
	if enviroment == "prod" {
		addr = "redis:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	cfg.Redis = &RedisInternal{
		Redis: rdb,
		mu:    sync.Mutex{},
	}

	_, err := rdb.Ping(context.Background()).Result()

	return err
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
