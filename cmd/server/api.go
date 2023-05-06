package main

import (
	"context"
	"crypto/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/travelreys/travelreys/pkg/api"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/email"
	"github.com/travelreys/travelreys/pkg/footprints"
	"github.com/travelreys/travelreys/pkg/images"
	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/media"
	"github.com/travelreys/travelreys/pkg/ogp"
	"github.com/travelreys/travelreys/pkg/storage"
	"github.com/travelreys/travelreys/pkg/trips"
	"github.com/travelreys/travelreys/pkg/tripssync"
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

	// Mail
	mailSvc := email.NewDefaultService()

	// Storage
	storageSvc, err := storage.NewDefaultStorageService(ctx)
	if err != nil {
		logger.Error("unable to connect storage service", zap.Error(err))
		return nil, err
	}

	// Auth
	authStore := auth.NewStore(ctx, db, rdb, logger)
	gp, err := auth.NewDefaultGoogleProvider()
	if err != nil {
		return nil, err
	}
	fb := auth.NewFacebookProvider()
	otp := auth.NewDefaultOTPProvider(authStore, rand.Reader)
	authSvc := auth.NewService(
		gp, fb, otp,
		authStore, cfg.SecureCookie,
		mailSvc, storageSvc, logger,
	)
	authSvcWithRBAC := auth.ServiceWithRBACMiddleware(authSvc, logger)

	// Images
	imageSvc := images.NewService(images.NewDefaultWebAPI(logger))
	imageSvc = images.ServiceWithRBACMiddleware(imageSvc, logger)

	// Maps
	mapsSvc, err := maps.NewDefaulService(logger)
	if err != nil {
		logger.Error("unable to connect map service", zap.Error(err))
		return nil, err
	}

	// Ogp
	ogpSvc := ogp.NewService()
	ogpSvc = ogp.ServiceWithRBACMiddleware(ogpSvc, logger)

	// Media
	mediaStore := media.NewStore(ctx, db, logger)
	mediaCDNProvider, err := media.NewDefaultCDNProvider()
	if err != nil {
		logger.Error("unable to connect cdn provider", zap.Error(err))
		return nil, err
	}
	mediaSvc := media.NewService(mediaStore, mediaCDNProvider, storageSvc, logger)
	mediaSvcWithRBAC := media.ServiceWithRBACMiddleware(mediaSvc, logger)

	// Trips
	tripStore := trips.NewStore(ctx, db, logger)
	tripSvc := trips.NewService(tripStore, authSvc, imageSvc, mediaSvc, storageSvc)
	tripSvc = trips.ServiceWithRBACMiddleware(tripSvc, logger)

	// TripSync
	store := tripssync.NewStore(rdb, logger)
	msgStore := tripssync.NewMessageStore(nc, rdb, logger)
	tobStore := tripssync.NewTOBMessageStore(nc, rdb, logger)

	svc := tripssync.NewService(store, msgStore, tobStore, tripStore)
	wsSvr := tripssync.NewWebsocketServer(svc, logger)

	// Footprints
	fpStore := footprints.NewStore(ctx, db, logger)
	fpSvc := footprints.NewService(fpStore, logger)
	fpSvc = footprints.ServiceWithRBACMiddleware(fpSvc, logger)

	r := mux.NewRouter()
	securityMW := api.NewSecureHeadersMiddleware(cfg.CORSOrigin)
	wrwMW := api.NewWrappedReponseWriterMiddleware()
	loggingMW := api.NewMuxLoggingMiddleware(logger)
	metricsMW := api.NewMetricsMiddleware()

	r.Use(securityMW.Handler)
	r.Use(wrwMW.Handler)
	r.Use(loggingMW.Handler)
	r.Use(metricsMW.Handler)

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/healthz", api.HealthzHandler)
	r.HandleFunc("/ws", wsSvr.HandleFunc)

	r.PathPrefix("/api/v1/auth").Handler(auth.MakeHandler(authSvcWithRBAC))
	r.PathPrefix("/api/v1/images").Handler(images.MakeHandler(imageSvc))
	r.PathPrefix("/api/v1/maps").Handler(maps.MakeHandler(mapsSvc))
	r.PathPrefix("/api/v1/media").Handler(media.MakeHandler(mediaSvcWithRBAC))
	r.PathPrefix("/api/v1/ogp").Handler(ogp.MakeHandler(ogpSvc))
	r.PathPrefix("/api/v1/trips").Handler(trips.MakeHandler(tripSvc))
	r.PathPrefix("/api/v1/footprints").Handler(footprints.MakeHandler(fpSvc))

	return &http.Server{
		Handler: r,
		Addr:    cfg.HTTPBindAddress(),
	}, nil
}
