package main

import (
	"log"
	"time"

	"github.com/arduino/go-apt-client"
)

func APT(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() (parts []Part) {
		pkgs, err := apt.ListUpgradable()
		if err != nil {
			log.Println(err)
			return
		}

		if len(pkgs) > 0 {
			parts = append(parts, IconPart("\uf019"))
		}

		return parts
	})
}
