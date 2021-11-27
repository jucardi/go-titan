package encoding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/net/errorx"
	"github.com/jucardi/go-titan/net/rest"
	"github.com/jucardi/go-titan/utils/reflectx"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

// NewResponseDecoder creates a new deserializer handler for API responses
func NewResponseDecoder(resp IResponse, fallbackMode ...Encoding) IResponseDecoder {
	mode := Json
	if len(fallbackMode) > 0 {
		mode = fallbackMode[0]
	}
	return &responseDeserializer{IResponse: resp, fallback: mode}
}

type responseDeserializer struct {
	IResponse
	fallback Encoding
}

// Unmarshal attempts to deserialize the response based on the content type.
//
// If the HTTPStatus code is > 400 it will attempt to deserialize the response body to an instance
// of `*errorx.Error` which will be returned as the return object.
//
// It also returns `*errorx.Error` if unmarshalling the response body fails
//
func (r *responseDeserializer) Unmarshal(obj interface{}, contentType ...string) error {
	target := obj

	if r.GetCode() >= 400 {
		target = &errorx.Error{}
	}

	respBytes, err := r.BodyBytes()
	if err != nil {
		return r.appendRequestDetails(
			errorx.WrapUnhandledf(err, "Unable to read response body  (Returned HTTP status: %d)", r.GetCode()),
		)
	}

	// If the response body is empty
	if len(respBytes) == 0 {
		if r.GetCode() >= 400 {
			return r.appendRequestDetails(
				errorx.New(r.GetCode(), http.StatusText(r.GetCode()), "Obtained error code with empty response body from the server"),
			)
		}
		return nil
	}

	if reflectx.IsNil(target) {
		return r.appendRequestDetails(
			errorx.Wrapf(err, "Failed to unmarshal response, target cannot be nil   (Returned HTTP status was: %d)", r.GetCode()),
		)
	}

	var ct string
	if len(contentType) > 0 {
		ct = contentType[0]
	} else {
		ct = r.Headers().Get(rest.HeaderContentType)
	}

	if msg, ok := target.(proto.Message); ok {
		// If the target is a protobuf message.

		if strings.Contains(ct, rest.ContentTypeProto) {
			// Verifies the Content-Type, if set and matches a the Protobuf content type, attempts to unmarshal as Proto
			err = proto.Unmarshal(respBytes, msg)
		} else if strings.Contains(ct, rest.ContentTypeJson) {
			// If the target object is a protobuf message but the Content-Type is JSON, uses the jsonpb package to unmarshal as JSON
			err = protojson.Unmarshal(respBytes, msg)
		} else if err = proto.Unmarshal(respBytes, msg); err != nil {
			// If the Content-Type was unrecognized and the target is a protobuf message, attempts proto.Unmarshal and if fails, attempts jsonpb.Unmarshal
			logx.Warnf("unrecognized Content-Type (%s), target is Protobuf but proto.Unmarshal failed. Attempting JSON with jsonpb.Unmarshal", ct)
			err = protojson.Unmarshal(respBytes, msg)
		}
	} else {
		// If the target is not a protobuf message

		if strings.Contains(ct, rest.ContentTypeProto) {
			// Logs an error if the Content-Type was set as proto but the target is not a protobuf message.
			logx.Errorf("detected Content-Type %s but the target element is not a protobuf message. Attempting to deserialize as json", rest.ContentTypeProto)
		} else if !strings.Contains(ct, rest.ContentTypeJson) {
			// Logs a warning if the Content-Type was unrecognized
			logx.Warnf("unrecognized Content-Type (%s), attempting to unmarshal as JSON", ct)
		}

		// Attempts to unmarshal as JSON
		err = json.Unmarshal(respBytes, target)
	}

	if err != nil {
		ret := errorx.New(r.GetCode(), "Failed to unmarshal response", "Failed to unmarshal response", err)
		ret.Fields = map[string]string{
			"response_body": string(respBytes),
		}
		return r.appendRequestDetails(ret)
	}

	if e, ok := target.(*errorx.Error); ok {
		return r.appendRequestDetails(e)
	}

	return nil
}

// UnmarshalError returns an error instance representing an error response obtained by the server.
// Returns nil if the status code was not an error status
//
// Attempts to deserialize the response body to an instance of `*errorx.Error`
//
// It also returns `*errorx.Error` if unmarshalling the response body fails
//
func (r *responseDeserializer) UnmarshalError() error {
	target := &errorx.Error{}
	if r.GetCode() < 400 {
		return nil
	}

	respBytes, err := r.BodyBytes()
	if err != nil {
		return errorx.WrapUnhandled(err, "(Returned status: %d) Unable to read response body")
	}

	// If the response body is empty
	if len(respBytes) == 0 {
		return errorx.New(r.GetCode(), http.StatusText(r.GetCode()), "Obtained error code with empty response body from the server")
	}

	if strings.Contains(r.Headers().Get(rest.HeaderContentType), rest.ContentTypeProto) {
		err = proto.Unmarshal(respBytes, target)
	} else {
		err = protojson.Unmarshal(respBytes, target)
	}

	if err == nil && target.Code > 0 {
		return target
	}

	ret := errorx.New(r.GetCode(), http.StatusText(r.GetCode()), string(respBytes))
	return ret
}

// UnmarshalJson attempts to deserialize a JSON response into the given obj.
func (r *responseDeserializer) UnmarshalJson(obj interface{}) error {
	respBytes, err := r.BodyBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(respBytes, obj)
}

// UnmarshalYaml attempts to deserialize a YAML response into the given obj.
func (r *responseDeserializer) UnmarshalYaml(obj interface{}) error {
	respBytes, err := r.BodyBytes()
	if err != nil {
		return err
	}
	return yaml.Unmarshal(respBytes, obj)
}

// UnmarshalProto attempts to deserialize a Protobuf Message response into the given obj.
func (r *responseDeserializer) UnmarshalProto(obj proto.Message) error {
	respBytes, err := r.BodyBytes()
	if err != nil {
		return err
	}
	return proto.Unmarshal(respBytes, obj)
}

func (r *responseDeserializer) appendRequestDetails(err *errorx.Error) *errorx.Error {
	req := r.Request()
	url := ""
	headers := ""

	if req == nil {
		return err
	}

	if req.URL != nil {
		url = req.URL.String()
	}

	var hs []string
	for k, v := range req.Header {
		hs = append(hs, fmt.Sprintf("%s:%s", k, strings.Join(v, ",")))
	}

	headers = strings.Join(hs, " | ")

	err.Fields = map[string]string{
		"url":     url,
		"uri":     req.RequestURI,
		"host":    req.Host,
		"method":  req.Method,
		"addr":    req.RemoteAddr,
		"headers": headers,
	}

	return err
}
