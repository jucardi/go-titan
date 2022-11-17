package monitor

import (
	"time"
)

// New creates a new default monitor which ticks at the provided interval. This monitor is at a
// `stopped` state until `Start()` or `StartAsync()` are called
func New(updateInterval time.Duration) IMonitor {
	ret := &service{updateInterval: updateInterval}
	return ret
}
