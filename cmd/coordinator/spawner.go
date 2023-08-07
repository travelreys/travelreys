package main

import (
	"context"

	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/media"
	"github.com/travelreys/travelreys/pkg/storage"
	"github.com/travelreys/travelreys/pkg/trips"
	"go.uber.org/zap"
)

func MakeCoordinatorSpanwer(logger *zap.Logger) (*trips.Spawner, error) {
	db, err := common.MakeDefaultMongoDatabase()
	if err != nil {
		logger.Error("cannot connect to mongo", zap.Error(err))
		return nil, err
	}

	etcd, err := common.MakeDefaultEtcdClient()
	if err != nil {
		logger.Error("cannot connect to etcd", zap.Error(err))
		return nil, err
	}

	nc, err := common.MakeDefaultNATSConn()
	if err != nil {
		logger.Error("cannot connect to nats", zap.Error(err))
		return nil, err
	}

	ctx := context.Background()

	// Storage
	storageSvc, err := storage.NewDefaultStorageService(ctx)
	if err != nil {
		logger.Error("unable to connect storage service", zap.Error(err))
		return nil, err
	}

	// Maps
	mapsSvc, err := maps.NewDefaulService(logger)
	if err != nil {
		logger.Error("unable to connect map service", zap.Error(err))
		return nil, err
	}

	// Media
	mediaStore := media.NewStore(ctx, db, logger)
	mediaCDNProvider, err := media.NewDefaultCDNProvider()
	if err != nil {
		logger.Error("unable to connect cdn provider", zap.Error(err))
		return nil, err
	}

	return trips.NewSpawner(
		mapsSvc,
		media.NewService(mediaStore, mediaCDNProvider, storageSvc, logger),
		trips.NewStore(ctx, db, logger),
		trips.NewSessionStore(etcd, logger),
		trips.NewSyncMsgBroadcastStore(nc, logger),
		trips.NewSyncMsgTOBStore(nc, logger),
		logger,
	), nil
}
