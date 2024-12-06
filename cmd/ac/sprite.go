package main

import (
	"bytes"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	spriteres "github.com/cloudkucooland/AlarmClock/resources/sprites"
)

type sprite struct {
	name  string
	raw   []byte
	image *ebiten.Image
	x     int          // current screen location
	y     int          // current screen location
	scale float32
	do    func(*sprite)
	ani   *spriteanimation
}

func (s *sprite) in(x, y int) bool {
	return false
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
	{
		name: "Indignent",
		raw:  spriteres.IndignentPNG,
	},
	{
		name: "Love",
		raw:  spriteres.LovePNG,
	},
	{
		name: "Mad",
		raw:  spriteres.MadPNG,
	},
	{
		name: "Pinwheel",
		raw:  spriteres.PinwheelPNG,
	},
	{
		name: "Spring",
		raw:  spriteres.SpringPNG,
	},
	{
		name: "Swan Mommy",
		raw:  spriteres.SwanmommyPNG,
	},
	{
		name: "Tea Time",
		raw:  spriteres.TeatimePNG,
	},
}

func getSprite(name string) *sprite {
	for x := range sprites {
		if sprites[x].name == name {
			return &sprites[x]
		}
	}
	return &sprites[0]
}

func loadSprites() error {
	for x := range sprites {
		img, _, err := image.Decode(bytes.NewReader(sprites[x].raw))
		if err != nil {
			return err
		}
		sprites[x].image = ebiten.NewImageFromImage(img)
		if sprites[x].do == nil {
			sprites[x].do = chirp
		}
	}
	return nil
}

func chirp(s *sprite) {
	fmt.Println("play sprite chirp", s.name)
}

type spriteanimation struct {
	in   bool
	step int
}
