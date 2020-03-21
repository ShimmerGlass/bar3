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

func DiskUtilisation(name string, interval time.Duration) Slot {
	var last *disk.IOCountersStat
	return NewTimedSlot(interval, func() []Part {
		allStats, err := disk.IOCounters(name)
		if err != nil {
			log.Println(err)
			return nil
		}
		stats, ok := allStats[name]
		if !ok {
			return nil
		}

		var readPs, writePs float64
		if last == nil {
			last = &stats
		} else {
			readD := stats.ReadBytes - last.ReadBytes
			writeD := stats.WriteBytes - last.WriteBytes
			last = &stats

			readPs = float64(readD) / interval.Seconds()
			writePs = float64(writeD) / interval.Seconds()
		}

		return []Part{
			TextPart(fmt.Sprintf(" %-7s", humanize.Bytes(uint64(readPs))), FontMono),
			TextPart(" "),
			TextPart(fmt.Sprintf(" %7s", humanize.Bytes(uint64(writePs))), FontMono),
		}
	})
}
