package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var (
	MongoConnectTimeout = 30 * time.Second
)

func MakeMongoDatabase(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoConnectTimeout)
	defer cancel()

	opts := options.Client().
		ApplyURI(uri).
		SetMinPoolSize(32).
		SetMaxConnecting(32).
		SetReadPreference(readpref.PrimaryPreferred()).
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority(), writeconcern.J(true)))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
}
