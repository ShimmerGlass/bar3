package main

import (
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/mem"
)

func RAM(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() string {
		v, _ := mem.VirtualMemory()
		return Comb(
			Style("\uf2db", ColorHighlight3),
			"  ",
			humanize.Bytes(v.Free+v.Cached),
		)
	})
}
