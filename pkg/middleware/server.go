package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"shortify/internal/config"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/unrolled/secure"
	"golang.org/x/sync/semaphore"
)

// SetupServer sets up a new gin engine with a semaphore and cors middleware
func SetupServer(rd *config.Config) (engine *gin.Engine) {

	gin.SetMode(gin.ReleaseMode)
	engine = gin.New()

	setupSemaphore(engine)
	setupCors(engine)
	setupRedisDB(engine, rd)
	setupSSL(engine)

	return engine
}

func setupSemaphore(engine *gin.Engine) {
	_, max := getMaxThrottlingRules()
	sema := semaphore.NewWeighted(max)
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

func setupRedisDB(engine *gin.Engine, cfg *config.Config) {
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
		max, _ := getMaxThrottlingRules()
		if requestCount >= int(max) {
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

// getMaxRequestCount returns the maximum number of requests allowed
func getMaxThrottlingRules() (countByIP int64, countGlobal int64) {
	defaultCountByIP := int64(5)
	defaultCountGlobal := int64(10)

	countByIP = getEnvAsInt64("MAX_REQUEST_COUNT_BY_IP", defaultCountByIP)
	countGlobal = getEnvAsInt64("MAX_REQUEST_COUNT_GLOBAL", defaultCountGlobal)

	return countByIP, countGlobal
}

// setupSSL is a function that sets up the SSL configuration for the server
func setupSSL(engine *gin.Engine) {
	engine.Use(func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":8080",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			log.Println("Error traying make a secure https: " + err.Error())
			return
		}
		c.Next()
	})
}

func getEnvAsInt64(name string, defaultValue int64) int64 {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return int64(value)
}
