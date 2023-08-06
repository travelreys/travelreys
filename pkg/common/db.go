package common

import (
	"context"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// NATS

const (
	NATSClientName = "travelreys"
)

var (
	DbReqTimeout = 3 * time.Second
)

func MakeNATSConn(url string) (*nats.Conn, error) {
	return nats.Connect(url, nats.Name(NATSClientName))
}

// Etcd

func MakeEtcdClient(uri string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{uri},
		DialTimeout: DbReqTimeout,
	})
}

// Mongo

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
