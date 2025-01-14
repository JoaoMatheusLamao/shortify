package config

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Config - a struct that holds a redis client
type Config struct {
	Redis *RedisInternal
	Mongo *mongo.Client
}

// NewConfig - a function that returns a new Config struct
func NewConfig() (*Config, error) {

	cfg := new(Config)

	err := cfg.newClientRedis()
	if err != nil {
		return cfg, err
	}

	err = cfg.newClientMongo()
	if err != nil {
		return cfg, err
	}

	return cfg, err
}

// CloseAll - a function that closes all connections
func (cfg *Config) CloseAll() {
	if cfg.Redis != nil {
		cfg.Redis.Redis.Close()
	}
}
