package persistence

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// InMemoryDB is a simple in-memory key-value store
type InMemoryDB struct {
	*cache.Cache
}

// NewInMemoryDB creates a new InMemoryDB
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Cache: cache.New(1*time.Minute, 2*time.Minute),
	}
}

// Handler Limit the number of requests from an IP (6 requests per minute, have to wait to make a new request)
func (db *InMemoryDB) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		val, ttl, found := db.Cache.GetWithExpiration(ip)
		if !found {
			db.Cache.Set(ip, 1, 60*time.Second)
			c.Next()
			return
		}

		requestCount := val.(int)
		if requestCount >= 20 {
			db.Cache.Increment(ip, 1)
			c.Writer.Header().Set("Retry-After", ttl.Sub(time.Now()).String())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		db.Cache.Increment(ip, 1)
		c.Next()
	}
}
