package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoInstance stores the client and database instance
type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var MI MongoInstance

// ConnectDB establishes a connection to MongoDB
func ConnectDB() {
	// Get MongoDB URI from environment variable
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not set")
	}

	// Define client options
	clientOptions := options.Client().ApplyURI(uri)

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping the database to check if the connection is successful
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Could not ping MongoDB:", err)
	}

	fmt.Println("✅ Database connected!")

	// Store the MongoDB instance globally
	MI = MongoInstance{
		Client: client,
		DB:     client.Database(os.Getenv("DATABASE_NAME")),
	}
}
