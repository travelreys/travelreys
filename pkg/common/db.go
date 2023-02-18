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

const (
	NATSClientName = "tiinyplanet"
)

func MakeNATSConn(url string) (*nats.Conn, error) {
	return nats.Connect(url, nats.Name(NATSClientName))
}

// Redis

func MakeRedisClient(uri string, isClusterMode bool) (redis.UniversalClient, error) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	if isClusterMode {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{opt.Addr},
			Password: opt.Password,
		}), nil
	}

	rdb := redis.NewClient(opt)
	return rdb, err
}

// Mongo

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
