package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/skratchdot/open-golang/open"

	"github.com/godbus/dbus"
	"github.com/zmb3/spotify"
)

func dBusTrack(conn *dbus.Conn) (id spotify.ID, title string, artist string, err error) {
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

	rid, _ := val["mpris:trackid"].Value().(string)
	if len(rid) > 0 {
		parts := strings.Split(rid, ":")
		id = spotify.ID(parts[len(parts)-1])
	}

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

	userTracks := cache.New(10*time.Minute, time.Minute)

	userHasTrack := func(id spotify.ID) bool {
		return false
	}

	go func() {
		const redirectURI = "http://localhost:8080/callback"

		var (
			auth = spotify.NewAuthenticator(redirectURI,
				spotify.ScopeUserReadPrivate,
				spotify.ScopePlaylistReadPrivate,
				spotify.ScopeUserLibraryRead,
				spotify.ScopeUserLibraryModify,
			)
			ch    = make(chan *spotify.Client)
			state = "abc123"
		)
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			tok, err := auth.Token(state, r)
			if err != nil {
				http.Error(w, "Couldn't get token", http.StatusForbidden)
				log.Fatal(err)
			}
			if st := r.FormValue("state"); st != state {
				http.NotFound(w, r)
				log.Fatalf("State mismatch: %s != %s\n", st, state)
			}
			// use the token to get an authenticated client
			client := auth.NewClient(tok)
			fmt.Fprintf(w, "Login Completed!<script>window.close();</script>")
			ch <- &client
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Got request for:", r.URL.String())
		})
		go http.ListenAndServe(":8080", nil)

		url := auth.AuthURL(state)
		open.Run(url)

		// wait for auth to complete
		client := <-ch

		userHasTrack = func(id spotify.ID) bool {
			v, ok := userTracks.Get(string(id))
			if ok {
				return v.(bool)
			}

			go func() {
				res, err := client.UserHasTracks(id)
				if err != nil {
					log.Println(err)
					return
				}

				if len(res) == 0 {
					return
				}

				userTracks.SetDefault(string(id), res[0])
				refresh <- struct{}{}
			}()

			return false
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1)

		go func() {
			for range sigs {
				id, _, _, err := dBusTrack(conn)
				if err != nil {
					log.Println(err)
					continue
				}

				res, err := client.UserHasTracks(id)
				if err != nil {
					log.Println(err)
					continue
				}

				if len(res) == 0 {
					continue
				}

				if !res[0] {
					err = client.AddTracksToLibrary(id)
					if err != nil {
						log.Println(err)
						continue
					}
				} else {
					err = client.RemoveTracksFromLibrary(id)
					if err != nil {
						log.Println(err)
						continue
					}
				}

				userTracks.Delete(string(id))
				refresh <- struct{}{}
			}
		}()
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

		id, title, artist, err := dBusTrack(conn)
		if err != nil {
			log.Println(err)
			return ""
		}

		if title == "" && artist == "" {
			return ""
		}

		hasTrack := userHasTrack(id)

		hasTrackLabel := ""
		if hasTrack {
			hasTrackLabel = iconC(" \uf004", "#ef0946")
		}

		return fmt.Sprintf("%s  %s %s %s %s",
			iconC("\uf001", ColorHighlight),
			html.EscapeString(elipsis(title, 30)),
			iconC("•", ColorInactive),
			html.EscapeString(elipsis(artist, 30)),
			hasTrackLabel,
		)
	}, refresh)
}
