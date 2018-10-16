package main

import (
	"log"
	"time"

	"github.com/arduino/go-apt-client"
)

func APT(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() string {
		pkgs, err := apt.ListUpgradable()
		if err != nil {
			log.Println(err)
			return ""
		}

		if len(pkgs) > 0 {
			return Icon("\uf019", ColorHighlight2)
		}

		return ""
	})
}
