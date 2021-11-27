package errorx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jucardi/go-strings/stringx"
	"github.com/jucardi/go-titan/errors"
	"github.com/jucardi/go-titan/runtime"
	"github.com/jucardi/go-titan/utils/reflectx"
)

// StatusText returns a text for the HTTP status code.
func (m *Error) StatusText() string {
	if m == nil || m.Code <= 0 {
		return ""
	}
	return http.StatusText(int(m.Code))
}

// Error returns the string representation of the error
func (m *Error) Error() string {
	if m == nil || m.Code <= 0 {
		return ""
	}
	return fmt.Sprintf("%d %s%s%s",
		m.Code,
		func() string {
			if m.Title == "" {
				return ""
			}
			return m.Title + ", "
		}(),
		m.Message,
		func() string {
			if len(m.Inner) == 0 {
				return ""
			}
			var details []string
			for _, d := range m.Inner {
				details = append(details, d.Error)
			}
			return " | Details: " + strings.Join(details, " ; ")
		}())
}

// Caller returns the caller information of the error (filename and line)
func (m *Error) Stack() []string {
	return m.Trace
}

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

	var (
		details []*InnerError
		errMsgs []string
	)

	if len(title) == 0 && len(errs) > 1 {
		ret.Title = "multiple errors occurred"
	}

	for _, e := range errs {
		if e == nil {
			continue
		}

		detail := &InnerError{
			Error: e.Error(),
		}

		if withCaller, ok := e.(errors.IErrorWithStack); ok {
			stack := withCaller.Stack()
			if len(stack) > 0 {
				detail.Caller = stack[0]
			}
			ret.Trace = append(ret.Trace, withCaller.Stack()...)
		}

		if data, e := json.Marshal(e); e != nil {
			detail.Details = string(data)
		}
		if len(msg) == 0 {
			errMsgs = append(errMsgs, e.Error())
		}
		details = append(details, detail)
	}

	if len(errMsgs) > 0 {
		ret.Message = strings.Join(errMsgs, "; ")
	}

	ret.Inner = details
	return ret
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

	for k, idFunc := range errorTypeIdentifiers {
		if idFunc(err.Error()) {
			return newError(skip+1, k, http.StatusText(k), msgx, err)
		}
	}

	return newError(skip+1, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), msgx, err)
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
		v.Title = http.StatusText(errType)
		return v
	}
	msgx := stringx.GetOrDefault("", msg...)

	return newError(skip+1, errType, http.StatusText(errType), msgx, err)
}
