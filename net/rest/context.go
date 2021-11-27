package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/net/errorx"
	"github.com/jucardi/go-titan/net/rest/config"
	"github.com/jucardi/go-titan/utils/reflectx"
)

// Context is a `gin.Context` wrapper that allows extending context functionality.
type Context struct {
	*gin.Context
	enc     Encoder
	reqBody []byte
	aborted bool
}

func NewContext(ctx *gin.Context, rewindable bool) *Context {
	ret := &Context{Context: ctx}
	ret.Writer = &respWriter{ResponseWriter: ctx.Writer}
	if rewindable {
		ret.init()
	}
	return ret
}

func (c *Context) init() {
	c.reqBody, _ = ioutil.ReadAll(c.Request.Body)
	c.RewindBody()
}

// RewindBody creates a new io.Writer with the cached request body bytes. Only works if `rewindable` was set
// to `true` when calling `NewContext`
func (c *Context) RewindBody() {
	if len(c.reqBody) == 0 {
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(c.reqBody))
}

// Copy returns a copy of the current context that can be safely used outside the request's scope.
// This has to be used when the context has to be passed to a goroutine.
func (c *Context) Copy() *Context {
	return &Context{
		Context: c.Context.Copy(),
	}
}

func (c *Context) Abort() {
	c.aborted = true
	c.Context.Abort()
}

func (c *Context) AbortWithCode(code int) {
	c.aborted = true
	c.Context.AbortWithStatus(code)
}

func (c *Context) AbortWithError(code int, err error) {
	c.Status(code)
	c.SendError(err)
	c.Abort()
}

// Send marshals the provided `obj` as the response and assigns the provided `code` and the HTTP response code.
//
// The response object will be encoded using the configured encoder (JSON by default).
//
// See `JSON` and `ProtoBuf` if a specific encoding is required regardless of the configuration.
//
func (c *Context) Send(code int, obj interface{}) {
	c.encode(code, obj)
}

// SendError aborts any pending middleware handlers and indents this error instance as a response object.
// If the error is implements the interface `errorx.IError`, it will use the `Status` value as
// HTTP Status code, otherwise a 500 HTTP code will be used.
//
// The response message will be encoded using the configured encoder (JSON by default).
//
// See `SendErrorJson` and `SendErrorProtobuf` if a specific encoding is required regardless of the configuration.
//
func (c *Context) SendError(err error) {
	c.setEncoding(encodingDefault).sendError(err)
}

// SendErrorJson aborts any pending middleware handlers and indents this error instance as a response object.
// If the error is implements the interface `errorx.IError`, it will use the `Status` value as
// HTTP Status code, otherwise a 500 HTTP code will be used.
//
// The response will be JSON encoded regardless of the encoding configuration.
//
func (c *Context) SendErrorJson(err error) {
	c.setEncoding(encoders.Json).sendError(err)
}

// SendErrorJsonIndent aborts any pending middleware handlers and indents this error instance as a response object.
// If the error is implements the interface `errorx.IError`, it will use the `Status` value as
// HTTP Status code, otherwise a 500 HTTP code will be used.
//
// The response will be JSON encoded regardless of the encoding configuration.
//
func (c *Context) SendErrorJsonIndent(err error) {
	c.setEncoding(encoders.IndentedJson).sendError(err)
}

// SendErrorProtobuf aborts any pending middleware handlers and indents this error instance as a response object.
// If the error is implements the interface `errorx.IError`, it will use the `Status()` value as
// HTTP Status code, otherwise a 500 HTTP code will be used.
//
// The response will be Protobuf encoded regardless of the encoding configuration.
//
func (c *Context) SendErrorProtobuf(err error) {
	c.setEncoding(encoders.Protobuf).sendError(err)
}

// SendOrErr marshals the provided response interface and assigns the provided status (or 200 OK if not provided)
// unless `err` is not nil, in which case, sends `err` as the response body. If `err` is `errorx.IError`, uses the
// assigned `Status()`, otherwise the response status will be (500) Internal Server Error.
//
// The response message (resp object or error) will be encoded using the configured encoder (JSON by default).
//
// See SendOrErrJson and SendOrErrProtobuf if a specific encoding is required regardless of the configuration.
//
func (c *Context) SendOrErr(resp interface{}, err error, httpStatus ...int) {
	c.sendOrErr(resp, err, httpStatus...)
}

// SendOrErrJson marshals the provided response interface and assigns the provided status (or 200 OK if not provided)
// unless `err` is not nil, in which case, sends `err` as the response body. If `err` is `errorx.IError`, uses the
// assigned `Status()`, otherwise the response status will be (500) Internal Server Error.
//
// The response message (resp object or error) will be encoded as JSON regardless of the configuration.
//
func (c *Context) SendOrErrJson(resp interface{}, err error, httpStatus ...int) {
	c.setEncoding(encoders.Json).sendOrErr(resp, err, httpStatus...)
}

// SendOrErrJsonIndent marshals the provided response interface and assigns the provided status (or 200 OK if not provided)
// unless `err` is not nil, in which case, sends `err` as the response body. If `err` is `errorx.IError`, uses the
// assigned `Status()`, otherwise the response status will be (500) Internal Server Error.
//
// The response message (resp object or error) will be encoded as JSON regardless of the configuration.
//
func (c *Context) SendOrErrJsonIndent(resp interface{}, err error, httpStatus ...int) {
	c.setEncoding(encoders.IndentedJson).sendOrErr(resp, err, httpStatus...)
}

