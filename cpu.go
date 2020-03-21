package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

func formatCPU(v float64) string {
	usage := int(math.Round(v))
	usageStr := strconv.Itoa(usage)
	switch len(usageStr) {
	case 3:
		usageStr = "00"
	case 1:
		usageStr = " " + usageStr
	}

	return fmt.Sprintf("%s%%", usageStr)
}

func maxCPUProcess(last map[int32]float64) (string, map[int32]float64) {
	pids, err := process.Pids()
	if err != nil {
		return "", nil
	}

	vals := map[int32]float64{}
	maxTime := float64(0)
	name := ""
	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		times, err := proc.Times()
		if err != nil {
			continue
		}
		time := times.System + times.User
		vals[pid] = time

		lastTime, ok := last[pid]
		if !ok {
			continue
		}

		timeD := time - lastTime

		if timeD > maxTime {
			maxTime = timeD
			name, _ = proc.Name()
		}
	}

	return name, vals
}

func CPU(interval time.Duration) Slot {
	var lastProcessTimes map[int32]float64

	return NewTimedSlot(interval, func() []Part {
		pg, _ := cpu.Percent(0, false)
		global := formatCPU(pg[0])

		pm, _ := cpu.Percent(0, true)
		maxCPU := 0.0
		for _, c := range pm {
			if c > maxCPU {
				maxCPU = c
			}
		}
		max := formatCPU(maxCPU)

		globalPart := TextPart(" "+global, FontMono)
		globalPart.Sat = pg[0] / 100

		maxPart := TextPart(" "+max, FontMono)
		maxPart.Sat = maxCPU / 100

		maxProcess, vals := maxCPUProcess(lastProcessTimes)
		lastProcessTimes = vals

		maxProcessPart := TextPart(fmt.Sprintf(" %8s", Elipsis(maxProcess, 8)), FontMono)
		maxProcessPart.Lum = .7

		return []Part{
			IconPart("\uf0e4"),
			globalPart,
			maxPart,
			maxProcessPart,
		}
	})
}
