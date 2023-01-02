package main

import (
	"github.com/awhdesmond/tiinyplanet/pkg/trips"
	"github.com/awhdesmond/tiinyplanet/pkg/tripssync"
	"github.com/awhdesmond/tiinyplanet/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func MakeCollabServer(cfg ServerConfig, logger *zap.Logger) (*grpc.Server, error) {
	nc, err := utils.MakeNATSConn(cfg.NatsURL)
	if err != nil {
		return nil, err
	}
	rdb, err := utils.MakeRedisClient(cfg.RedisURL, cfg.RedisClusterMode)
	if err != nil {
		return nil, err
	}

	db, err := utils.MakeMongoDatabase(cfg.MongoURL, cfg.MongoDBName)
	if err != nil {
		return nil, err
	}

	tripStore := trips.NewStore(db)
	collabStore := tripssync.NewStore(tripStore, nc, rdb)

	pxy, err := tripssync.NewProxy(collabStore)
	if err != nil {
		return nil, err
	}

	collabSvr := tripssync.MakeServer(pxy, logger)

	baseServer := grpc.NewServer()
	tripssync.RegisterCollabServiceServer(baseServer, collabSvr)
	reflection.Register(baseServer)
	return baseServer, nil
}
