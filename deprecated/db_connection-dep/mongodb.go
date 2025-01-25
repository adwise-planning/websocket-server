package db_connection

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoCtx = context.Background()

func InitializeMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(MongoCtx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(MongoCtx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	MongoClient = client
	log.Println("Connected to MongoDB")
}
