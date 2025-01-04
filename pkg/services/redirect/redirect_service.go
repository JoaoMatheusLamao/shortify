package redirect

import (
	"context"
	"shortify/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// FindOriginalURLAndRedirect is a function that redirects the user to the original URL
func FindOriginalURLAndRedirect(rd *config.Redis) gin.HandlerFunc {

	return func(c *gin.Context) {

		shortURL := c.Param("shortURL")

		originalURL, err := rd.Get(context.Background(), shortURL).Result()
		if err == redis.Nil {
			c.JSON(404, gin.H{"error": "URL not found"})
			return
		} else if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error: " + err.Error()})
			return
		}

		c.Redirect(302, originalURL)
	}
}
