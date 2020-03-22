package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/process"
)

type ProcessWatcher struct {
	MaxCPU string
	MaxRAM string
}

func (w *ProcessWatcher) watch() {
	var lastProcessTimes map[int32]float64

	for range time.Tick(2 * time.Second) {
		pids, err := process.Pids()
		if err != nil {
			continue
		}

		cpuVals := map[int32]float64{}
		maxTime := float64(0)
		maxRAM := uint64(0)
		maxCPUName := ""
		maxRAMName := ""
		for _, pid := range pids {
			proc, err := process.NewProcess(pid)
			if err != nil {
				continue
			}

			procName, err := proc.Name()
			if err != nil {
				continue
			}

			// cpu
			times, err := proc.Times()
			if err != nil {
				continue
			}
			time := times.System + times.User
			cpuVals[pid] = time

			lastTime, ok := lastProcessTimes[pid]
			if !ok {
				continue
			}

			timeD := time - lastTime

			if timeD > maxTime {
				maxTime = timeD
				maxCPUName = procName
			}

			// ram
			ram, err := proc.MemoryInfo()
			if err != nil {
				continue
			}

			if ram.RSS > maxRAM {
				maxRAM = ram.RSS
				maxRAMName = procName
			}
		}

		lastProcessTimes = cpuVals
		w.MaxCPU = maxCPUName
		w.MaxRAM = maxRAMName
	}
}

func ProcessPart(name string) Part {
	part := TextPart(fmt.Sprintf(" %6s", Elipsis(name, 6)), FontMono)
	part.Lum = .7

	return part
}
