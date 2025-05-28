package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"os"
	"time"
)

func Connect() *mongo.Client {
	uri := os.Getenv("MONGODB_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(uri).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
		SetReadConcern(readconcern.Majority())

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("mongo.Connect error: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo.Ping error: %v", err)
	}

	coll := client.Database("notifications").Collection("subscribers")
	idxModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := coll.Indexes().CreateOne(ctx, idxModel); err != nil {
		log.Printf("could not create index: %v", err)
	}

	return client
}
