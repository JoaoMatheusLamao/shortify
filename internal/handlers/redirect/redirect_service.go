package redirect

import (
	"context"
	"shortify/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// FindOriginalURLAndRedirect is a function that redirects the user to the original URL
func FindOriginalURLAndRedirect(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		shortURL := c.Param("shortURL")

		originalURL, err := getOriginalURLFromRedis(shortURL, cfg)
		if err != nil {
			if err == redis.Nil {
				c.JSON(404, gin.H{"error": "URL not found"})
			} else {
				c.JSON(500, gin.H{"error": "Internal server error: " + err.Error()})
			}
			return
		}

		go cfg.Mongo.UpsertAnalyticsInMongo(shortURL, originalURL, c.ClientIP())

		c.Redirect(302, originalURL)
	}
}

// getOriginalURLFromRedis is a function that retrieves the original URL from Redis
func getOriginalURLFromRedis(shortURL string, cfg *config.Config) (string, error) {
	originalURL, err := cfg.Redis.Get(context.Background(), shortURL).Result()
	if err != nil {
		return "", err
	}
	return originalURL, nil
}
