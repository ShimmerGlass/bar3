package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/disk"
)

func Storage(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() []Part {
		stat, err := disk.Usage("/")
		if err != nil {
			log.Println(err)
			return nil
		}
		return []Part{
			IconPart("\uf0a0"),
			TextPart(fmt.Sprintf(" %2.0f%%", stat.UsedPercent), FontMono),
		}
	})
}
