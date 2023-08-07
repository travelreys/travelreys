package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/travelreys/travelreys/pkg/api"
	"go.uber.org/zap"
)

const (
	cfgFlagLogLevel = "log-level"
)

func main() {
	viper.SetDefault(cfgFlagLogLevel, "info")

	pflag.String(cfgFlagLogLevel, "", "log level")
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	// Logger

	logger, _ := api.InitZap(viper.GetString(cfgFlagLogLevel))
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	// Make Coordinator Spawner
	spawner, err := MakeCoordinatorSpanwer(logger)
	if err != nil {
		logger.Panic("error initialising api server", zap.Error(err))
	}

	// Run Coordinator Spawner
	go func() {
		hostname, _ := os.Hostname()
		logger.Info("starting spawner", zap.String("hostname", hostname))
		if err := spawner.Run(); err != nil {
			logger.Panic("error starting spawner", zap.Error(err))
			os.Exit(1)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
}
