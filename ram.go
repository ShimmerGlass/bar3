package main

import (
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/mem"
)

func RAM(interval time.Duration, pw *ProcessWatcher) Slot {
	return NewTimedSlot(interval, func() []Part {
		v, _ := mem.VirtualMemory()
		return []Part{
			IconPart("\uf2db"),
			TextPart("  "),
			TextPart(humanize.Bytes(v.Free + v.Cached)),
			TextPart(" "),
			ProcessPart(pw.MaxRAM),
		}
	})
}
