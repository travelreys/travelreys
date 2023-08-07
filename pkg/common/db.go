package common

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	NATSClientName = "travelreys"
)

var (
	DbReqTimeout = 3 * time.Second

	mongoURL    = os.Getenv("TRAVELREYS_MONGO_URL")
	mongoDBName = os.Getenv("TRAVELREYS_MONGO_DBNAME")
	natsURL     = os.Getenv("TRAVELREYS_NATS_URL")
	etcdURL     = os.Getenv("TRAVELREYS_ETCD_URL")
	redisURL    = os.Getenv("TRAVELREYS_REDIS_URL")
)

// NATS

func MakeDefaultNATSConn() (*nats.Conn, error) {
	return MakeNATSConn(natsURL)
}

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

func MakeDefaultRedisClient() (redis.UniversalClient, error) {
	return MakeRedisClient(redisURL)
}

func MakeRedisClient(uri string) (redis.UniversalClient, error) {
	opts, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewUniversalClient(redisOptsToUnivOpts(opts))
	return rdb, nil
}

// Etcd

func MakeDefaultEtcdClient() (*clientv3.Client, error) {
	return MakeEtcdClient(etcdURL)
}

func MakeEtcdClient(uri string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{uri},
		DialTimeout: DbReqTimeout,
	})
}

// Mongo

func MakeDefaultMongoDatabase() (*mongo.Database, error) {
	return MakeMongoDatabase(mongoURL, mongoDBName)
}

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

func MongoIsDupError(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}
