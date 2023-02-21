package common

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// NATS

const (
	NATSClientName = "tiinyplanet"
)

func MakeNATSConn(url string) (*nats.Conn, error) {
	return nats.Connect(url, nats.Name(NATSClientName))
}

// Redis
func redisOptsToUnivOpts(opts *redis.Options) *redis.UniversalOptions {
	return &redis.UniversalOptions{
		Addrs:    []string{opts.Addr},
		Password: opts.Password,
	}
}

func MakeRedisClient(uri string) (redis.UniversalClient, error) {
	opts, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewUniversalClient(redisOptsToUnivOpts(opts))
	return rdb, nil
}

// Mongo

var (
	DbReqTimeout = 3 * time.Second
)

func MakeMongoDatabase(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbReqTimeout)
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
