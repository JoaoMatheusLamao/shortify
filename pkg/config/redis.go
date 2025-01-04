package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// NewRedis is a function that returns a new Redis client
func NewRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()

	return rdb, err
}
