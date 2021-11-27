package recovery

import (
	"runtime"

	"github.com/jucardi/go-titan/logx"
)

const StackTraceSize = 4096

// RecoverHandler defines the contract for a function to be triggered when a recovery occurs
type RecoverHandler func(rval interface{})

type dump struct {
	Error      interface{} `json:"error"                 yaml:"yaml"`
	StackTrace string      `json:"stack_trace,omitempty" yaml:"stack_trace,omitempty"`
}

var (
	handlers []RecoverHandler
)

// Recovery handles catching unexpected panics in the process. This should be use for processes outside
// of an HTTP request if using the router provided in the `net/router` package. For HTTP panic handling,
// if using the router in `net/router`, it already registers a middleware panic handler.
//
// You must use "defer Recovery()" before your main loop
//
func Recovery() {
	if rval := recover(); rval != nil {
		logx.WithObj(dump{Error: rval, StackTrace: Stack()}).Error("PANIC RECOVERED")
		TriggerHandlers(rval)
	}
}

// TriggerHandlers executes all registered recovery handlers providing the recovery object obtained.
// Recovery automatically does this, use this only when manually needing to do `rval := recover()1
func TriggerHandlers(rval interface{}) {
	for _, h := range handlers {
		h(rval)
	}
}

// Register registers a handler that will be triggered in the event of a recovery.
func Register(handler RecoverHandler) {
	handlers = append(handlers, handler)
}

// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack() string {
	buf := make([]byte, StackTraceSize)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}
