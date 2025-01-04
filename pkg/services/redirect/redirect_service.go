package redirect

import (
	"context"
	"errors"
	"shortify/pkg/config"

	"github.com/gin-gonic/gin"
)

// RedirectControler is a function that redirects the user to the original URL
func RedirectControler(rd *config.Redis) gin.HandlerFunc {

	return func(c *gin.Context) {

		shortURL := c.Param("shortURL")

		originalURL, err := redirect(shortURL, rd)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		c.Redirect(302, originalURL)
	}
}

// redirect is a function that redirects the user to the original URL
func redirect(id string, rd *config.Redis) (string, error) {
	url, err := rd.Get(context.Background(), id).Result()
	if err != nil {
		return "", errors.New("Error getting the short URL from the database: " + err.Error())
	}

	return url, nil
}
