package main

import (
	"context"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"github.com/tiinyplanet/tiinyplanet/pkg/tripssync"
	"go.uber.org/zap"
)

func MakeCoordinatorSpanwer(cfg ServerConfig, logger *zap.Logger) (*tripssync.CoordinatorSpawner, error) {
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

	tripStore := trips.NewStore(context.Background(), db, logger)
	store := tripssync.NewStore(rdb)
	msgStore := tripssync.NewMessageStore(nc, rdb)
	tobStore := tripssync.NewTOBMessageStore(nc, rdb)

	return tripssync.NewCoordinatorSpawner(store, msgStore, tobStore, tripStore, logger), nil
}
