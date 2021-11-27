package cid

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jucardi/go-titan/net/rest"
)

const (
	HeaderCorrelationId    = "X-CID"
	HeaderCorrelationTrace = "X-CID-Trace"
	correlationIdStore     = "cid"
	correlationTraceStore  = "cid-trace"
)

var curProcessName string

// Handler is handler that will capture and store the X-CID header on every
// inbound request, appending this app/service binary name as part of the sequence chain.
//
// If the header doesn't exist it will generate a new identifier and
// add this application as the first chain in the sequence.
//
// The correlation identifier has the following format:
// <uuid>|app1.app2.app3 ...
func Handler(c *rest.Context) {
	cidHandler(c)
}

// GetCid returns the current correlation identifier if it was set during the current
// request scope otherwise it generates a new one.
func GetCid(c *rest.Context) string {
	if correlationId, exists := c.Get(correlationIdStore); exists {
		return correlationId.(string)
	}
	return cidHandler(c)
}

// Wrap will add the X-CID header to the specified request. It returns
// the same request for function chaining
func Wrap(c *rest.Context, r http.Request) http.Request {
	r.Header.Set(HeaderCorrelationId, GetCid(c))
	return r
}

func cidHandler(c *rest.Context) string {
	cid, trace := getCorrelationId(c), addTrace(c)
	c.Writer.Header().Set(HeaderCorrelationId, cid)
	c.Writer.Header().Set(HeaderCorrelationTrace, trace)
	c.Set(correlationIdStore, cid)
	c.Set(correlationTraceStore, trace)
	c.Next()
	return cid
}

func getCorrelationId(c *rest.Context) string {
	correlationId := c.Request.Header.Get(HeaderCorrelationId)
	if correlationId == "" {
		correlationId = uuid.New().String()
	}
	return correlationId
}

func addTrace(c *rest.Context) string {
	trace := c.Request.Header.Get(HeaderCorrelationTrace)
	if trace == "" {
		return processName()
	}

	return fmt.Sprintf("%s | %s", trace, processName())
}

// TODO: Change to be configurable
func processName() string {
	if curProcessName == "" {
		s := strings.Split(os.Args[0], "/")
		process := s[len(s)-1]
		curProcessName = process
	}
	return curProcessName
}
