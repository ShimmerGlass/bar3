package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
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

func CPU(interval time.Duration) Slot {
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

		return []Part{
			IconPart("\uf0e4"),
			globalPart,
			maxPart,
		}
	})
}
