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

// CreateShortenURL is a controller that creates a short URL
func CreateShortenURL(rd *config.Redis) gin.HandlerFunc {

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

		shortURL, err := generateShortURL(input.OriginalURL, rd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		protocolAndHost := utils.GetCurrentProtocolAndHost(c)
		response.ShortURL = fmt.Sprintf("%s/r/%s", protocolAndHost, shortURL)

		c.JSON(http.StatusOK, response)
	}
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
