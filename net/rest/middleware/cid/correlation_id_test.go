package cid

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jucardi/go-testx/assert"
	"github.com/jucardi/go-titan/net/rest"
)

const (
	testUri = "/test"
)

func TestCorrelationIdentifierResponse(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, testUri, nil)

	router := createRouter()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code, "Expecting status code: 204")

	// Make sure we are including the X-CID in the header response for all requests
	assert.True(t, res.Header().Get(HeaderCorrelationId) != "", "Expecting X-CID header")
}

func TestCorrelationIdentifierIsSameAsRequest(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, testUri, nil)
	req.Header.Set(HeaderCorrelationId, "some-cid")

	router := createRouter()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code, "Expecting status code: 204")

	// Make sure we the response has our X-CID requests prefix
	c := res.Header().Get(HeaderCorrelationId)
	assert.True(t, strings.HasPrefix(c, "some-cid"), "Expecting X-CID header to prefix request")
}

func createRouter() *gin.Engine {
	router := gin.New()
	router.Use(func(context *gin.Context) {
		Handler(rest.NewContext(context, false))
	})
	router.GET(testUri, func(c *gin.Context) {
		c.String(200, "OK")
	})
	return router
}