// SendOrErrProtobuf marshals the provided response interface and assigns the provided status (or 200 OK if not provided)
// unless `err` is not nil, in which case, sends `err` as the response body. If `err` is `errorx.IError`, uses the
// assigned `Status()`, otherwise the response status will be (500) Internal Server Error.
//
// The response message (resp object or error) will be encoded as Protobuf regardless of the configuration.
//
func (c *Context) SendOrErrProtobuf(resp interface{}, err error, httpStatus ...int) {
	c.setEncoding(encoders.Protobuf).sendOrErr(resp, err, httpStatus...)
}

// StatusOrErr assigns the provided status code (or 200 OK if not provided) unless `err` is not nil, in which case,
// sends`err` as the response body. If `err` is `errorx.IError`, uses the assigned `Status()`, otherwise the
// response status will be (500) Internal Server Error. If `err` is nil, the response body will be empty
//
// The response message will be encoded using the configured encoder (JSON by default).
//
// See StatusOrErrJson and StatusOrErrProtobuf if a specific encoding is required regardless of the configuration.
//
func (c *Context) StatusOrErr(err error, httpStatus ...int) {
	c.statusOrErr(err, httpStatus...)
}

// StatusOrErrJson assigns the provided status code (or 200 OK if not provided) unless `err` is not nil, in which case,
// sends`err` as the response body. If `err` is `errorx.IError`, uses the assigned `Status()`, otherwise the
// response status will be (500) Internal Server Error. If `err` is nil, the response body will be empty
//
// The response message will be encoded as JSON regardless of the configuration.
//
func (c *Context) StatusOrErrJson(err error, httpStatus ...int) {
	c.setEncoding(encoders.Json).statusOrErr(err, httpStatus...)
}

// StatusOrErrProtobuf assigns the provided status code (or 200 OK if not provided) unless `err` is not nil, in which case,
// sends`err` as the response body. If `err` is `errorx.IError`, uses the assigned `Status()`, otherwise the
// response status will be (500) Internal Server Error. If `err` is nil, the response body will be empty
//
// The response message will be encoded as Protobuf regardless of the configuration.
//
func (c *Context) StatusOrErrProtobuf(err error, httpStatus ...int) {
	c.setEncoding(encoders.Protobuf).statusOrErr(err, httpStatus...)
}

// DumpRequest dumps the request data contained in the context.
func (c *Context) DumpRequest() string {
	c.RewindBody()
	data, _ := httputil.DumpRequest(c.Request, len(c.reqBody) > 0)
	return string(data)
}

func (c *Context) sendError(err error) {
	logx.Trace("sending error")
	var e *errorx.Error

	if v, ok := err.(*errorx.Error); ok {
		e = v
	} else {
		e = errorx.Wrap(e)
	}

	if e == nil {
		return
	}
	httpStatus := int(e.Code)

	// http.StatusText returns empty if the code is unknown. If so, assigns code 500 for being an unhandled error type.
	if http.StatusText(httpStatus) == "" {
		httpStatus = http.StatusInternalServerError
	}

	_ = c.Error(err)

	if !config.Rest().Response.ErrorBodies {
		c.AbortWithStatus(httpStatus)
		return
	}

	var ex *errorx.Error

	if !config.Rest().Response.ErrorStackTrace {
		// Cloning the error, sanitizing stack data before sending it to the client
		ex = &errorx.Error{
			Code:      e.Code,
			Timestamp: e.Timestamp,
			Title:     e.Title,
			Message:   e.Message,
		}
		if len(ex.Inner) > 0 {
			ex.Inner = streams.From(e.Inner).Map(func(i interface{}) interface{} {
				x := i.(*errorx.InnerError)
				return &errorx.InnerError{
					Error:   x.Error,
					Details: x.Details,
				}
			}).ToArray().([]*errorx.InnerError)
		}
	} else {
		ex = e
	}

	c.encode(httpStatus, ex)
	c.Abort()
}

func (c *Context) sendOrErr(resp interface{}, err error, httpStatus ...int) {
	status := http.StatusOK

	if len(httpStatus) > 0 {
		status = httpStatus[0]
	}

	if !reflectx.IsNil(err) {
		c.sendError(errorx.Wrap(err))
	} else if resp == nil {
		c.Status(status)
	} else {
		t := reflect.TypeOf(resp)
		switch t.Kind() {
		case reflect.Ptr, reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
			c.encode(status, resp)
		default:
			c.String(status, "%s", resp)
		}
	}
}

func (c *Context) statusOrErr(err error, httpStatus ...int) {
	status := http.StatusOK

	if len(httpStatus) > 0 {
		status = httpStatus[0]
	}

	if !reflectx.IsNil(err) {
		c.sendError(err)
	} else {
		c.Status(status)
	}
}

func (c *Context) setEncoding(enc Encoder) *Context {
	if c.enc == nil {
		c.enc = enc
	}
	return c
}

func (c *Context) encode(code int, obj interface{}) {
	enc := c.enc

	// Encoding not explicitly set, using encoding based on the configuration
	if enc == nil {
		if e, ok := encodingMap[config.Rest().Response.Encoding]; ok {
			enc = e
		} else {
			logx.Warn("unrecognized primary encoding configured, ", config.Rest().Response.Encoding)
			enc = encodingDefault
		}
	}

	err := enc(c, code, obj)
	if err == nil {
		return
	}
	logx.Error("failed to encode response with primary encoding, ", err.Error())
	if e, ok := encodingMap[config.Rest().Response.FallbackEncoding]; ok {
		enc = e
	} else {
		logx.Warn("unrecognized fallback encoding configured, ", config.Rest().Response.FallbackEncoding)
		enc = encodingFallback
	}
	if err := enc(c, code, obj); err != nil {
		logx.Error("failed to encode response with secondary encoding, ", err.Error())
	}
}
