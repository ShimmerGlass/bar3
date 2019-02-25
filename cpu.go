package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func CPU(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() string {
		p, _ := cpu.Percent(0, false)
		usage := int(math.Round(p[0]))
		usageStr := strconv.Itoa(usage)
		switch len(usageStr) {
		case 3:
			usageStr = "00"
		case 1:
			usageStr = " " + usageStr
		}
		return Comb(
			Icon("\uf0e4", ColorHighlight3),
			Style(
				fmt.Sprintf(" %s%%", usageStr),
				Grad(p[0],
					GradStop{0, ColorText},
					GradStop{50, ColorText},
					GradStop{75, ColorWarning},
					GradStop{100, ColorDanger},
				),
				FontMono,
			),
		)
	})
}
