package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path"
	"strings"
	"syscall"
	"time"

	"golang.org/x/oauth2"

	"github.com/patrickmn/go-cache"

	"github.com/skratchdot/open-golang/open"

	"github.com/godbus/dbus"
	"github.com/zmb3/spotify"
)

const spotifyRedirectURI = "http://localhost:8080/callback"

var spotifyAuthenticator = spotify.NewAuthenticator(spotifyRedirectURI,
	spotify.ScopeUserReadPrivate,
	spotify.ScopePlaylistReadPrivate,
	spotify.ScopeUserLibraryRead,
	spotify.ScopeUserLibraryModify,
)

func spotifyTrack(conn *dbus.Conn) (id spotify.ID, title string, artist string, err error) {
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

func spotifyPlayStatus(conn *dbus.Conn) (playing bool, err error) {
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

func spotifyClient() (*spotify.Client, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	tokenFilePath := path.Join(usr.HomeDir, ".spotify")

	client, err := spotifyClientSaved(tokenFilePath)
	if err == nil && client != nil {
		return client, nil
	}
	if err != nil {
		log.Println(err)
	}

	client, err = spotifyClientAcquire()
	if err != nil {
		return nil, err
	}

	tokenFile, err := os.OpenFile(tokenFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		log.Println(err)
		return client, nil
	}
	defer tokenFile.Close()

	token, err := client.Token()
	if err != nil {
		log.Println(err)
	} else {
		err = json.NewEncoder(tokenFile).Encode(token)
		if err != nil {
			log.Println(err)
			return client, nil
		}
	}

	return client, nil
}

func spotifyClientSaved(path string) (*spotify.Client, error) {
	tokenFile, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return nil, nil
	}

	token := &oauth2.Token{}
	err = json.NewDecoder(tokenFile).Decode(token)
	if err != nil {
		return nil, err
	}

	client := spotifyAuthenticator.NewClient(token)
	return &client, nil
}

func spotifyClientAcquire() (*spotify.Client, error) {
	ch := make(chan spotify.Client)
	state := "abc123"

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		tok, err := spotifyAuthenticator.Token(state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}

		// use the token to get an authenticated client
		client := spotifyAuthenticator.NewClient(tok)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "Login Completed!<script>window.close();</script>")

		ch <- client
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go http.ListenAndServe(":8080", nil)

	url := spotifyAuthenticator.AuthURL(state)
	time.Sleep(2 * time.Second)
	open.Run(url)

	// wait for auth to complete
	client := <-ch
	return &client, nil
}

func Music() Slot {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Println("Failed to connect to session bus:", err)
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
		client, err := spotifyClient()
		if err != nil {
			log.Println(err)
			return
		}

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
				id, _, _, err := spotifyTrack(conn)
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

	return NewTimedSlot(time.Minute, func() []Part {
		playing, err := spotifyPlayStatus(conn)
		if err != nil {
			log.Println(err)
			return nil
		}

		if !playing {
			return nil
		}

		id, title, artist, err := spotifyTrack(conn)
		if err != nil {
			log.Println(err)
			return nil
		}

		if title == "" && artist == "" {
			return nil
		}

		hasTrack := userHasTrack(id)

		var hasTrackIcon string
		if hasTrack {
			hasTrackIcon = "\uf004"
		} else {
			hasTrackIcon = "\uf08a"
		}

		parts := []Part{
			IconPart("\uf001"),
			TextPart("  "),
			TextPart(Elipsis(html.UnescapeString(title), 30)),
			TextPart("  "),
		}
		if len(artist) > 0 {
			parts = append(parts,
				IconPart("•"),
				TextPart(Elipsis(html.UnescapeString(artist), 30)),
				TextPart("  "),
			)
		}
		parts = append(parts,
			IconPart(hasTrackIcon),
		)

		return parts
	}, refresh)
}
