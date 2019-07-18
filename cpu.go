package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func CPU(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() []Part {
		p, _ := cpu.Percent(0, false)
		usage := int(math.Round(p[0]))
		usageStr := strconv.Itoa(usage)
		switch len(usageStr) {
		case 3:
			usageStr = "00"
		case 1:
			usageStr = " " + usageStr
		}

		txt := TextPart(fmt.Sprintf(" %s%%", usageStr), FontMono)
		txt.Sat = p[0] / 100

		RainbowPanSpeed = p[0] / 2

		return []Part{
			IconPart("\uf0e4"),
			txt,
		}
	})
}
