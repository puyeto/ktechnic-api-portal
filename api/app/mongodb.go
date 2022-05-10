package app

import (
	"context"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB ...
var MongoDB *mongo.Database

// InitializeMongoDB Initialize MongoDB Connection
func InitializeMongoDB(dbURL, dbName string, logger *logrus.Logger) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect(ctx)

	logger.Printf("Mongo DB initialized %v", dbName)
	return client.Database(dbName)
}

// CreateIndexMongo create a mongodn index
func CreateIndexMongo(colName string) (string, error) {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"datetimestamp": -1, // index in ascending order
		}, Options: nil,
	}
	collection := MongoDB.Collection(colName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return collection.Indexes().CreateOne(ctx, mod)
}
