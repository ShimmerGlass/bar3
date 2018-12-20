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
		err := w.CurrentByName("Paris, FR")
		if err != nil {
			log.Println(err)
			return ""
		}

		var icon string
		if len(w.Weather) > 0 {
			switch w.Weather[0].Icon {
			case "01d":
				icon = "\ue30d"
			case "02d":
				icon = "\ue302"
			case "03d", "03n":
				icon = "\ue33d"
			case "04d", "04n":
				icon = "\ue312"
			case "09d", "09n":
				icon = "\ue318"
			case "10d":
				icon = "\ue309"
			case "11d", "11n":
				icon = "\ue31d"
			case "13d":
				icon = "\ue30a"
			case "50d":
				icon = "\ue3ae"
			case "01n":
				icon = "\ue32b"
			case "02n":
				icon = "\ue32e"
			case "10n":
				icon = "\ue334"
			case "13n":
				icon = "\ue327"
			case "50n":
				icon = "\ue346"
			}
		}

		log.Printf("weather: %+v", w.Weather)

		return Comb(
			Icon(icon, ColorHighlight2),
			fmt.Sprintf("  %.1f°", w.Main.Temp),
		)
	})
}
