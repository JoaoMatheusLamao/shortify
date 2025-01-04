package config

import (
	"github.com/redis/go-redis/v9"
)

// Config - a struct that holds a redis client
type Config struct {
	Redis *redis.Client
}

// NewConfig - a function that returns a new Config struct
func NewConfig() (*Config, error) {

	cfg := new(Config)

	err := cfg.newRedis()

	return cfg, err
}

// CloseAll - a function that closes all connections
func (cfg *Config) CloseAll() {
	if cfg.Redis != nil {
		cfg.Redis.Close()
	}
}
