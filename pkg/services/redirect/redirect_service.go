package redirect

import (
	"context"
	"fmt"
	"log"
	"shortify/pkg/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOriginalURLAndRedirect is a function that redirects the user to the original URL
func FindOriginalURLAndRedirect(cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {

		shortURL := c.Param("shortURL")

		originalURL, err := cfg.Redis.Get(context.Background(), shortURL).Result()
		if err == redis.Nil {
			c.JSON(404, gin.H{"error": "URL not found"})
			return
		} else if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error: " + err.Error()})
			return
		}

		c.Redirect(302, originalURL)

		go insertAnalyticsInMongo(originalURL, c.ClientIP(), cfg)
	}
}

// insertAnalyticsInMongo is a function that inserts the analytics data in MongoDB
func insertAnalyticsInMongo(originalURL, clientIP string, cfg *config.Config) {

	collection := cfg.Mongo.Database("shortify").Collection("short_urls")
	safeIP := strings.ReplaceAll(clientIP, ".", "_")

	filter := bson.D{{Key: "original_url", Value: originalURL}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "last_find", Value: time.Now()},
		}},
		{Key: "$inc", Value: bson.D{
			{Key: "counter_find", Value: 1},
			{Key: fmt.Sprintf("ip_counter_find.%s", safeIP), Value: 1},
		}},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Println("Error updating analytics in the MongoDB database: " + err.Error())
	}
}
