package utils

import (
	"go-last/rank"
	"go-last/resources"
	"image"
	"log"
	"net/http"
	"sync"
)

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
