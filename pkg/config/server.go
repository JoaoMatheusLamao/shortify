package config

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/semaphore"
)

// SetupServer sets up a new gin engine with a semaphore and cors middleware
func SetupServer(rd *Config) (engine *gin.Engine) {
	engine = gin.Default()

	setupSemaphore(engine)
	setupCors(engine)
	setupRedisDB(engine, rd)

	return engine
}

func setupSemaphore(engine *gin.Engine) {
	sema := semaphore.NewWeighted(10)
	engine.Use(func(c *gin.Context) {
		if err := sema.Acquire(c.Request.Context(), 1); err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		defer sema.Release(1)
		c.Next()
	})
}

func setupCors(engine *gin.Engine) {
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}

func setupRedisDB(engine *gin.Engine, cfg *Config) {
	cfg.Redis.FlushAll(context.Background())
	engine.Use(func(c *gin.Context) {
		ip := c.ClientIP()
		val, err := cfg.Redis.Get(context.Background(), ip).Result()
		if err == redis.Nil {
			log.Println("First request from IP: ", ip)
			err = cfg.Redis.Set(context.Background(), ip, 1, 60*time.Second).Err()
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

		log.Println("New request from IP: ", ip, " with count: ", val)
		requestCount, _ := strconv.Atoi(val)
		if requestCount >= 20 {
			ttl, err := cfg.Redis.TTL(context.Background(), ip).Result()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			c.Writer.Header().Set("Retry-After", ttl.String())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		err = cfg.Redis.Incr(context.Background(), ip).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.Next()
	})
}
