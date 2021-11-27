package rest

import (
	"github.com/jucardi/go-titan/errors"
	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/net/rest/config"
)

type Encoder func(c *Context, code int, obj interface{}) error

var (
	encoders = &Encoders{}

	encodingDefault  Encoder = encoders.Auto
	encodingFallback Encoder = encoders.Json

	contentTypeMap = map[string]Encoder{
		ContentTypeJson:  encoders.Json,
		ContentTypeProto: encoders.Protobuf,
	}

	encodingMap = map[config.Encoding]Encoder{
		config.EncodingAuto:         encoders.Auto,
		config.EncodingJson:         encoders.Json,
		config.EncodingIndentedJson: encoders.IndentedJson,
		config.EncodingProto:        encoders.Protobuf,
	}
)

type Encoders struct {
}

func (e Encoders) Auto(c *Context, code int, obj interface{}) error {
	if c.aborted {
		c.Status(code)
		return nil
	}
	logx.Trace("attempting to use AUTO encoder")
	if c.Request == nil {
		return errors.New("request is nil, unable to automatically determine encoding, using default encoder")
	}

	responseType := c.Request.Header.Get(HeaderResponseType)
	contentType := c.Request.Header.Get(HeaderContentType)

	var enc Encoder

	// If a Response-Type header was provided in the request, attempts to determine the response encoding based on the Response-Type header value
	if responseType != "" {
		if e, ok := contentTypeMap[responseType]; ok {
			enc = e
		}
	}

	// If no Response-Type header was provided, attempts to determine the response encoding based on the Content-Type header of the request.
	if enc == nil {
		if e, ok := contentTypeMap[contentType]; ok {
			enc = e
		}
	}

	// If an encoding was determined from the request
	if enc != nil {
		return enc(c, code, obj)
	}

	if responseType != "" {
		return errors.New("unable to determine encoder based on the request Response-Type header: ", responseType)
	}

	if contentType != "" {
		return errors.New("unable to determine encoder based on the request Content-Type header: ", contentType)
	}

	return errors.New("failed to determine auto encoding")
}

func (Encoders) Json(c *Context, code int, obj interface{}) (_ error) {
	if c.aborted {
		c.Status(code)
		return
	}
	logx.Trace("using JSON encoder")
	c.JSON(code, obj)
	return
}

func (Encoders) IndentedJson(c *Context, code int, obj interface{}) (_ error) {
	if c.aborted {
		c.Status(code)
		return
	}
	logx.Trace("using Indented JSON encoder")
	c.IndentedJSON(code, obj)
	return
}

func (Encoders) Protobuf(c *Context, code int, obj interface{}) (_ error) {
	if c.aborted {
		c.Status(code)
		return
	}
	logx.Trace("using Protobuf encoder")
	c.ProtoBuf(code, obj)
	return
}
