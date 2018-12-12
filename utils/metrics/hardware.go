package metrics

import (
	"runtime"
	"time"
)

type HardwareStats struct {
	Memory MemoryStats `json:"memory"`
	Cpu    CpuStats    `json:"cpu"`
}

type MemoryStats struct {
	Total           string  `json:"total"`
	Used            float32 `json:"used"`
	Idle            float32 `json:"idle"`
	ProcessAlloc    string  `json:"process_alloc"`
	ProcessMaxAlloc string  `json:"process_max_alloc"`
}

type CpuStats struct {
	MaxThreads int     `json:"max_threads"`
	Used       float32 `json:"used"`
	Idle       float32 `json:"idle"`
}

type updateThrottle struct {
	ticker *time.Ticker
	cache  HardwareStats
}

var throttle *updateThrottle

func (t *updateThrottle) start() {
	t.cache = getData()

	go func(t *updateThrottle) {
		for {
			select {
			case <-t.ticker.C:
				// log().Trace("refreshing the modules data from filesystem")
				t.cache = getData()
			}
		}
	}(t)
}

func SetThrottle(seconds int) {
	if throttle != nil {
		throttle.ticker.Stop()
		throttle = nil
	}
	if seconds <= 0 {
		return
	}
	t := &updateThrottle{
		ticker: time.NewTicker(time.Duration(seconds) * time.Second),
	}
	t.start()
	throttle = t
}

func GetHardwareStats() HardwareStats {
	if throttle != nil {
		return throttle.cache
	}
	return getData()
}

func getData() HardwareStats {
	stats := getHardwareStats()
	stats.Cpu.MaxThreads = runtime.NumCPU()
	appendProcessMemInfo(&stats.Memory)
	return stats
}

func appendProcessMemInfo(result *MemoryStats) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	result.ProcessAlloc = getMemoryString(float64(stats.Alloc))
	result.ProcessMaxAlloc = getMemoryString(float64(stats.TotalAlloc))
}
