package main

import (
	"time"
)

func Date() Slot {
	return NewTimedSlot(time.Minute, func() string {
		return Comb(
			Style("\uf073  ", ColorHighlight2),
			time.Now().Format("Mon, _2 Jan 15:04"),
		)
	})
}
