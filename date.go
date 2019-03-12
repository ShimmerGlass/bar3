package main

import (
	"time"
)

func Date() Slot {
	s := make(BasicSlot)

	go func() {
		for {
			s <- []Part{
				IconPart("\uf073"),
				TextPart(time.Now().Format("  Mon, _2 Jan 15:04")),
			}
			now := time.Now()
			time.Sleep(time.Minute - time.Second*time.Duration(now.Second()) + time.Second - time.Nanosecond*time.Duration(now.Nanosecond()))
		}
	}()
	return s
}
