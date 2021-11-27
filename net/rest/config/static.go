package config

import (
	"github.com/jucardi/go-titan/configx"
	"github.com/jucardi/go-titan/logx"
)

const (
	configKey  = "rest"
	configName = "rest-cfg"
)

var (
	singleton = defaultConfig()
)

func init() {
	configx.AddOnReloadCallback(func(cfg configx.IConfig) {
		config := defaultConfig()

		logx.WithObj(
			cfg.MapToObj(configKey, config),
		).Fatal("unable to map service configuration")

		singleton = config
	}, configName)
}

// Rest returns rest configuration
func Rest() *RestConfig {
	return singleton
}

func AddReloadCallback(callback func(config *RestConfig)) {
	configx.AddOnReloadCallback(func(_ configx.IConfig) {
		callback(Rest())
	})
}

func defaultConfig() *RestConfig {
	return &RestConfig{
		HttpPort:    8080,
		AdminPort:   15000,
		ContextPath: "",
		Response: ResponseConfig{
			Encoding:         EncodingAuto,
			FallbackEncoding: EncodingJson,
			ErrorBodies:      true,
			ErrorStackTrace:  false,
		},
		Reporting:        ReportingConfig{},
		RequestLimitSize: 5242880,
	}
}
