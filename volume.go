package main

import (
	"log"
	"math"
	"strings"
	"time"

	"github.com/godbus/dbus"
	"github.com/sqp/pulseaudio"
)

type pulselistener struct {
	out chan struct{}
}

func (two *pulselistener) DeviceVolumeUpdated(path dbus.ObjectPath, values []uint32) {
	two.out <- struct{}{}
}

func (two *pulselistener) DeviceMuteUpdated(path dbus.ObjectPath, state bool) {
	two.out <- struct{}{}
}

func Volume() Slot {
	pulse, err := pulseaudio.New()
	if err != nil {
		log.Fatal(err)
	}

	app := &pulselistener{make(chan struct{})}

	errs := pulse.Register(app)
	if len(errs) > 0 {
		log.Fatal(err)
	}

	go pulse.Listen()

	return NewTimedSlot(time.Minute, func() string {
		sinks, err := pulse.Core().ListPath("Sinks")
		if err != nil {
			log.Fatal(err)
		}

		var vol float64

		if len(sinks) >= 0 {

			var muted bool
			err = pulse.Device(sinks[0]).Get("Mute", &muted)
			if err != nil {
				log.Fatal(err)
			}
			if !muted {
				var volumes []uint32
				err = pulse.Device(sinks[0]).Get("Volume", &volumes)
				if err != nil {
					log.Fatal(err)
				}

				var volumeSteps uint32
				err = pulse.Device(sinks[0]).Get("VolumeSteps", &volumeSteps)
				if err != nil {
					log.Fatal(err)
				}

				var volTotal uint32
				for _, v := range volumes {
					volTotal += v
				}

				vol = float64(volTotal) / float64(len(volumes)) / float64(volumeSteps)
			}
		}

		var pattern string
		if vol == 0 {
			pattern = iconC("\uf026", ColorInactive)
		} else {
			pattern = iconC("\uf028", ColorHighlight)
		}

		barSize := float64(10)
		pattern += color(ColorInactive, " [ ")
		pattern += color(ColorHighlight, strings.Repeat("/", int(math.Round(vol*barSize))))
		pattern += color(ColorInactive, strings.Repeat("/", int(math.Round(barSize-vol*barSize))))
		pattern += color(ColorInactive, " ]")

		return bold(pattern)
	}, app.out)
}