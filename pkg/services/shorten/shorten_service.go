package shorten

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"shortify/pkg/config"
	"shortify/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// ShortenControler is a controller that returns a 200 status code
func ShortenControler(rd *config.Redis) gin.HandlerFunc {

	var input struct {
		OriginalURL string `json:"original_url" binding:"required"`
	}

	var response struct {
		ShortURL string `json:"short_url"`
	}

	return func(c *gin.Context) {

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		protocolAndHost := utils.GetCurrentProtocolAndHost(c)

		shortURL, err := shorten(input.OriginalURL, rd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response.ShortURL = fmt.Sprintf("%s/r/%s", protocolAndHost, shortURL)

		c.JSON(http.StatusOK, response)
	}

}

// shorten is a service that returns a short URL
func shorten(originalURL string, rd *config.Redis) (string, error) {

	shortURL, err := generateShortURL(originalURL, rd)
	if err != nil {
		return "", err
	}

	return shortURL, nil

}

func generateShortURL(originalURL string, rd *config.Redis) (string, error) {

	id := utils.GenerateUniqueID(originalURL)
	_, err := rd.Get(context.Background(), id).Result()
	if err == redis.Nil {
		err = rd.Set(context.Background(), id, originalURL, 5*time.Hour).Err()
		if err != nil {
			return "", errors.New("Error saving the short URL to the database: " + err.Error())
		}
		return id, nil
	} else if err != nil {
		return "", errors.New("Error getting the short URL from the database: " + err.Error())
	}

	return id, nil
}
