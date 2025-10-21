package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	// MongoDB connection string
	mongoURI := "mongodb://localhost:27017"
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	// Verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB unreachable: ", err)
	}

	fmt.Println("MongoDB connected ✅")
	return client
}

// GetDatabase returns the database instance
func GetDatabase(client *mongo.Client) *mongo.Database {
	return client.Database("alumni_db")
}
