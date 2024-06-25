package imagegen

import (
	"go-last/config"
	"go-last/rank"
	"go-last/resources"
	"go-last/utils"
	"image"
	"image/draw"
	"strconv"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
	_ "golang.org/x/image/webp"
)

// TODO: Make this cleaner
func drawText(img *image.NRGBA, txt []string, rect image.Rectangle, x int, y int) {
	c := freetype.NewContext()
	c.SetDPI(50)
	c.SetFont(resources.Font)
	c.SetFontSize(26)
	c.SetClip(rect)
	c.SetDst(img)
	c.SetHinting(font.HintingVertical)
	for i, str := range txt {
		c.SetSrc(image.NewUniform(config.Font.AltColor))
		c.DrawString(str, freetype.Pt(x+1, y+((i+1)*16)+1))
		c.SetSrc(image.NewUniform(config.Font.Color))
		c.DrawString(str, freetype.Pt(x+2, y+((i+1)*16)))
	}
}

func drawArt(
	album rank.Album, final *image.NRGBA,
	caption bool, plays bool,
	x int, y int) {
	rect := image.Rect(x, y, x+config.ArtSize, y+config.ArtSize)
	draw.Draw(final, rect, album.Art.Img, image.Point{0, 0}, draw.Src)
	if caption || !album.Art.Present || plays {
		write := []string{}
		if caption || !album.Art.Present {
			write = append(write, album.Name, album.Artist)
		}
		if plays {
			write = append(write, "Plays: "+strconv.Itoa(album.Plays))
		}
		drawText(
			final,
			write,
			rect, x, y)
	}
}

func Grid(albums []rank.Album, gridSize int, caption bool, plays bool) (*image.NRGBA, error) {
	imagesize := config.ArtSize * gridSize
	fin_img := image.NewNRGBA(image.Rect(0, 0, imagesize, imagesize))

	albumStream := make(chan rank.Album)
	go utils.DownloadImages(albums, albumStream)

	for album := range albumStream {
		x := (album.Rank % gridSize) * config.ArtSize
		y := (album.Rank / gridSize) * config.ArtSize
		drawArt(album, fin_img, caption, plays, x, y)
	}

	return fin_img, nil
}
