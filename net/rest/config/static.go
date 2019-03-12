package config

var (
	singleton = newConfig()
	callbacks []func(config *RestConfig)
)

func Rest() *RestConfig {
	if singleton == nil {
		singleton = newConfig()
		triggerCallbacks(singleton)
	}
	return singleton
}

func Set(cfg *RestConfig) {
	if cfg == nil {
		return
	}
	singleton = cfg
	triggerCallbacks(singleton)
}

func AddReloadCallback(callback func(config *RestConfig)) {
	callbacks = append(callbacks, callback)
}

func triggerCallbacks(config *RestConfig) {
	for _, callback := range callbacks {
		callback(config)
	}
}

func newConfig() *RestConfig {
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
