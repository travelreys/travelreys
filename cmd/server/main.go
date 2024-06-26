package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/travelreys/travelreys/pkg/api"
	"go.uber.org/zap"
)

const (
	cfgFlagHost         = "host"
	cfgFlagHostname     = "hostname"
	cfgFlagPort         = "port"
	cfgFlagLogLevel     = "log-level"
	cfgFlagCORSOrigin   = "cors-origin"
	cfgFlagSecureCookie = "secure-cookie"

	envVarPrefix = "TRAVELREYS"
)

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	CORSOrigin   string `mapstructure:"cors-origin"`
	SecureCookie bool   `mapstructure:"secure-cookie"`
}

func (cfg ServerConfig) HTTPBindAddress() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

func main() {
	hostname, _ := os.Hostname()

	viper.SetDefault(cfgFlagHost, "")
	viper.SetDefault(cfgFlagHostname, hostname)
	viper.SetDefault(cfgFlagPort, "2022")
	viper.SetDefault(cfgFlagLogLevel, "info")
	viper.SetDefault(cfgFlagCORSOrigin, "*")
	viper.SetDefault(cfgFlagSecureCookie, true)

	viper.SetEnvPrefix(envVarPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	pflag.String(cfgFlagHost, "", "host address to bind server")
	pflag.String(cfgFlagPort, "", "http server port")
	pflag.String(cfgFlagLogLevel, "", "log level")
	pflag.Bool(cfgFlagSecureCookie, true, "secure cookie")
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

	logger.Info("server configuration", zap.String("config", fmt.Sprintf("%+v", srvCfg)))

	// Make Servers
	apiSrv, err := MakeAPIServer(srvCfg, logger)
	if err != nil {
		logger.Panic("error initialising api server", zap.Error(err))
	}

	// Run Servers

	go func() {
		logger.Info("starting api server",
			zap.String("host", srvCfg.Host),
			zap.String("port", srvCfg.Port),
		)
		if err := http.ListenAndServe(srvCfg.HTTPBindAddress(), apiSrv.Handler); err != nil {
			logger.Panic("error starting api server", zap.Error(err))
			os.Exit(1)
		}
	}()

	// graceful shutdown
	stopCh := api.SetupSignalHandler()
	sd, _ := api.NewShutdown(logger)
	sd.Graceful(stopCh, apiSrv)
}
