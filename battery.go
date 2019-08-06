package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/distatus/battery"
)

func Battery(interval time.Duration) Slot {
	return NewTimedSlot(interval, func() []Part {
		batteries, err := battery.GetAll()
		if err != nil {
			log.Println("battery: ", err)
			return nil
		}

		if len(batteries) == 0 {
			return nil
		}

		var icons = [][]string{
			{
				"\uf579",
				"\uf579",
				"\uf57a",
				"\uf57b",
				"\uf57c",
				"\uf57d",
				"\uf57e",
				"\uf57f",
				"\uf580",
				"\uf581",
			},
			{
				"\uf585",
				"\uf585",
				"\uf585",
				"\uf586",
				"\uf587",
				"\uf587",
				"\uf588",
				"\uf589",
				"\uf590",
				"\uf584",
			},
		}

		for _, b := range batteries {
			var si []string
			if b.State == battery.Charging {
				si = icons[1]
			} else {
				si = icons[0]
			}

			ratio := b.Current / b.Full
			if ratio < 0 {
				ratio = 0
			}
			if ratio > 1 {
				ratio = 1
			}

			return []Part{
				IconPart(si[int(math.Round(ratio*9))]),
				TextPart(" "),
				TextPart(fmt.Sprintf("%.0f%%", math.Round(ratio*100)), FontMono),
			}
		}

		return nil
	})
}
