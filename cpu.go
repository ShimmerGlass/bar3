package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func CPU(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() string {
		p, _ := cpu.Percent(0, false)
		return fmt.Sprintf(`%s  <span font_family="Noto Mono">%2.0f%%</span>`, iconC("\uf0e4", ColorHighlight3), p[0])
	})
}
