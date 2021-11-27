package encoding

import (
	"net/http"

	"google.golang.org/protobuf/proto"
)

// IResponseDecoder defines the contract for an API response decoder
type IResponseDecoder interface {
	// Unmarshal attempts to deserialize the response based on the content type.
	//
	// If the HTTPStatus code is > 400 it will attempt to deserialize the response body to an instance
	// of `*errorx.Error` which will be returned as the return object.
	//
	// It also returns `*errorx.Error` if unmarshalling the response body fails
	//
	Unmarshal(obj interface{}, contentType ...string) error

	// UnmarshalError returns an error instance representing an error response obtained by the server.
	// Returns nil if the status code was not an error status
	//
	// Attempts to deserialize the response body to an instance of `*errorx.Error`
	//
	// It also returns `*errorx.Error` if unmarshalling the response body fails
	//
	UnmarshalError() error

	// UnmarshalJson attempts to deserialize a JSON response into the given obj.
	UnmarshalJson(obj interface{}) error

	// UnmarshalYaml attempts to deserialize a YAML response into the given obj.
	UnmarshalYaml(obj interface{}) error

	// UnmarshalProto attempts to deserialize a Protobuf Message response into the given obj.
	UnmarshalProto(obj proto.Message) error
}

// IResponse defines the required functions a response object needs to implement to work with the
// response deserializer defined in IResponseDecoder
type IResponse interface {
	// GetCode returns the status code of the response
	GetCode() int

	// BodyBytes returns a byte array obtained from the response body
	BodyBytes() ([]byte, error)

	// Headers returns the headers contained in a response object
	Headers() http.Header

	// Request returns the request used to obtain this response
	Request() *http.Request
}
