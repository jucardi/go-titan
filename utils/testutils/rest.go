package testutils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/net/rest"
	"github.com/jucardi/go-titan/net/rest/router"
	"github.com/jucardi/go-titan/utils/paths"
	"google.golang.org/protobuf/proto"
)

type RouteRegistrationFunc func(r router.IRouter)

var rs router.IEngine

// Prepare should be called during init of a integration test file.  It will handle
// invoking the RouterCreateFunc to create the initial routes/stack as defined in the
// calling microservice.  It will setup any routes by invoking the args containing
// RouteRegistrationFunc.
func Prepare(router router.IEngine, routes ...RouteRegistrationFunc) {
	rs = router

	for _, route := range routes {
		route(router)
	}
}

// Serve starts the server for testing and returns the ResponseRecorder used for asserts
func Serve(req *http.Request) *ResponseRecorder {
	resp := NewResponseRecorder()
	rs.ServeHTTP(resp, req)
	return resp
}

// ContextURI returns the current ContextPath + specified path
func ContextURI(apiVersion int, path string) string {
	return paths.Combine("/", rs.ContextPath(), fmt.Sprintf("v%d", apiVersion), path)
}

// NewJsonGetRequest returns a new GET request with the `Content-Type:application/json` header to indicate
// that the response should be expected in a Json format
func NewJsonGetRequest(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(rest.HeaderContentType, rest.ContentTypeJson)
	return req
}

// NewProtoGetRequest returns a new GET request with the `Content-Type:application/x-protobuf` header to indicate
// that the response should be expected in a Protobuf format
func NewProtoGetRequest(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(rest.HeaderContentType, rest.ContentTypeProto)
	return req
}

// RequestFromJsonFile loads a JSON file into a request body, and creates an http.Request with the
// provided HTTP method and URL. Automatically sets the 'Content-Type' header to 'application/json'
//
//    - {method}: Indicates the HTTP method
//    - {url}:    The URL to send the request to
//    - {file}:   The file to load into the request body
//    - {token}:  The authorization token to be used. Add Bearer of no token type is provided
//
func RequestFromJsonFile(method, url, file string, token ...string) *http.Request {
	reader := FileToReader(file)
	ret := RequestFromReader(method, url, reader, token...)
	ret.Header.Set(rest.HeaderContentType, rest.ContentTypeJson)
	return ret
}

// RequestFromProtoObj serializes the proto message, uses the resulting byte array as the request body, and creates an http.Request with the
// provided HTTP method and URL. Automatically sets the 'Content-Type' header to 'application/x-protobuf'
//
//    - {method}: Indicates the HTTP method
//    - {url}:    The URL to send the request to
//    - {body}:   The protobuf message
//    - {token}:  The authorization token to be used. Add Bearer of no token type is provided
//
func RequestFromProtoObj(method, url string, body proto.Message, token ...string) *http.Request {
	data, err := proto.Marshal(body)
	logx.WithObj(err).Fatal("failed to marshal into protobuf message")
	ret := RequestFromBytes(method, url, data, token...)
	ret.Header.Add(rest.HeaderContentType, rest.ContentTypeProto)
	return ret
}

// RequestFromBytes loads byte array data as the request body, and creates an http.Request with the
// provided HTTP method and URL.
//
//    - {method}: Indicates the HTTP method
//    - {url}:    The URL to send the request to
//    - {body}:   The byte array for the request body
//    - {token}:  The authorization token to be used. Add Bearer of no token type is provided
//
func RequestFromBytes(method, url string, body []byte, token ...string) *http.Request {
	return RequestFromReader(method, url, bytes.NewBuffer(body), token...)
}

// RequestNoBody creates a new  http.Request with the provided HTTP method and URL and an empty request
// body.
//
//    - {method}: Indicates the HTTP method
//    - {url}:    The URL to send the request to
//    - {token}:  The authorization token to be used. Add Bearer of no token type is provided
//
func RequestNoBody(method, url string, token ...string) *http.Request {
	return RequestFromReader(method, url, nil, token...)
}

// RequestFromReader loads the provided io.Reader as request body, and creates an http.Request with the
// provided HTTP method and URL.
//
//    - {method}: Indicates the HTTP method
//    - {url}:    The URL to send the request to
//    - {body}:   The request body
//    - {token}:  The authorization token to be used. Add Bearer of no token type is provided
//
func RequestFromReader(method, url string, body io.Reader, token ...string) *http.Request {
	ret, err := http.NewRequest(method, url, body)
	logx.WithObj(err).Fatal("failed to create request")

	if len(token) > 0 && token[0] != "" {
		t := token[0]
		if len(strings.Split(t, " ")) <= 1 {
			t = "Bearer " + t
		}
		ret.Header.Set(rest.HeaderAuthorization, t)
	}
	return ret
}
