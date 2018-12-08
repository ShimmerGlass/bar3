package main

import (
	"fmt"
	"log"
	"os"
	"time"

	owm "github.com/briandowns/openweathermap"
)

func Weather(interval time.Duration) Slot {
	w, err := owm.NewCurrent("C", "en", os.Getenv("OWM_API_KEY")) // fahrenheit (imperial) with Russian output
	if err != nil {
		log.Println(err)
	}

	return NewTimedSlot(interval, func() string {
		if w == nil {
			return ""
		}
		err := w.CurrentByName("Paris,FR")
		if err != nil {
			log.Println(err)
		}

		var thunder, rain, snow, fog, clear, clouds bool
		for _, c := range w.Weather {
			switch c.ID / 100 {
			case 2:
				thunder = true
			case 3, 5:
				rain = true
			case 6:
				snow = true
			case 7:
				fog = true
			case 8:
				clouds = true
			}
			if c.ID == 800 {
				clear = true
			}
		}

		var icon string
		switch {
		case thunder && rain:
			icon = "\ufb7c"
		case thunder:
			icon = "\ufa92"
		case snow && rain:
			icon = "\ufb7d"
		case snow:
			icon = "\ufa97"
		case rain:
			icon = "\ufa95"
		case fog:
			icon = "\ufa90"
		case clouds:
			icon = "\ufa8f"
		case clear:
			icon = "\ufa98"
		}

		return Comb(
			Icon(icon, ColorHighlight2),
			fmt.Sprintf(" %.1f°", w.Main.Temp),
		)
	})
}
