package testutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/jucardi/go-titan/net/encoding"
)

// NewResponseRecorder returns a ResponseRecorder which is an extension of httptest.ResponseRecorder.  This
// provides additional simplicity for decoding JSON into a struct
func NewResponseRecorder() *ResponseRecorder {
	ret := &ResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
	ret.IResponseDecoder = encoding.NewResponseDecoder(ret)
	return ret
}

// ResponseRecorder is an extension of httptest.ResponseRecorded
type ResponseRecorder struct {
	*httptest.ResponseRecorder
	encoding.IResponseDecoder
}

// Request returns the request that was sent to obtain this Response.
func (r *ResponseRecorder) Request() *http.Request {
	return r.ResponseRecorder.Result().Request
}

// GetCode returns the status code of the response
func (r *ResponseRecorder) GetCode() int {
	return r.Code
}

func (r *ResponseRecorder) BodyBytes() ([]byte, error) {
	return r.Body.Bytes(), nil
}

// Headers returns the headers contained in a response object
func (r *ResponseRecorder) Headers() http.Header {
	return r.Header()
}

// Decode will decode the current request body into the target
//
// Deprecated: Decode exists for historical compatibility
// and should not be used. Use the Unmarshal functions instead
//
func (r *ResponseRecorder) Decode(target interface{}) error {
	return json.Unmarshal(r.Body.Bytes(), target)
}
