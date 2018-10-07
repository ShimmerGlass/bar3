package main

import (
	"time"
)

func Date() Slot {
	s := make(BasicSlot)

	go func() {
		for {
			s <- Comb(
				Icon("\uf073  ", ColorHighlight2),
				time.Now().Format("Mon, _2 Jan 15:04"),
			)
			now := time.Now()
			time.Sleep(time.Minute - time.Second*time.Duration(now.Second()) + time.Second - time.Nanosecond*time.Duration(now.Nanosecond()))
		}
	}()
	return s
}
