package main

import (
	"github.com/awhdesmond/tiinyplanet/pkg/trips"
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

	db, err := utils.MakeMongoDatabase(cfg.MongoURL, cfg.MongoDBName)
	if err != nil {
		return nil, err
	}

	collabStore := trips.NewCollabStore(nc, rdb)
	tripStore := trips.NewStore(db)

	svc, err := trips.NewCollabService(collabStore, tripStore)
	if err != nil {
		return nil, err
	}

	collabSvr := trips.MakeCollabServer(svc, logger)

	baseServer := grpc.NewServer()
	trips.RegisterCollabServiceServer(baseServer, collabSvr)
	reflection.Register(baseServer)
	return baseServer, nil
}
