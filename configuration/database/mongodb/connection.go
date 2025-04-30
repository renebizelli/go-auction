package mongodb

import (
	"context"
	"fmt"
	"os"
	"renebizelli/go-leilao/configuration/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_URL = "mongodb://localhost:27017"
	MONGODB_DB  = "auction"
)

func NewMongoDBConnection(ctx context.Context) (*mongo.Database, error) {

	mongoDbBURL := os.Getenv("MONGODB_URL")
	mongoDbDatabase := os.Getenv("MONGODB_DB")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDbBURL))

	if err != nil {
		logger.Error("Error connecting to MongoDB: %v", err)
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		logger.Error("Error pinging MongoDB: %v", err)
		return nil, err
	}

	logger.Info(fmt.Sprintf("Connected to MongoDB at %s", mongoDbBURL))

	return client.Database(mongoDbDatabase), nil
}
