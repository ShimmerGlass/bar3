package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func CPU(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() string {
		p, _ := cpu.Percent(0, false)
		return Comb(
			Icon("\uf0e4", ColorHighlight3),
			Style(
				fmt.Sprintf(" %2.0f%%", p[0]),
				Grad(p[0],
					GradStop{0, ColorText},
					GradStop{25, ColorText},
					GradStop{50, ColorWarning},
					GradStop{100, ColorDanger},
				),
				FontMono,
			),
		)
	})
}
