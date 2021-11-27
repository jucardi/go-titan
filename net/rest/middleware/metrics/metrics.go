package metrics

import (
	"sync"
	"time"

	"github.com/jucardi/go-titan/net/rest"
	"github.com/jucardi/go-titan/utils/metrics"
)

const (
	LatencyContextKey = "__latency__"
)

var (
	routerStats = map[string]*RouteStats{}
	mux         sync.Mutex
)

type RouteStats struct {
	RouteData map[int]*ResponseStats `json:"route_data" yaml:"route_data"`
	mux       sync.Mutex
}

func (r *RouteStats) getStatusData(status int) *ResponseStats {
	r.mux.Lock()
	defer r.mux.Unlock()

	if data, ok := r.RouteData[status]; ok {
		return data
	}
	responseData := &ResponseStats{}
	r.RouteData[status] = responseData
	return responseData
}

type ResponseStats struct {
	Count        int64  `json:"count"         yaml:"count"`
	MaxDuration  string `json:"max_duration"  yaml:"max_duration"`
	MeanDuration string `json:"mean_duration" yaml:"mean_duration"`

	maxDuration  time.Duration
	meanDuration time.Duration
}

func getRouteData(route string) *RouteStats {
	mux.Lock()
	defer mux.Unlock()

	if data, ok := routerStats[route]; ok {
		return data
	}
	routeData := &RouteStats{
		RouteData: map[int]*ResponseStats{},
	}

	routerStats[route] = routeData
	return routeData
}

func Handler(c *rest.Context) {
	key := c.Request.URL.Path
	start := time.Now()
	c.Next()
	latency := time.Since(start)
	status := c.Writer.Status()

	// Consolidate by 2xx, 3xx, 4xx and 5xx status codes
	switch {
	case status >= 200 && status < 300:
		status = 200
	case status >= 300 && status < 400:
		status = 300
	case status >= 400 && status < 500:
		status = 400
	case status >= 500:
		status = 500
	}

	routeData := getRouteData(key)
	responseData := routeData.getStatusData(status)

	if latency > responseData.maxDuration {
		responseData.maxDuration = latency
		responseData.MaxDuration = latency.String()
	}

	totalLatency := responseData.meanDuration*time.Duration(responseData.Count) + latency
	responseData.Count += 1
	totalLatency = totalLatency / time.Duration(responseData.Count)
	responseData.meanDuration = totalLatency
	responseData.MeanDuration = totalLatency.String()

	c.Set(LatencyContextKey, latency)
}

func GetMeasuredLatency(c *rest.Context) time.Duration {
	dur, exists := c.Get(LatencyContextKey)
	if !exists {
		return time.Duration(0)
	}
	return dur.(time.Duration)
}

func GetStats() map[string]*RouteStats {
	mux.Lock()
	defer mux.Unlock()

	ret := map[string]*RouteStats{}
	for k, v := range routerStats {
		cloned := &RouteStats{
			RouteData: map[int]*ResponseStats{},
		}
		for x, y := range v.RouteData {
			val := *y
			cloned.RouteData[x] = &val
		}
		ret[k] = cloned
	}
	return ret
}

func GetHardwareStats() metrics.HardwareStats {
	return metrics.GetHardwareStats()
}
