package shorten

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"shortify/internal/config"
	"shortify/internal/utils"

	"time"

	"github.com/gin-gonic/gin"
)

// CreateShortenURL is a controller that creates a short URL
func CreateShortenURL(cfg *config.Config) gin.HandlerFunc {

	var request struct {
		OriginalURL   string `json:"original_url" binding:"required"`
		ExpirationMin int    `json:"expiration_min"`
	}

	var response struct {
		ShortURL string `json:"short_url"`
	}

	return func(c *gin.Context) {

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		expiration := 5 * time.Minute
		if request.ExpirationMin > 0 {
			expiration = time.Duration(request.ExpirationMin) * time.Minute
		}

		shortURL, err := generateShortURL(request.OriginalURL, expiration, cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		protocolAndHost := utils.GetCurrentProtocolAndHost(c)
		response.ShortURL = fmt.Sprintf("%s/r/%s", protocolAndHost, shortURL)

		c.JSON(http.StatusOK, response)

		go cfg.Mongo.UpsertAnalyticsInMongo(shortURL, request.OriginalURL, c.ClientIP())
	}
}

func generateShortURL(originalURL string, expiration time.Duration, cfg *config.Config) (string, error) {

	hash := utils.GenerateUniqueID(originalURL)

	_, err := cfg.Redis.Get(context.Background(), hash).Result()
	if err != nil {
		err = cfg.Redis.Set(context.Background(), hash, originalURL, expiration).Err()
		if err != nil {
			return "", errors.New("Error saving the short URL to the database: " + err.Error())
		}
		return hash, nil
	}

	cfg.Redis.Expire(context.Background(), hash, expiration).Err()

	return hash, nil
}
