package main

import (
	"go-last/config"
	"go-last/imagegen"
	"go-last/rank"
	"go-last/resources"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
)

func collage(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
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

func main() {
	config.Init()
	resources.Init()
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServerFS(resources.Static))
	mux.HandleFunc("/collage", collage)

	err := http.ListenAndServe(":"+config.Port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
