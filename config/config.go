package config

import (
	"image/color"
	"os"

	lfm "github.com/pazuzu156/lastfm-go"
)

var Default struct {
	GridSize int
	Period   string
}

var Font struct {
	Color    color.RGBA
	AltColor color.RGBA
}

var Client *lfm.API

var MinGrid int
var MaxGrid int

const ArtSize = 300

const Port = "4391"

func Init() {

	MinGrid = 2
	MaxGrid = 20

	Default.GridSize = 3
	Default.Period = "7day"

	Client = lfm.New(
		os.Getenv("LFM_KEY"),
		os.Getenv("LFM_SECRET"),
	)

	Font.Color = color.RGBA{255, 255, 255, 255}
	Font.AltColor = color.RGBA{0, 0, 0, 255}
}
