package database

import (
	"context"
	"fmt"
	"time"

	"clean-arch/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoDB *mongo.Database
var MongoClient *mongo.Client

// ConnectMongo connects to MongoDB and sets package-level MongoClient and MongoDB
func ConnectMongo(ctx context.Context, env *config.Env) error {
	if env == nil {
		return fmt.Errorf("env is nil")
	}
	clientOpts := options.Client().ApplyURI(env.MongoURI)

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		return fmt.Errorf("mongo.NewClient: %w", err)
	}

	// connect with the parent ctx (so cancellation from caller works)
	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("mongo.Connect: %w", err)
	}

	// ping with short timeout
	ctxPing, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(ctxPing, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctxPing)
		return fmt.Errorf("mongo.Ping: %w", err)
	}

	MongoClient = client
	MongoDB = client.Database(env.MongoDB)
	return nil
}

// CloseMongo disconnects the global Mongo client
func CloseMongo(ctx context.Context) error {
	if MongoClient == nil {
		return nil
	}
	return MongoClient.Disconnect(ctx)
}
