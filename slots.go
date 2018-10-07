package main

import (
	"log"
)

type SlotEntry struct {
	Slot
	LastStatus string
}

func Run(w *Writer, sep string, slots ...Slot) error {
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
				changed <- struct{}{}
			}
		}(i, entry)
	}

	for range changed {
		status := ""

		for _, s := range entries {
			if len(s.LastStatus) > 0 {
				if len(status) > 0 {
					status += sep
				}
				status += s.LastStatus
			}
		}

		err := w.Write(Block{
			FullText: status,
			Markup:   MarkupPango,
			Color:    colorFmt(ColorText),
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
