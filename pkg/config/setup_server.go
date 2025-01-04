package config

import (
	"net/http"
	"shortify/pkg/persistence"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
)

// SetupServer sets up a new gin engine with a semaphore and cors middleware
func SetupServer() (engine *gin.Engine, db *persistence.InMemoryDB) {
	engine = gin.Default()

	setupSemaphore(engine)
	setupCors(engine)

	db = setupInMemoryDB(engine)

	return engine, db
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

func setupInMemoryDB(engine *gin.Engine) (db *persistence.InMemoryDB) {
	db = persistence.NewInMemoryDB()
	db.Flush()
	engine.Use(db.Handler())
	return db
}
