package configx

import (
	"github.com/jucardi/go-titan/definitions/logging"
)

var (
	defaultLogger = logging.New()
	logger        logging.ILogger
)

// SetLogger sets an implementation of ILogger to be used as the logger for the
// configx package
func SetLogger(l logging.ILogger) {
	logger = l
}

// LogCallbacks returns a handler that allows to register individual callbacks
// to be used by the beans package to report errors, info messages and/or debug
// messages
func LogCallbacks() logging.ILogCallback {
	return defaultLogger
}

func log() logging.ILogger {
	if logger != nil {
		return logger
	}
	return defaultLogger
}
