package logx

import (
	"math/bits"
	"strings"

	"github.com/jucardi/go-strings/stringx"
)

// These are the different logging levels.
const (
	// LevelPanic level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	LevelPanic Level = 0x1 << iota

	// LevelFatal level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	LevelFatal

	// LevelError level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LevelError

	// LevelWarn level. Non-critical entries that deserve eyes.
	LevelWarn

	// LevelInfo level. General operational entries about what's going on inside the
	// application.
	LevelInfo

	// LevelDebug level. Usually only enabled when debugging. Very verbose logging.
	LevelDebug

	// LevelTrace level. Designates finer-grained informational events than the Debug.
	LevelTrace
)

const (
	LevelAll       = LevelErrors | LevelNonErrors
	LevelErrors    = LevelError | LevelFatal | LevelPanic
	LevelNonErrors = LevelTrace | LevelDebug | LevelInfo | LevelWarn
)

// Level type
type Level uint32

// IsActive is a convenient utility that indicates whether a provided log level is active based on an existing log level
// configuration.
//
//   Eg:  DefaultLogger.GetLevel().IsActive(logx.LevelDebug)
//
// In the example above, if Debug level logging is enabled, it will return true (same if the argument passed is Info,
// Warn, Error or Panic)
func (l Level) IsActive(level Level) bool {
	return l >= level
}

// Convert the Level to a string. E.g. LevelPanic becomes "PANIC".
func (l Level) String() string {
	var ret []string

	for _, lvl := range l.Split() {
		switch lvl {
		case LevelTrace:
			ret = append(ret, "trace")
		case LevelDebug:
			ret = append(ret, "debug")
		case LevelInfo:
			ret = append(ret, "info")
		case LevelWarn:
			ret = append(ret, "warn")
		case LevelError:
			ret = append(ret, "error")
		case LevelFatal:
			ret = append(ret, "fatal")
		case LevelPanic:
			ret = append(ret, "panic")
		}
	}

	if len(ret) == 0 {
		return "unknown"
	}
	return strings.Join(ret, ",")
}

// Priority returns the highest priority of level contained in the value
func (l Level) Priority() Level {
	return Level(0x1 << uint(bits.TrailingZeros(uint(l))))
}

// Split returns the individual log level values contained by a single level instance
func (l Level) Split() []Level {
	var ret []Level
	if l&LevelTrace == LevelTrace {
		ret = append(ret, LevelTrace)
	}
	if l&LevelDebug == LevelDebug {
		ret = append(ret, LevelDebug)
	}
	if l&LevelInfo == LevelInfo {
		ret = append(ret, LevelInfo)
	}
	if l&LevelWarn == LevelWarn {
		ret = append(ret, LevelWarn)
	}
	if l&LevelError == LevelError {
		ret = append(ret, LevelError)
	}
	if l&LevelFatal == LevelFatal {
		ret = append(ret, LevelFatal)
	}
	if l&LevelPanic == LevelPanic {
		ret = append(ret, LevelPanic)
	}
	return ret
}

// ParseLevel takes a string level and returns the log level constant.
func ParseLevel(lvl string) Level {
	var ret Level
	for _, val := range strings.Split(lvl, ",") {
		switch stringx.New(val).ToLower().TrimSpace().S() {
		case "all":
			return LevelAll
		case "errors":
			ret = ret | LevelErrors
		case "non-errors":
			ret = ret | LevelNonErrors
		case "panic":
			ret = ret | LevelPanic
		case "fatal":
			ret = ret | LevelFatal
		case "error":
			ret = ret | LevelError
		case "warn", "warning":
			ret = ret | LevelWarn
		case "info":
			ret = ret | LevelInfo
		case "debug":
			ret = ret | LevelDebug
		case "trace":
			ret = ret | LevelTrace
		}
	}

	return ret
}
