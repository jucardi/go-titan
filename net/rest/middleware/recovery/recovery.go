package recovery

import (
	"fmt"

	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/net/errorx"
	"github.com/jucardi/go-titan/net/rest"
	"github.com/jucardi/go-titan/utils/recovery"
)

type dump struct {
	Error       interface{} `json:"error"                  yaml:"yaml"`
	HttpRequest string      `json:"http_request,omitempty" yaml:"http_request,omitempty"`
	StackTrace  string      `json:"stack_trace,omitempty"  yaml:"stack_trace,omitempty"`
}

// Handler is a router handler that will catch unexpected panics and recover. The recovery process will
// log information related to the panic including the stack trace and will relay an HttpError with status
// (500) Internal Server Error and with details of the error.
func Handler(c *rest.Context) {
	defer func() {
		rval := recover()
		if rval == nil {
			return
		}

		var msg string

		if e, ok := rval.(error); ok {
			msg = e.Error()
		} else {
			msg = fmt.Sprint(rval)
		}

		logx.WithObj(
			dump{
				Error:       rval,
				HttpRequest: c.DumpRequest(),
				StackTrace:  recovery.Stack(),
			}).Errorf("PANIC RECOVERED | %s", msg)

		recovery.TriggerHandlers(rval)

		var err *errorx.Error
		if v, ok := rval.(error); ok {
			err = errorx.Wrap(v)
		} else {
			err = errorx.NewUnhandled(fmt.Sprint(rval))
		}
		c.SendError(err)
	}()
	c.Next()
}
