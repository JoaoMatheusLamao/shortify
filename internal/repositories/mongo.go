package repositories

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

// MongoInternal is a struct that contains a MongoDB client
type MongoInternal struct {
	client *mongo.Client
}

// NewMongoInternal is a function that returns a new MongoInternal struct
func NewMongoInternal() (*MongoInternal, error) {

	uri := os.Getenv("MONGO_URI")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, err
	}

	return &MongoInternal{
		client: client,
	}, nil

}

// UpsertAnalyticsInMongo is a function that inserts or updates the analytics data in MongoDB
func (m *MongoInternal) UpsertAnalyticsInMongo(shortURL, originalURL, clientIP string) {

	collection := m.client.Database("shortify").Collection("short_urls")
	safeIP := strings.ReplaceAll(clientIP, ".", "_")

	filter := bson.D{{Key: "original_url", Value: originalURL}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "last_find", Value: time.Now()},
			{Key: "short_url", Value: shortURL},
			{Key: "last_creation", Value: time.Now()},
		}},
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "first_creation", Value: time.Now()},
		}},
		{Key: "$inc", Value: bson.D{
			{Key: "counter_find", Value: 1},
			{Key: fmt.Sprintf("ip_counter_find.%s", safeIP), Value: 1},
			{Key: "counter_creation", Value: 1},
			{Key: fmt.Sprintf("ip_counter_creation.%s", safeIP), Value: 1},
		}},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Println("Error updating analytics in the MongoDB database: " + err.Error())
	}
}
