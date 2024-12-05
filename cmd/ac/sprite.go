package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	// "github.com/cloudkucooland/AlarmClock/resources"
	spriteres "github.com/cloudkucooland/AlarmClock/resources/sprites"
)

type sprite struct {
	name  string
	raw   []byte
	image *ebiten.Image
	alpha *image.Alpha // alpha chanel to quickly determine "in"
	x     int          // current screen location
	y     int          // current screen location
	do    func(*sprite)
}

// probably remove all this biz. use the control.in()
func (s *sprite) in(x, y int) bool {
	return s.alpha.At(x-s.x, y-s.y).(color.Alpha).A > 0
}

// if has a location other than 0,0, it must be on sceen
func (s *sprite) onscreen() bool {
	if s.x == 0 && s.y == 0 {
		return false
	}
	return true
}

// should this be a map key'd off name?
var sprites = []sprite{
	{
		name: "Artist",
		raw:  spriteres.ArtistPNG,
	},
	{
		name: "Baby",
		raw:  spriteres.BabyPNG,
	},
	{
		name: "Bathtime",
		raw:  spriteres.BathtimePNG,
	},
	{
		name: "Confused",
		raw:  spriteres.ConfusedPNG,
	},
	{
		name: "Happy",
		raw:  spriteres.HappyPNG,
	},
	{
		name: "Love",
		raw:  spriteres.LovePNG,
	},
}

func loadSprites() error {
	for x := range sprites {
		img, _, err := image.Decode(bytes.NewReader(sprites[x].raw))
		if err != nil {
			return err
		}
		sprites[x].image = ebiten.NewImageFromImage(img)
		b := img.Bounds()
		sprites[x].alpha = image.NewAlpha(b)
		for j := b.Min.Y; j < b.Max.Y; j++ {
			for i := b.Min.X; i < b.Max.X; i++ {
				sprites[x].alpha.Set(i, j, img.At(i, j))
			}
		}
		if sprites[x].do == nil {
			sprites[x].do = chirp
		}
	}
	return nil
}

func chirp(s *sprite) {
	fmt.Println("play chirp", s.name)
}
