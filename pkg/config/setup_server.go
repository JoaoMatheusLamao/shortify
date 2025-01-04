package config

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
)

// SetupServer sets up a new gin engine with a semaphore and cors middleware
func SetupServer(rd *Redis) (engine *gin.Engine) {
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

func setupRedisDB(engine *gin.Engine, rd *Redis) {
	rd.FlushAll(context.Background())
	engine.Use(rd.Handler())
}
