package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/go-titan/net/rest"
	"github.com/jucardi/go-titan/net/rest/config"
	"github.com/jucardi/go-titan/net/rest/middleware/cid"
	"github.com/jucardi/go-titan/net/rest/middleware/headers"
	"github.com/jucardi/go-titan/net/rest/middleware/limits"
	"github.com/jucardi/go-titan/net/rest/middleware/logging"
	"github.com/jucardi/go-titan/net/rest/middleware/metrics"
	"github.com/jucardi/go-titan/net/rest/middleware/recovery"
)

var (
	httpRewindableBody = false
)

// UseCommonMiddleware applies the common middleware we use in microservices to the specified engine.
// The middleware added is Recover, Logging, Handler and Correlation ID
func UseCommonMiddleware(router IRouter) {
	router.Use(
		limits.Handler,
		logging.Handler,
		metrics.Handler,
		recovery.Handler,
		headers.Handler,
		cid.Handler,
	)
}

func init() {
	config.AddReloadCallback(func(config *config.RestConfig) {
		httpRewindableBody = config.Reporting.Body
	})
}

func wrapContext(context *gin.Context) *rest.Context {
	return rest.NewContext(context, httpRewindableBody)
}

func convertHandlers(handlers []HandlerFunc) []gin.HandlerFunc {
	if x := streams.From(handlers).Map(func(i interface{}) interface{} {
		h := i.(HandlerFunc)
		return h.toGinHandler()
	}).ToArray(); x != nil {
		return x.([]gin.HandlerFunc)
	}
	return nil
}

func mergeHandlerGroups(handlers ...[]HandlerFunc) []gin.HandlerFunc {
	var result []HandlerFunc
	for _, group := range handlers {
		for _, h := range group {
			result = append(result, h)
		}
	}
	return convertHandlers(result)
}
