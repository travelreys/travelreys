package main

import (
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"github.com/tiinyplanet/tiinyplanet/pkg/tripssync"
	"github.com/tiinyplanet/tiinyplanet/pkg/utils"
	"go.uber.org/zap"
)

func MakeCoordinatorSpanwer(cfg ServerConfig, logger *zap.Logger) (*tripssync.CoordinatorSpawner, error) {

	db, err := utils.MakeMongoDatabase(cfg.MongoURL, cfg.MongoDBName)
	if err != nil {
		return nil, err
	}

	nc, err := utils.MakeNATSConn(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	rdb, err := utils.MakeRedisClient(cfg.RedisURL, cfg.RedisClusterMode)
	if err != nil {
		return nil, err
	}

	tripStore := trips.NewStore(db)
	sesnStore := tripssync.NewSessionStore(rdb)
	smStore := tripssync.NewSyncMessageStore(nc, rdb)
	tobStore := tripssync.NewTOBMessageStore(nc, rdb)

	return tripssync.NewCoordinatorSpawner(
		sesnStore, smStore, tobStore, tripStore, logger,
	), nil
}
