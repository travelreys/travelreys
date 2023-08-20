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
	"github.com/travelreys/travelreys/pkg/finance"
	"github.com/travelreys/travelreys/pkg/images"
	"github.com/travelreys/travelreys/pkg/invites"
	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/media"
	"github.com/travelreys/travelreys/pkg/ogp"
	"github.com/travelreys/travelreys/pkg/social"
	"github.com/travelreys/travelreys/pkg/storage"
	"github.com/travelreys/travelreys/pkg/trips"
	"go.uber.org/zap"
)

func MakeAPIServer(cfg ServerConfig, logger *zap.Logger) (*http.Server, error) {
	// Databases, external services and persistent storage
	db, err := common.MakeDefaultMongoDatabase()
	if err != nil {
		logger.Error("cannot connect to mongo db", zap.Error(err))
		return nil, err
	}
	nc, err := common.MakeDefaultNATSConn()
	if err != nil {
		logger.Error("cannot connect to NATS", zap.Error(err))
		return nil, err
	}
	rdb, err := common.MakeDefaultRedisClient()
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
	authSvcWithVal := auth.SvcWithValidationMw(authSvc, logger)
	authSvcForAPI := auth.SvcWithRBACMw(authSvc, logger)
	authSvcForAPI = auth.SvcWithValidationMw(authSvcForAPI, logger)

	// Images
	imageSvc := images.NewService(images.NewDefaultWebAPI(logger))
	imageSvcForAPI := images.SvcWithRBACMw(imageSvc, logger)

	// Maps
	mapsSvc, err := maps.NewDefaulService(logger)
	if err != nil {
		logger.Error("unable to connect map service", zap.Error(err))
		return nil, err
	}

	// Finance
	finStore := finance.NewStore(rdb, logger)
	finSvc := finance.NewService(finStore, logger)
	finSvcForAPI := finance.SvcWithRBACMw(finSvc, logger)

	// Ogp
	ogpSvc := ogp.NewService()
	ogpSvc = ogp.SvcWithRBACMw(ogpSvc, logger)
	ogpSvcForAPI := ogp.SvcWithValidation(ogpSvc, logger)

	// Media
	mediaStore := media.NewStore(ctx, db, logger)
	mediaCDNProvider, err := media.NewDefaultCDNProvider()
	if err != nil {
		logger.Error("unable to connect cdn provider", zap.Error(err))
		return nil, err
	}
	mediaSvc := media.NewService(mediaStore, mediaCDNProvider, storageSvc, logger)

	// Trips
	tripStore := trips.NewStore(ctx, db, logger)
	tripSvc := trips.NewService(
		tripStore,
		authSvcWithVal,
		imageSvc,
		mediaSvc,
		storageSvc,
		logger,
	)
	tripSvcWithVal := trips.SvcWithValidationMw(tripSvc, logger)
	tripSvcForAPI := trips.SvcWithRBACMw(tripSvc, logger)
	tripSvcForAPI = trips.SvcWithValidationMw(tripSvcForAPI, logger)
	tripSyncSvc := trips.NewSyncService(
		tripStore,
		trips.NewSessionStore(rdb, logger),
		trips.NewSyncMsgStore(nc, logger),
	)
	wsSvr := trips.NewWebsocketServer(tripSyncSvc, logger)

	// Trips Invite
	inviteStore := invites.NewStore(ctx, db, logger)
	inviteSvc := invites.NewService(
		authSvcWithVal,
		tripSyncSvc,
		mailSvc,
		inviteStore,
		logger,
	)
	inviteSvc = invites.SvcWithValidationMw(inviteSvc, logger)
	inviteSvc = invites.SvcWithRBACMw(inviteSvc, tripSvcWithVal, authSvcWithVal, logger)

	// Social
	socialStore := social.NewStore(ctx, db, logger)
	socialSvc := social.NewService(socialStore, authSvcWithVal, tripSvcWithVal, mailSvc, logger)
	socialSvc = social.SvcWithRBACMw(socialSvc, tripSvcWithVal, logger)

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

	r.PathPrefix("/api/v1/auth").Handler(auth.MakeHandler(authSvcForAPI))
	r.PathPrefix("/api/v1/images").Handler(images.MakeHandler(imageSvcForAPI))
	r.PathPrefix("/api/v1/finance").Handler(finance.MakeHandler(finSvcForAPI))
	r.PathPrefix("/api/v1/maps").Handler(maps.MakeHandler(mapsSvc))
	r.PathPrefix("/api/v1/ogp").Handler(ogp.MakeHandler(ogpSvcForAPI))
	r.PathPrefix("/api/v1/social").Handler(social.MakeHandler(socialSvc))
	r.PathPrefix("/api/v1/trips").Handler(trips.MakeHandler(tripSvcForAPI))
	r.PathPrefix("/api/v1/invites").Handler(invites.MakeHandler(inviteSvc))

	return &http.Server{
		Handler: r,
		Addr:    cfg.HTTPBindAddress(),
	}, nil
}
