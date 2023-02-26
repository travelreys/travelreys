package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tiinyplanet/tiinyplanet/pkg/api"
	"github.com/tiinyplanet/tiinyplanet/pkg/auth"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/flights"
	"github.com/tiinyplanet/tiinyplanet/pkg/images"
	"github.com/tiinyplanet/tiinyplanet/pkg/maps"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"github.com/tiinyplanet/tiinyplanet/pkg/tripssync"
	"go.uber.org/zap"
)

func MakeAPIServer(cfg ServerConfig, logger *zap.Logger) (*http.Server, error) {
	// Databases, external services and persistent storage
	db, err := common.MakeMongoDatabase(cfg.MongoURL, cfg.MongoDBName)
	if err != nil {
		logger.Error("cannot connect to mongo db", zap.Error(err))
		return nil, err
	}
	nc, err := common.MakeNATSConn(cfg.NatsURL)
	if err != nil {
		logger.Error("cannot connect to NATS", zap.Error(err))
		return nil, err
	}
	rdb, err := common.MakeRedisClient(cfg.RedisURL)
	if err != nil {
		logger.Error("cannot connect to redis", zap.Error(err))
		return nil, err
	}

	ctx := context.Background()

	// Auth
	gp, err := auth.NewDefaultGoogleProvider()
	if err != nil {
		return nil, err
	}
	authStore := auth.NewStore(ctx, db, logger)
	authSvc := auth.NewService(gp, authStore, logger)
	authSvc = auth.ServiceWithRBACMiddleware(authSvc, logger)

	// Flights
	flightsSvc := flights.NewService(flights.NewDefaultWebAPI(logger))
	flightsSvc = flights.ServiceWithRBACMiddleware(flightsSvc, logger)

	// Images
	imageSvc := images.NewService(images.NewDefaultWebAPI(logger))
	imageSvc = images.ServiceWithRBACMiddleware(imageSvc, logger)

	// Maps
	mapsSvc, err := maps.NewDefaulService(logger)
	if err != nil {
		logger.Error("unable to connect map service", zap.Error(err))
		return nil, err
	}

	// Trips
	tripStore := trips.NewStore(ctx, db, logger)
	tripSvc := trips.NewService(tripStore, authSvc, imageSvc)
	tripSvc = trips.ServiceWithRBACMiddleware(tripSvc, logger)

	r := mux.NewRouter()
	securityMW := api.NewSecureHeadersMiddleware(cfg.CORSOrigin)
	wrwMW := api.NewWrappedReponseWriterMiddleware()
	loggingMW := api.NewMuxLoggingMiddleware(logger)
	metricsMW := api.NewMetricsMiddleware()

	// TripSync
	store := tripssync.NewStore(rdb, logger)
	msgStore := tripssync.NewMessageStore(nc, rdb, logger)
	tobStore := tripssync.NewTOBMessageStore(nc, rdb, logger)

	svc := tripssync.NewService(store, msgStore, tobStore, tripStore)
	wsSvr := tripssync.NewWebsocketServer(svc, logger)

	r.Use(securityMW.Handler)
	r.Use(wrwMW.Handler)
	r.Use(loggingMW.Handler)
	r.Use(metricsMW.Handler)

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/healthz", api.HealthzHandler)
	r.HandleFunc("/ws", wsSvr.HandleFunc)

	r.PathPrefix("/api/v1/auth").Handler(auth.MakeHandler(authSvc))
	r.PathPrefix("/api/v1/flights").Handler(flights.MakeHandler(flightsSvc))
	r.PathPrefix("/api/v1/images").Handler(images.MakeHandler(imageSvc))
	r.PathPrefix("/api/v1/maps").Handler(maps.MakeHandler(mapsSvc))
	r.PathPrefix("/api/v1/trips").Handler(trips.MakeHandler(tripSvc))

	return &http.Server{
		Handler: r,
		Addr:    cfg.HTTPBindAddress(),
	}, nil
}
