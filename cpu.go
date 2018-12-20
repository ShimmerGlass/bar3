package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func CPU(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() string {
		p, _ := cpu.Percent(0, false)
		usage := p[0]
		if usage >= 100 {
			usage = 99
		}
		return Comb(
			Icon("\uf0e4", ColorHighlight3),
			Style(
				fmt.Sprintf(" %2.0f%%", usage),
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
