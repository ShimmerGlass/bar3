package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func CPU(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() string {
		p, _ := cpu.Percent(0, false)
		return fmt.Sprintf(
			`%s  <span font_family="Noto Mono">%s</span>`,
			iconC("\uf0e4", ColorHighlight3),
			colorize(
				colorSlide(p[0], ColorSlideStop{0, ColorText}, ColorSlideStop{50, ColorWarning}, ColorSlideStop{100, ColorDanger}),
				fmt.Sprintf("%2.0f%%", p[0]),
			),
		)
	})
}
