package config

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Redis - a struct that holds a redis client
type Redis struct {
	*redis.Client
}

// NewRedis is a function that returns a new Redis client
func NewRedis() (*Redis, error) {

	addr := "localhost:6379"

	enviroment := os.Getenv("ENVIROMENT_EXEC")
	if enviroment == "prod" {
		addr = "redis:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := rdb.Ping(context.Background()).Result()

	return &Redis{rdb}, err
}

// Handler Limit the number of requests from an IP (6 requests per minute, have to wait to make a new request)
func (db *Redis) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		val, err := db.Get(context.Background(), ip).Result()
		if err == redis.Nil {
			err = db.Set(context.Background(), ip, 1, 60*time.Second).Err()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error()})
				return
			}
			c.Next()
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error ", "message": err.Error()})
			return
		}

		requestCount, _ := strconv.Atoi(val)
		if requestCount >= 20 {
			ttl, err := db.TTL(context.Background(), ip).Result()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			c.Writer.Header().Set("Retry-After", ttl.String())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		err = db.Incr(context.Background(), ip).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.Next()
	}
}
