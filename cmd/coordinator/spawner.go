package main

import (
	"context"

	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/trips"
	"github.com/travelreys/travelreys/pkg/tripssync"
	"go.uber.org/zap"
)

func MakeCoordinatorSpanwer(cfg ServerConfig, logger *zap.Logger) (*tripssync.Spawner, error) {
	db, err := common.MakeMongoDatabase(cfg.MongoURL, cfg.MongoDBName)
	if err != nil {
		logger.Error("cannot connect to mongo", zap.Error(err))
		return nil, err
	}

	nc, err := common.MakeNATSConn(cfg.NatsURL)
	if err != nil {
		logger.Error("cannot connect to nats", zap.Error(err))
		return nil, err
	}

	rdb, err := common.MakeRedisClient(cfg.RedisURL)
	if err != nil {
		logger.Error("cannot connect to redi", zap.Error(err))
		return nil, err
	}
	// Maps
	mapsSvc, err := maps.NewDefaulService(logger)
	if err != nil {
		logger.Error("unable to connect map service", zap.Error(err))
		return nil, err
	}

	tripStore := trips.NewStore(context.Background(), db, logger)
	store := tripssync.NewStore(rdb, logger)
	msgStore := tripssync.NewMessageStore(nc, rdb, logger)
	tobStore := tripssync.NewTOBMessageStore(nc, rdb, logger)

	return tripssync.NewSpawner(mapsSvc, store, msgStore, tobStore, tripStore, logger), nil
}
