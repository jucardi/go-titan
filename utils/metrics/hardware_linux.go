package metrics

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	memInfoLoc = "/proc/meminfo"
	cpuInfoLoc = "/proc/stat"

	// For testing purposes:
	//
	// memInfoLoc = "./test_assets/meminfo"
	// cpuInfoLoc = "./test_assets/stat"
)

func getHardwareStats() HardwareStats {
	return HardwareStats{
		Cpu:    getCpuStats(),
		Memory: getMemoryStats(),
	}
}

func getMemoryStats() MemoryStats {
	ret := MemoryStats{}
	contents, err := ioutil.ReadFile(memInfoLoc)
	if err != nil {
		return ret
	}
	lines := strings.Split(string(contents), "\n")
	var total, idle float64
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			total, _ = strconv.ParseFloat(fields[1], 64)
		case "MemAvailable:":
			idle, _ = strconv.ParseFloat(fields[1], 64)
		}
	}
	ret.Total = getMemoryString(float64(total * 1000))
	ret.Idle = float32(100 * idle / total)
	ret.Used = 100 - ret.Idle

	return ret
}

func getCpuStats() CpuStats {
	idle, total, count := getCpuData()
	usedP := 100 - float32(100*idle/total)
	idleP := float32(100) - usedP

	return CpuStats{
		MaxThreads: int(count),
		Used:       usedP,
		Idle:       idleP,
	}
}

func getCpuData() (idle, total, count uint64) {
	contents, err := ioutil.ReadFile(cpuInfoLoc)
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val
				if i == 4 {
					idle = val
				}
			}
		} else if strings.HasPrefix(fields[0], "cpu") {
			count++
		}
	}
	return
}
