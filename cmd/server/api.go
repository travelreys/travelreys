package main

import (
	"net/http"
	"os"

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
		logger.Error("MakeAPIServer", zap.Error(err))
		return nil, err
	}
	nc, err := common.MakeNATSConn(cfg.NatsURL)
	if err != nil {
		logger.Error("MakeAPIServer", zap.Error(err))
		return nil, err
	}
	rdb, err := common.MakeRedisClient(cfg.RedisURL, cfg.RedisClusterMode)
	if err != nil {
		logger.Error("MakeAPIServer", zap.Error(err))
		return nil, err
	}

	// Auth
	gp, err := auth.NewGoogleProvider(
		os.Getenv("TIINYPLANET_OAUTH_GOOGLE_SECRET_FILE"),
	)
	if err != nil {
		return nil, err
	}
	authStore := auth.NewStore(db, logger)
	authSvc := auth.NewService(gp, authStore, logger)
	authSvc = auth.ServiceWithRBACMiddleware(authSvc, logger)

	// Flights
	skyscanner := flights.NewSkyscannerAPI(
		os.Getenv("TIINYPLANET_SKYSCANNER_APIKEY"),
		os.Getenv("TIINYPLANET_SKYSCANNER_APIHOST"),
	)
	flightsSvc := flights.NewService(skyscanner)
	flightsSvc = flights.ServiceWithRBACMiddleware(flightsSvc, logger)

	// Images
	unsplash := images.NewWebImageAPI(
		os.Getenv("TIINYPLANET_UNSPLASH_ACCESSKEY"),
		logger,
	)
	imageSvc := images.NewService(unsplash)
	imageSvc = images.ServiceWithRBACMiddleware(imageSvc, logger)

	// Maps
	mapsSvc, err := maps.NewService(
		os.Getenv("TIINYPLANET_GOOGLE_MAPS_APIKEY"),
		logger,
	)
	if err != nil {
		logger.Error("MakeAPIServer", zap.Error(err))
		return nil, err
	}

	// Trips
	tripStore := trips.NewStore(db, logger)
	tripSvc := trips.NewService(tripStore, authSvc, imageSvc)
	tripSvc = trips.ServiceWithRBACMiddleware(tripSvc, logger)

	r := mux.NewRouter()
	securityMW := api.NewSecureHeadersMiddleware(cfg.CORSOrigin)
	wrwMW := api.NewWrappedReponseWriterMiddleware()
	loggingMW := api.NewMuxLoggingMiddleware(logger)
	metricsMW := api.NewMetricsMiddleware()

	// TripSync
	sesnStore := tripssync.NewSessionStore(rdb)
	smStore := tripssync.NewSyncMessageStore(nc, rdb)
	tobStore := tripssync.NewTOBMessageStore(nc, rdb)

	pxy, err := tripssync.NewProxy(sesnStore, smStore, tobStore, tripStore)
	if err != nil {
		return nil, err
	}
	proxyServer := tripssync.MakeProxyServer(pxy, logger)

	r.Use(securityMW.Handler)
	r.Use(wrwMW.Handler)
	r.Use(loggingMW.Handler)
	r.Use(metricsMW.Handler)

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/healthz", api.HealthzHandler)
	r.HandleFunc("/ws", proxyServer.HandleFunc)

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
