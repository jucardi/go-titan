package configx

import (
	"fmt"
	"os"
)

/*
The logx package attempts to stay agnostic to different logging mechanisms, so it provides two choices
to handle logging messages.

The first choice is to assign an implementation of the ILogger interface
defined in this file, where a logger instance can be set for use, as long as it implements the 4 logging
functions (Debug, Info, Warn, Error).

The second choice to handle logging messages is to simply register a callback in each logging level desired
so sn external function will be called with the message when a logging event occurs. These handlers can be
registered by calling the LogCallbacks() function.

Eg.  To register an error handler:

    LogCallback().SetErrorCallback(func(msg string) {
       myCustomLogger.LogErrorMessage(msg)
    })

By default, only error messages will be printed to the Stdout
*/

type MessageCallback func(msg string)

// ILogCallback defines the contract for logging callback assignments for this package
type ILogCallback interface {
	// SetWarnCallback sets a callback that will be triggered when warning messages occur
	SetWarnCallback(callback MessageCallback)
	// SetErrorCallback sets a callback that will be triggered when an error occurs
	SetErrorCallback(callback MessageCallback)
	// SetInfoCallback sets a callback that will be triggered when an info message is generated
	SetInfoCallback(callback MessageCallback)
	// SetDebugCallback sets a callback that will be triggered when a debug message is generated
	SetDebugCallback(callback MessageCallback)
}

// ILogger defines the contract for a full logger that can be used by this package
type ILogger interface {
	// Error logs an error to the logger
	Error(args ...interface{})
	// Warn logs a warning message
	Warn(args ...interface{})
	// Info logs an info message to the logger
	Info(args ...interface{})
	// Debug logs a debug message to the logger
	Debug(args ...interface{})
}

var (
	defaultLogger = &loggingHandler{}
	logger        ILogger
)

// SetLogger sets an implementation of ILogger to be used as the logger for the
// configx package
func SetLogger(l ILogger) {
	logger = l
}

// LogCallbacks returns a handler that allows to register individual callbacks
// to be used by the beans package to report errors, info messages and/or debug
// messages
func LogCallbacks() ILogCallback {
	return defaultLogger
}

type loggingHandler struct {
	onWarnHandler  MessageCallback
	onErrHandler   MessageCallback
	onInfoHandler  MessageCallback
	onDebugHandler MessageCallback
}

func log() ILogger {
	if logger != nil {
		return logger
	}
	return defaultLogger
}

func (l *loggingHandler) SetWarnCallback(callback MessageCallback) {
	l.onWarnHandler = callback
}

func (l *loggingHandler) SetErrorCallback(callback MessageCallback) {
	l.onErrHandler = callback
}

func (l *loggingHandler) SetInfoCallback(callback MessageCallback) {
	l.onInfoHandler = callback
}

func (l *loggingHandler) SetDebugCallback(callback MessageCallback) {
	l.onDebugHandler = callback
}

func (l *loggingHandler) Warn(args ...interface{}) {
	if l.onWarnHandler != nil {
		l.onWarnHandler(fmt.Sprint(args...))
	}
}

func (l *loggingHandler) Error(args ...interface{}) {
	if l.onErrHandler != nil {
		l.onErrHandler(fmt.Sprint(args...))
	} else {
		newArgs := append([]interface{}{"ERROR (configx) - "}, args...)
		_, _ = fmt.Fprintln(os.Stderr, newArgs...)
	}
}

func (l *loggingHandler) Info(args ...interface{}) {
	if l.onInfoHandler != nil {
		l.onInfoHandler(fmt.Sprint(args...))
	}
}

func (l *loggingHandler) Debug(args ...interface{}) {
	if l.onDebugHandler != nil {
		l.onDebugHandler(fmt.Sprint(args...))
	}
}
