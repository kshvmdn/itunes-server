package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/gorilla/mux"
)

type Track struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

type Status struct {
	Status  string `json:"status"`
	Current Track  `json:"current"`
}

type Route struct {
	Endpoint string
	Command  string
}

func ReaderToString(out io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	b := buf.Bytes()
	s := *(*string)(unsafe.Pointer(&b))
	return strings.Trim(s, "\n")
}

func ExecItunes(cmd string) string {
	fullCmd := fmt.Sprintf("tell Application \"iTunes\" %s", cmd)
	c := exec.Command("/usr/bin/osascript", "-e", fullCmd)

	out, _ := c.StdoutPipe()

	defer c.Wait()

	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	return ReaderToString(out)
}

func Index(w http.ResponseWriter, r *http.Request) {
	status := ExecItunes("to player state as string")

	var name string
	var artist string
	var album string

	if status == "playing" {
		name = ExecItunes("to name of current track as string")
		artist = ExecItunes("to artist of current track as string")
		album = ExecItunes("to album of current track as string")
	}

	json.NewEncoder(w).Encode(Status{status, Track{name, artist, album}})
}

func ListTracks(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var tracks []Track

	start, err := strconv.Atoi(r.FormValue("skip"))
	if err != nil {
		start = 0
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil || limit > 100 {
		limit = 100
	}

	trackIds := ExecItunes("to get the id of (every track)")

	end := start + limit
	if end > len(trackIds) {
		end = len(trackIds) - 1
	}

	wg.Add(end - start)

	for _, id := range strings.Split(trackIds, ", ")[start:end] {
		go func(id string) {
			defer wg.Done()
			name := ExecItunes(fmt.Sprintf("to get the name of track id %s", id))
			artist := ExecItunes(fmt.Sprintf("to get the artist of track id %s", id))
			album := ExecItunes(fmt.Sprintf("to get the album of track id %s", id))
			tracks = append(tracks, Track{name, artist, album})
		}(id)
	}

	wg.Wait()
	json.NewEncoder(w).Encode(tracks)
}

func PlaySong(w http.ResponseWriter, r *http.Request) {
	songName := mux.Vars(r)["song_name"]
	cmd := fmt.Sprintf("to play (every track of playlist \"Library\" whose name is \"%s\")", songName)
	ExecItunes(cmd)
	http.Redirect(w, r, "/", 302)
}

func PlayArtist(w http.ResponseWriter, r *http.Request) {
	artistName := mux.Vars(r)["artist_name"]
	cmd := fmt.Sprintf("to play (every track of playlist \"Library\" whose artist is \"%s\")", artistName)
	ExecItunes(cmd)
	http.Redirect(w, r, "/", 302)
}

func PlayAlbum(w http.ResponseWriter, r *http.Request) {
	albumName := mux.Vars(r)["album_name"]
	cmd := fmt.Sprintf("to play (every track of playlist \"Library\" whose album is \"%s\")", albumName)
	ExecItunes(cmd)
	http.Redirect(w, r, "/", 302)
}

func PlayRandom(w http.ResponseWriter, r *http.Request) {
	trackIds := strings.Split(ExecItunes("to get the id of (every track)"), ", ")
	cmd := fmt.Sprintf("to play (every track of playlist \"Library\" whose id is \"%s\")", trackIds[rand.Intn(len(trackIds))])
	ExecItunes(cmd)
	http.Redirect(w, r, "/", 302)
}

func Routes() []Route {
	return []Route{
		{"open", "to open"},
		{"exit", "to quit"},
		{"play", "to play"},
		{"pause", "to pause"},
		{"stop", "to stop"},
		{"next", "to next track"},
		{"prev", "to previous track"},
		{"mute", "to set mute to true"},
		{"unmute", "to set mute to false"},
	}
}

func main() {
	portPtr := flag.String("port", "8080", "Port to run the server on.")
	disableTrackListPtr := flag.Bool("no-track-list", false, "Disable the track list endpoint.")
	flag.Parse()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)

	if !*disableTrackListPtr {
		router.HandleFunc("/tracks", ListTracks)
	}

	router.HandleFunc("/play/track/{track_name}", PlaySong)
	router.HandleFunc("/play/artist/{artist_name}", PlayArtist)
	router.HandleFunc("/play/album/{album_name}", PlayAlbum)
	router.HandleFunc("/shuffle", PlayRandom)

	for _, route := range Routes() {
		func(endpoint string, cmd string) {
			router.HandleFunc(fmt.Sprintf("/%s", endpoint), func(w http.ResponseWriter, r *http.Request) {
				ExecItunes(cmd)
				http.Redirect(w, r, "/", 302)
			})
		}(route.Endpoint, route.Command)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *portPtr), router))
}
