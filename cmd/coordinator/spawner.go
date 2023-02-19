package main

import (
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"github.com/tiinyplanet/tiinyplanet/pkg/tripssync"
	"go.uber.org/zap"
)

func MakeCoordinatorSpanwer(cfg ServerConfig, logger *zap.Logger) (*tripssync.CoordinatorSpawner, error) {
	db, err := common.MakeMongoDatabase(cfg.MongoURL, cfg.MongoDBName)
	if err != nil {
		logger.Error("MakeCoordinatorSpanwer", zap.Error(err))
		return nil, err
	}

	nc, err := common.MakeNATSConn(cfg.NatsURL)
	if err != nil {
		logger.Error("MakeCoordinatorSpanwer", zap.Error(err))
		return nil, err
	}

	rdb, err := common.MakeRedisClient(cfg.RedisURL, cfg.RedisClusterMode)
	if err != nil {
		logger.Error("MakeCoordinatorSpanwer", zap.Error(err))
		return nil, err
	}

	tripStore := trips.NewStore(db, logger)
	sesnStore := tripssync.NewSessionStore(rdb)
	smStore := tripssync.NewSyncMessageStore(nc, rdb)
	tobStore := tripssync.NewTOBMessageStore(nc, rdb)

	return tripssync.NewCoordinatorSpawner(
		sesnStore, smStore, tobStore, tripStore, logger,
	), nil
}
