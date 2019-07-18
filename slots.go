package main

import (
	"log"
	"time"
)

var RainbowPanSpeed = 1.0

type SlotEntry struct {
	Slot
	LastStatus []Part
}

func Run(w *Writer, sep Part, slots ...Slot) error {
	entries := make([]*SlotEntry, len(slots))
	changed := make(chan struct{})

	for i, s := range slots {
		entry := &SlotEntry{
			Slot: s,
		}
		entries[i] = entry

		go func(s *SlotEntry) {
			s.Signal()
		}(entry)

		go func(i int, s *SlotEntry) {
			for status := range s.Out() {
				s.LastStatus = status
				select {
				case changed <- struct{}{}:
				default:
				}
			}
		}(i, entry)
	}

	timer := time.NewTicker(100 * time.Millisecond)

	lastHue := 0
	lastHueTime := time.Now()
	for range changed {
		<-timer.C
		status := []Part{}

		for _, s := range entries {
			if len(s.LastStatus) > 0 {
				if len(status) > 0 {
					status = append(status, sep)
				}
				status = append(status, s.LastStatus...)
			}
		}

		newHue := (lastHue + int(time.Since(lastHueTime).Seconds()*RainbowPanSpeed)) % 360
		lastHue = newHue
		lastHueTime = time.Now()
		err := w.Write(Block{
			FullText: Render(float64(newHue), float64(newHue), status),
			Markup:   MarkupPango,
			Color:    ColorText.String(),
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
