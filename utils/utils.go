package utils

import (
	"fmt"
	"go-last/rank"
	"go-last/resources"
	"image"
	"log"
	"net/http"
	"sync"
)

func GenAltText(albums []rank.Album, size int, plays bool) string {
	var altText string
	var maxDigits int

	len := len(albums)
	for len > 0 {
		len /= 10
		maxDigits++
	}

	altText += fmt.Sprintf("%dx%d album grid:", size, size)

	for i, album := range albums {
		altText += "\n" + fmt.Sprintf("%0*d", maxDigits, i + 1)+ ". " + album.Name + " - " + album.Artist;
		if plays {
			altText += fmt.Sprintf(" (%d plays)", album.Plays)
		}
	}
	return altText
}


func DownloadImages(albums []rank.Album, albumStream chan<- rank.Album) {
	var wg sync.WaitGroup
	wg.Add(len(albums))
	for i, album := range albums {
		go func(i int, album rank.Album) {
			download, err := http.Get(album.Art.URL)
			if err != nil || download.StatusCode != http.StatusOK {
				album.Art.Img = resources.Placeholder
			} else {
				defer download.Body.Close()
				album.Art.Img, _, err = image.Decode(download.Body)
				if err != nil {
					log.Fatal(err)
				}
				album.Art.Present = true
			}

			albumStream <- album
			wg.Done()
		}(i, album)
	}
	wg.Wait()
	close(albumStream)
}
