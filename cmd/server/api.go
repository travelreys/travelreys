package main

import (
	"net/http"

	"github.com/awhdesmond/tiinyplanet/pkg/images"
	"github.com/awhdesmond/tiinyplanet/pkg/trips"
	"github.com/awhdesmond/tiinyplanet/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func MakeAPIServer(cfg ServerConfig, logger *zap.Logger) (*http.Server, error) {

	db, err := utils.MakeMongoDatabase(cfg.MongoURL, cfg.MongoDBName)
	if err != nil {
		return nil, err
	}

	// Images
	unsplash := images.NewWebImageAPI()
	imageSvc := images.NewService(unsplash)

	// Trips
	tripStore := trips.NewStore(db)
	tripSvc := trips.NewService(tripStore)
	tripSvc = trips.ServiceWithLoggingMiddleware(tripSvc, logger)

	r := mux.NewRouter()
	securityMW := utils.NewSecureHeadersMiddleware(cfg.CORSOrigin)
	wrwMW := utils.NewWrappedReponseWriterMiddleware()
	loggingMW := utils.NewMuxLoggingMiddleware(logger)
	metricsMW := utils.NewMetricsMiddleware()

	r.Use(securityMW.Handler)
	r.Use(wrwMW.Handler)
	r.Use(loggingMW.Handler)
	r.Use(metricsMW.Handler)

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/healthz", utils.HealthzHandler)

	r.PathPrefix("/api/v1/trips").Handler(trips.MakeHandler(tripSvc))
	r.PathPrefix("/api/v1/images").Handler(images.MakeHandler(imageSvc))

	return &http.Server{
		Handler: r,
		Addr:    cfg.HTTPBindAddress(),
	}, nil
}
