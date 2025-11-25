package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

// ConnectMongo reads MONGO_URI & MONGO_DBNAME from env and connect
func ConnectMongo() error {
	uri := os.Getenv("MONGO_URI")
	dbname := os.Getenv("MONGO_DBNAME")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	if dbname == "" {
		dbname = "appdb"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}

	// ping
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	MongoDB = client.Database(dbname)
	log.Printf("connected to mongo: %s (db=%s)\n", uri, dbname)
	return nil
}
