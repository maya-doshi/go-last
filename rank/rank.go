package rank

import (
	"image"
	"strconv"
	"time"

	lfm "github.com/pazuzu156/lastfm-go"
)

type AlbumArt struct {
	Img     image.Image
	URL     string
	Present bool
}

type Album struct {
	Name   string
	Artist string
	Plays  int
	Rank   int
	Length time.Duration
	Art    AlbumArt
}

func PlayCount(cli *lfm.API, user string, period string, count int) ([]Album, error) {
	var top []Album
	LFMTop, err := cli.User.GetTopAlbums(map[string]interface{}{
		"user":   user,
		"period": period,
		"limit":  strconv.Itoa(count),
	})

	if err != nil {
		return top, err
	}

	for i, album := range LFMTop.Albums {
		if i >= count {
			break
		}
		albumImageIndex := len(album.Images) - 1
		plays, err := strconv.Atoi(album.PlayCount)
		if err != nil {
			plays = 0
		}
		new := Album{
			Name:   album.Name,
			Artist: album.Artist.Name,
			Rank:   i,
			Plays:  plays,
			Art: AlbumArt{
				URL:     album.Images[albumImageIndex].URL,
				Present: false,
			},
		}
		top = append(top, new)
	}

	return top, nil
}

func Time() ([]Album, error) {
	// TODO:
	return nil, nil
}
