package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tiinyplanet/tiinyplanet/pkg/images"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"github.com/tiinyplanet/tiinyplanet/pkg/tripssync"
	"github.com/tiinyplanet/tiinyplanet/pkg/utils"
	"go.uber.org/zap"
)

func MakeAPIServer(cfg ServerConfig, logger *zap.Logger) (*http.Server, error) {

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

	// Images
	unsplash := images.NewWebImageAPI()
	imageSvc := images.NewService(unsplash)

	// Trips
	tripStore := trips.NewStore(db)
	tripSvc := trips.NewService(tripStore, imageSvc)
	tripSvc = trips.ServiceWithLoggingMiddleware(tripSvc, logger)

	r := mux.NewRouter()
	securityMW := utils.NewSecureHeadersMiddleware(cfg.CORSOrigin)
	wrwMW := utils.NewWrappedReponseWriterMiddleware()
	loggingMW := utils.NewMuxLoggingMiddleware(logger)
	metricsMW := utils.NewMetricsMiddleware()

	// Collab
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
	r.HandleFunc("/healthz", utils.HealthzHandler)
	r.HandleFunc("/ws", proxyServer.HandleFunc)
	r.PathPrefix("/api/v1/trips").Handler(trips.MakeHandler(tripSvc))
	r.PathPrefix("/api/v1/images").Handler(images.MakeHandler(imageSvc))

	return &http.Server{
		Handler: r,
		Addr:    cfg.HTTPBindAddress(),
	}, nil
}
