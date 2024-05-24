package adapters

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient represents a MongoDB client.
type MongoDBClient struct {
	Client *mongo.Client
}

// NewMongoDBClient creates a new MongoDB client.
func NewMongoDBClient(ctx context.Context, uri string) (*MongoDBClient, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return &MongoDBClient{Client: client}, nil
}
