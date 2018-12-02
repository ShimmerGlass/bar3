package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/disk"
)

func Storage(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() string {
		stat, err := disk.Usage("/")
		if err != nil {
			log.Println(err)
			return ""
		}
		return Comb(
			Icon("\uf0a0", ColorHighlight3),
			fmt.Sprintf("  %2.0f%%", stat.UsedPercent),
		)
	})
}
