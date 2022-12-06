package errorx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/go-strings/stringx"
	"github.com/jucardi/go-titan/errors"
	"github.com/jucardi/go-titan/net/rest/httpx"
	"github.com/jucardi/go-titan/runtime"
	"github.com/jucardi/go-titan/utils/reflectx"
)

var (
	DefaultMultiStatusClientErrorCode = httpx.StatusMultipleClientErrors
	DefaultMultiStatusServerErrorCode = httpx.StatusMultipleServerErrors
)

// New creates a new Error by the given title, message and error.
// - If `title` is empty and multiple errors are passed, the title will be set to a message
//   indicating that multiple errors occurred.
// - If `message` is empty, the message will be the result of joining all `err.Error()`
func New(code int, title, msg string, errs ...error) *Error {
	return newError(1, code, title, msg, errs...)
}

func newError(skip int, code int, title, msg string, errs ...error) *Error {
	ret := &Error{
		Code:      int32(code),
		Title:     title,
		Message:   msg,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Trace:     []string{runtime.Caller(skip + 1)},
	}

	details := generateInnerList(errs...)

	switch len(details) {
	case 0:
		break
	case 1:
		ret.ErrFlags = details[0].ErrFlags
		ret.Trace = details[0].Trace
		if msg == "" {
			ret.Message = details[0].Error()
		}
		break
	default:
		var (
			msgs  []string
			flags int32
		)

		streams.
			From(details).
			ForEach(func(i interface{}) {
				x := i.(*InnerError)
				msgs = append(msgs, x.Message)
				flags = flags | x.ErrFlags
			})

		ret.Inner = details
		ret.ErrFlags = flags

		if title == "" {
			ret.Title = "multiple errors occurred"
		}

		if msg == "" {
			ret.Message = strings.Join(msgs, "; ")
		}
	}

	return ret
}

func generateInnerList(errs ...error) (ret []*InnerError) {
	for _, e := range errs {
		if e == nil {
			continue
		}

		detail := &InnerError{
			Message: e.Error(),
		}

		if withCaller, ok := e.(errors.IErrorWithStack); ok {
			stack := withCaller.Stack()
			if len(stack) > 0 {
				detail.Caller = stack[0]
				detail.Trace = append(detail.Trace, stack...)
			}
		}

		if withFlags, ok := e.(errors.IErrorWithFlags); ok {
			detail.ErrFlags = int32(withFlags.Flags())
		}

		if data, e := json.Marshal(e); e != nil {
			detail.Details = string(data)
		}
		ret = append(ret, detail)
	}
	return
}

// Wrap attempts to guess the error type based on the `Error()` message from the provided error.
// Returns the error instance if it already implements IError.
// See the values in the constants: ErrCodeUnauthorized, ErrCodeNotFound, ErrCodeInvalidRequest, ErrCodeConflict for status matching.
//
//  - msg:  (Optional) Message tp append to the error, if not provided, it will use `Error()` from the provided error.
//
func Wrap(err error, msg ...string) *Error {
	return wrap(1, err, msg...)
}

func Wrapf(err error, format string, args ...interface{}) *Error {
	return wrap(1, err, fmt.Sprintf(format, args...))
}

func wrap(skip int, err error, msg ...string) *Error {
	if reflectx.IsNil(err) {
		return nil
	}

	msgx := stringx.GetOrDefault("", msg...)

	if v, ok := err.(*Error); ok && v != nil {
		if msgx != "" {
			v.Message = msgx + ", " + v.Message
		}
		if len(v.Trace) > 0 {
			v.Trace = append(runtime.StackLines(skip+1), v.Trace...)
		} else {
			v.Trace = runtime.StackLines(skip + 1)
		}

		return v
	}

	if code := Flags().ToStatus(errors.GetFlags(err)); code != 0 {
		return newError(skip+1, code, httpx.StatusText(code), msgx, err)
	}

	for k, idFunc := range errorTypeIdentifiers {
		if idFunc(err.Error()) {
			return newError(skip+1, k, httpx.StatusText(k), msgx, err)
		}
	}

	return newError(skip+1, http.StatusInternalServerError, httpx.StatusText(http.StatusInternalServerError), msgx, err)
}

func WrapWithCode(err error, code int, msg ...string) *Error {
	return wrapWithCode(1, err, code, msg...)
}

func WrapWithCodef(err error, errType int, format string, args ...interface{}) *Error {
	return wrapWithCode(1, err, errType, fmt.Sprintf(format, args...))
}

func wrapWithCode(skip int, err error, errType int, msg ...string) *Error {
	if reflectx.IsNil(err) {
		return nil
	}

	if v, ok := err.(*Error); ok {
		v.Code = int32(errType)
		v.Title = httpx.StatusText(errType)
		return v
	}
	msgx := stringx.GetOrDefault("", msg...)

	return newError(skip+1, errType, httpx.StatusText(errType), msgx, err)
}
