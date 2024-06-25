package resources

import (
	"bytes"
	"embed"
	"image"
	"io/fs"
	"log"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

//go:embed static/*
var rawStatic embed.FS

//go:embed placeholder.png
var rawPlaceholder []byte

//go:embed unifont-15.1.05.ttf
var rawFont []byte

var Font *truetype.Font
var Placeholder image.Image
var Static fs.FS

func Init() {
	var err error
	Font, err = freetype.ParseFont(rawFont)
	if err != nil {
		log.Fatal(err)
	}
	Placeholder, _, err = image.Decode(bytes.NewReader(rawPlaceholder))
	if err != nil {
		log.Fatal(err)
	}

	Static, err = fs.Sub(rawStatic, "static")
	if err != nil {
		log.Fatal(err)
	}
}
