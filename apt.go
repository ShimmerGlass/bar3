package main

import (
	"fmt"
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
			return Comb(
				Icon("\uf49e", ColorHighlight3),
				" ",
				fmt.Sprint(len(pkgs)),
			)
		}

		return ""
	})
}
