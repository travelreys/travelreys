package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tiinyplanet/tiinyplanet/pkg/api"
	"go.uber.org/zap"
)

const (
	cfgFlagHost     = "host"
	cfgFlagHostname = "hostname"
	cfgFlagLogLevel = "log-level"

	cfgFlagMongoURL         = "mongo-url"
	cfgFlagMongoDBName      = "mongo-dbname"
	cfgFlagNatsURL          = "nats-url"
	cfgFlagRedisURL         = "redis-url"
	cfgFlagRedisClusterMode = "redis-cluster-mode"

	envVarPrefix = "TIINYPLANET"
)

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`

	CORSOrigin string `mapstructure:"cors-origin"`

	NatsURL          string `mapstructure:"nats-url"`
	RedisURL         string `mapstructure:"redis-url"`
	RedisClusterMode bool   `mapstructure:"redis-cluster-mode"`
	MongoURL         string `mapstructure:"mongo-url"`
	MongoDBName      string `mapstructure:"mongo-dbname"`
}

func main() {
	hostname, _ := os.Hostname()

	viper.SetDefault(cfgFlagHost, "")
	viper.SetDefault(cfgFlagHostname, hostname)
	viper.SetDefault(cfgFlagLogLevel, "info")

	viper.SetDefault(cfgFlagNatsURL, "")
	viper.SetDefault(cfgFlagRedisURL, "")
	viper.SetDefault(cfgFlagRedisClusterMode, false)
	viper.SetDefault(cfgFlagMongoURL, "")
	viper.SetDefault(cfgFlagMongoDBName, "")

	viper.SetEnvPrefix(envVarPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	pflag.String(cfgFlagHost, "", "host address to bind server")
	pflag.String(cfgFlagLogLevel, "", "log level")
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	// Logger

	logger, _ := api.InitZap(viper.GetString(cfgFlagLogLevel))
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	var srvCfg ServerConfig
	if err := viper.Unmarshal(&srvCfg); err != nil {
		logger.Panic("config unmarshal failed", zap.Error(err))
	}
	fmt.Printf("%+v\n", srvCfg)

	logger.Info("spawner configuration", zap.String("config", fmt.Sprintf("%+v", srvCfg)))

	// Make Coordinator Spawner
	spawner, err := MakeCoordinatorSpanwer(srvCfg, logger)
	if err != nil {
		logger.Panic("error initialising api server", zap.Error(err))
	}

	// Run Coordinator Spawner
	go func() {
		logger.Info("starting spawner", zap.String("host", srvCfg.Host))
		if err := spawner.Run(); err != nil {
			logger.Panic("error starting spawner", zap.Error(err))
			os.Exit(1)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
}
