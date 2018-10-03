package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"strings"
	"time"

	"github.com/godbus/dbus"
)

func dBusTrack(conn *dbus.Conn) (title string, artist string, err error) {
	res := dbus.Variant{}

	err = conn.Object(
		"org.mpris.MediaPlayer2.spotify",
		"/org/mpris/MediaPlayer2",
	).Call(
		"org.freedesktop.DBus.Properties.Get",
		0,
		"org.mpris.MediaPlayer2.Player",
		"Metadata",
	).Store(&res)
	if err != nil {
		return
	}

	val, ok := res.Value().(map[string]dbus.Variant)
	if !ok {
		return
	}

	artistVar, ok := val["xesam:artist"].Value().([]string)
	if ok {
		artist = strings.Join(artistVar, ", ")
	}

	title, _ = val["xesam:title"].Value().(string)

	return
}

func dBusPlaying(conn *dbus.Conn) (playing bool, err error) {
	res := dbus.Variant{}

	err = conn.Object(
		"org.mpris.MediaPlayer2.spotify",
		"/org/mpris/MediaPlayer2",
	).Call(
		"org.freedesktop.DBus.Properties.Get",
		0,
		"org.mpris.MediaPlayer2.Player",
		"PlaybackStatus",
	).Store(&res)
	if err != nil {
		return
	}

	return res.Value().(string) == "Playing", nil
}

func Music() Slot {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}

	refresh := make(chan struct{})

	go func() {
		for {
			err = conn.BusObject().Call(
				"org.freedesktop.DBus.AddMatch",
				0,
				"type='signal',path='/org/mpris/MediaPlayer2',interface='org.freedesktop.DBus.Properties',sender='org.mpris.MediaPlayer2.spotify'",
			).Err
			if err != nil {
				log.Println(err)
			}
			c := make(chan *dbus.Signal, 10)
			conn.Signal(c)

			for range c {
				refresh <- struct{}{}
			}

			time.Sleep(time.Second)
		}
	}()

	return NewTimedSlot(time.Second, func() string {
		playing, err := dBusPlaying(conn)
		if err != nil {
			log.Println(err)
			return ""
		}

		if !playing {
			return ""
		}

		title, artist, err := dBusTrack(conn)
		if err != nil {
			log.Println(err)
			return ""
		}

		if title == "" && artist == "" {
			return ""
		}

		return fmt.Sprintf("%s  %s %s %s",
			iconC("\uf001", ColorHighlight),
			html.EscapeString(title),
			iconC("•", ColorInactive),
			html.EscapeString(artist),
		)
	}, refresh)
}
