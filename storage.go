package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dustin/go-humanize"
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

func DiskUtilisation(names []string, interval time.Duration) Slot {
	last := map[string]disk.IOCountersStat{}
	return NewTimedSlot(interval, func() []Part {
		allStats, err := disk.IOCounters(names...)
		if err != nil {
			log.Println(err)
			return nil
		}

		var readPs, writePs float64
		for name, stats := range allStats {
			lastStats, ok := last[name]
			if !ok {
				last[name] = stats
			} else {
				readD := stats.ReadBytes - lastStats.ReadBytes
				writeD := stats.WriteBytes - lastStats.WriteBytes
				last[name] = stats

				readPs += float64(readD) / interval.Seconds()
				writePs += float64(writeD) / interval.Seconds()
			}
		}

		return []Part{
			TextPart(fmt.Sprintf("%-7s", humanize.Bytes(uint64(readPs))), FontMono),
			TextPart(" "),
			TextPart(fmt.Sprintf("%7s", humanize.Bytes(uint64(writePs))), FontMono),
		}
	})
}
