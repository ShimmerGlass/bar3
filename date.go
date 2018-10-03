package main

import (
	"fmt"
	"time"
)

func Date() Slot {
	return NewTimedSlot(time.Minute, func() string {
		return fmt.Sprintf("%s  %s", iconC("\uf073", ColorHighlight2), time.Now().Format("Mon, _2 Jan 15:04"))
	})
}
