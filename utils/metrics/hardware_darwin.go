package metrics

import (
	"os/exec"
	"strconv"
	"strings"
)

const (
	cpuPrefix    = "CPU usage:"
	memoryPrefix = "PhysMem:"
)

func getHardwareStats() HardwareStats {
	cmd := exec.Command("top", "-F", "-R", "-o", "cpu", "-l", "1", "-n0")
	result, err := cmd.Output()
	ret := HardwareStats{}
	if err != nil {
		return ret
	}

	lines := strings.Split(string(result), "\n")

	for _, line := range lines {
		// CPU Calc
		if strings.HasPrefix(line, cpuPrefix) {
			fields := strings.Split(strings.Replace(line, cpuPrefix, "", 1), ",")
			for _, cpuReport := range fields {
				cpuReport = strings.TrimSpace(cpuReport)
				if strings.HasSuffix(cpuReport, "idle") {
					split := strings.Split(cpuReport, "%")
					if result, err := strconv.ParseFloat(split[0], 32); err == nil {
						ret.Cpu.Idle = float32(result)
						ret.Cpu.Used = float32(100 - result)
					}
					break
				}
			}
			continue
		}

		// Memory Calc
		if strings.HasPrefix(line, memoryPrefix) {
			fields := strings.Split(strings.Replace(line, memoryPrefix, "", 1), ",")
			var total, used float64

			for _, memReport := range fields {
				split := strings.Fields(memReport)
				valueStr := strings.NewReplacer("K", "000", "M", "000000", "G", "000000000").Replace(split[0])
				value, _ := strconv.ParseFloat(valueStr, 64)
				total += value

				if strings.Contains(memReport, " used") {
					used = value
				}
				if strings.Contains(memReport, "wired") {
					split := strings.Fields(strings.Split(memReport, "(")[1])
					valueStr := strings.NewReplacer("K", "000", "M", "000000", "G", "000000000").Replace(split[0])
					value, _ := strconv.ParseFloat(valueStr, 64)
					used -= value
				}
			}

			ret.Memory.Total = getMemoryString(total)
			ret.Memory.Used = float32(100 * used / total)
			ret.Memory.Idle = 100 - ret.Memory.Used

			continue
		}
	}

	return ret
}

func getCpuStats() CpuStats {
	return getHardwareStats().Cpu
}

func getMemoryStats() MemoryStats {
	return getHardwareStats().Memory
}
