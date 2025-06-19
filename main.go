package main

import (
	"fmt"
	"go-last/config"
	"go-last/imagegen"
	"go-last/rank"
	"go-last/resources"
	"go-last/utils"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
)

func getArgs(w http.ResponseWriter, r *http.Request) (int, bool, bool, []rank.Album, error) {
	user := r.URL.Query().Get("user")
	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		return 0, false, false, nil, fmt.Errorf("user parameter is required")
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = config.Default.GridSize
	} else if size < config.MinGrid {
		size = config.MinGrid
	} else if size > config.MaxGrid {
		size = config.MaxGrid
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = config.Default.Period
	}

	captionArg := r.URL.Query().Get("captions")
	caption := true
	if captionArg != "on" {
		caption = false
	}

	playsArg := r.URL.Query().Get("plays")
	plays := true
	if playsArg != "on" {
		plays = false
	}

	albums, err := rank.PlayCount(config.Client, user, period, size*size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 0, false, false, nil, fmt.Errorf("grid gen failed")
	}

	return size, caption, plays, albums, nil
}

func collage(w http.ResponseWriter, r *http.Request) {
	size, caption, plays, albums, err := getArgs(w, r)
	if err != nil {
		return
	}

	img, err := imagegen.Grid(albums, size, caption, plays)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg;")
	jpeg.Encode(w, img, &jpeg.Options{Quality: 95})
	img = nil
	return
}

func altText(w http.ResponseWriter, r *http.Request) {
	size, _, plays, albums, err := getArgs(w, r)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain;")
	fmt.Fprintf(w, "%s", utils.GenAltText(albums, size, plays))
	return
}

func main() {
	config.Init()
	resources.Init()
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServerFS(resources.Static))
	mux.HandleFunc("/collage", collage)
	mux.HandleFunc("/altText", altText)

	err := http.ListenAndServe(":"+config.Port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
