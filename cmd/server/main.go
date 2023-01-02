package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/awhdesmond/tiinyplanet/pkg/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	cfgFlagHost     = "host"
	cfgFlagHostname = "hostname"
	cfgFlagPort     = "port"
	cfgFlagGRPCPort = "grpc-port"
	cfgFlagLogLevel = "log-level"

	cfgFlagCORSOrigin = "cors-origin"

	cfgFlagMongoURL         = "mongo-url"
	cfgFlagMongoDBName      = "mongo-dbname"
	cfgFlagNatsURL          = "nats-url"
	cfgFlagRedisURL         = "redis-url"
	cfgFlagRedisClusterMode = "redis-cluster-mode"

	envVarPrefix = "TIINYPLANET"
)

type ServerConfig struct {
	GRPCPort string `mapstructure:"grpc-port"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`

	CORSOrigin string `mapstructure:"cors-origin"`

	NatsURL          string `mapstructure:"nats-url"`
	RedisURL         string `mapstructure:"redis-url"`
	RedisClusterMode bool   `mapstructure:"redis-cluster-mode"`
	MongoURL         string `mapstructure:"mongo-url"`
	MongoDBName      string `mapstructure:"mongo-dbname"`
}

func (cfg ServerConfig) HTTPBindAddress() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

func (cfg ServerConfig) GRPCBindAddress() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.GRPCPort)
}

func main() {
	hostname, _ := os.Hostname()

	viper.SetDefault(cfgFlagHost, "")
	viper.SetDefault(cfgFlagHostname, hostname)
	viper.SetDefault(cfgFlagPort, "2022")
	viper.SetDefault(cfgFlagGRPCPort, "2023")
	viper.SetDefault(cfgFlagLogLevel, "info")
	viper.SetDefault(cfgFlagCORSOrigin, "*")
	viper.SetDefault(cfgFlagNatsURL, "")
	viper.SetDefault(cfgFlagRedisURL, "")
	viper.SetDefault(cfgFlagRedisClusterMode, false)
	viper.SetDefault(cfgFlagMongoURL, "")
	viper.SetDefault(cfgFlagMongoDBName, "")

	viper.SetEnvPrefix(envVarPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	pflag.String(cfgFlagHost, "", "host address to bind server")
	pflag.String(cfgFlagPort, "", "http server port")
	pflag.String(cfgFlagGRPCPort, "", "grpc server port")
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	// Logger

	logger, _ := utils.InitZap(viper.GetString(cfgFlagLogLevel))
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	var srvCfg ServerConfig
	if err := viper.Unmarshal(&srvCfg); err != nil {
		logger.Panic("config unmarshal failed", zap.Error(err))
	}

	logger.Info("server configuration", zap.String("config", fmt.Sprintf("%+v", srvCfg)))

	// Make Servers
	grpcSrv, err := MakeCollabServer(srvCfg, logger)
	if err != nil {
		logger.Panic("error initialising grpc server", zap.Error(err))
	}
	grpcListener, err := net.Listen("tcp", srvCfg.GRPCBindAddress())
	if err != nil {
		logger.Panic("error creating grpc listener", zap.Error(err))
	}

	go func() {
		logger.Info("starting grpc server",
			zap.String("host", srvCfg.Host),
			zap.String("port", srvCfg.GRPCPort),
		)
		if err := grpcSrv.Serve(grpcListener); err != nil {
			logger.Panic("error starting grpc server", zap.Error(err))
			os.Exit(1)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
}
