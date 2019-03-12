package endpoints

import (
	"expvar"
	"github.com/jucardi/go-titan/net/rest/middleware/metrics"
	"time"

	"github.com/gin-gonic/gin"
)

// AddLogLevel adds the `/metrics` endpoint to the given router.
func AddMetrics(router *gin.Engine) {
	router.GET("/memory-details", gin.WrapH(expvar.Handler()))
	router.GET("/metrics", func(context *gin.Context) {
		start := time.Now()
		resp := map[string]interface{}{
			"hardware":  metrics.GetHardwareStats(),
			"resources": metrics.GetStats(),
		}
		latency := time.Since(start)
		resp["metrics_latency"] = latency.String()
		context.IndentedJSON(200, resp)
	})
}
