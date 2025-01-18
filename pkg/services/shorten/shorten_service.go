package shorten

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"shortify/pkg/config"
	"shortify/pkg/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		go insertShortURLInMongo(shortURL, request.OriginalURL, c.ClientIP(), cfg)
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

// insertShortURLInMongo is a function that inserts a short URL in the MongoDB database
func insertShortURLInMongo(shortURL string, originalURL string, ip string, cfg *config.Config) {

	collection := cfg.Mongo.Database("shortify").Collection("short_urls")

	filter := bson.D{{Key: "original_url", Value: originalURL}}
	safeIP := strings.ReplaceAll(ip, ".", "_")
	update := bson.D{
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "short_url", Value: shortURL},
			{Key: "first_creation", Value: time.Now()},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "last_creation", Value: time.Now()},
		}},
		{Key: "$inc", Value: bson.D{
			{Key: "counter_creation", Value: 1},
			{Key: fmt.Sprintf("ip_counter_creation.%s", safeIP), Value: 1},
		}},
	}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Println("Error inserting or updating the short URL in the MongoDB database: " + err.Error())
	}
}
