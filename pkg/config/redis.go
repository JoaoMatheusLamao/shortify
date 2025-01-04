package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

// newRedis is a function that returns a new Redis client
func (cfg *Config) newRedis() error {

	addr := "localhost:6379"

	enviroment := os.Getenv("ENVIROMENT_EXEC")
	if enviroment == "prod" {
		addr = "redis:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	cfg.Redis = rdb

	_, err := rdb.Ping(context.Background()).Result()

	return err
}
