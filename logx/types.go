package logx

type ILogger interface {
	IEntry

	// SetLevel sets the logging level
	SetLevel(level Level)
	// GetLevel gets the logging level
	GetLevel() Level
}

type IEntry interface {
	IMessageHandler
	WithObj(obj interface{}) IEntry
	WithFields(fields map[string]interface{}) IEntry
	WithField(key string, val interface{}) IEntry
}

// IMessageHandler contains the DefaultLogger functions that are meant for logging messages only.
type IMessageHandler interface {
	// Trace logs a message at level Trace on the DefaultLogger.
	Trace(args ...interface{})
	// Tracef logs a message at level Trace on the DefaultLogger.
	Tracef(format string, args ...interface{})

	// Debug logs a message at level Debug on the DefaultLogger.
	Debug(args ...interface{})
	// Debugf logs a message at level Debug on the DefaultLogger.
	Debugf(format string, args ...interface{})

	// Info logs a message at level Info on the DefaultLogger.
	Info(args ...interface{})
	// Infof logs a message at level Info on the DefaultLogger.
	Infof(format string, args ...interface{})

	// Warn logs a message at level Warn on the DefaultLogger.
	Warn(args ...interface{})
	// Warnf logs a message at level Warn on the DefaultLogger.
	Warnf(format string, args ...interface{})

	// Error logs a message at level Error on the DefaultLogger.
	Error(args ...interface{})
	// Errorf logs a message at level Error on the DefaultLogger.
	Errorf(format string, args ...interface{})

	// Fatal logs a message at level Fatal on the DefaultLogger.
	Fatal(args ...interface{})
	// Fatalf logs a message at level Fatal on the DefaultLogger.
	Fatalf(format string, args ...interface{})

	// Panic logs a message at level Panic on the DefaultLogger.
	Panic(args ...interface{})
	// Panicf logs a message at level Panic on the DefaultLogger.
	Panicf(format string, args ...interface{})
}
