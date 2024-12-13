package main

import (
	"bytes"
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	spriteres "github.com/cloudkucooland/AlarmClock/resources/sprites"
)

const (
	spriteScale = 1.5 // the default scale for the icons
)

type sprite struct {
	name  string
	raw   []byte
	image *ebiten.Image
	loc   image.Point
	scale float64
	do    func(*Game)
	ani   *spriteanimation
}

func (s *sprite) in(x, y int) bool {
	h := float64(32) // get from image
	w := float64(32) // get from image
	return (x >= s.loc.X &&
		float64(x) <= float64(s.loc.X)+w*s.scale) &&
		(y >= s.loc.Y && float64(y) <= float64(s.loc.Y)+h*s.scale)
}

func (s *sprite) setLocation(x, y int) {
	s.loc.X = x
	s.loc.Y = y
}

func (s *sprite) setScale(scale float64) {
	s.scale = scale
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

func getSprite(name string, do func(*Game)) *sprite {
	out := sprite{
		name: "Uninitialized",
		ani:  &spriteanimation{},
		do:   do,
	}

	for x := range sprites {
		if sprites[x].name == name {
			out.name = sprites[x].name
			out.ani = &spriteanimation{}
			img, _, err := image.Decode(bytes.NewReader(sprites[x].raw))
			if err != nil {
				panic(err.Error())
			}
			out.image = ebiten.NewImageFromImage(img)
			if out.do == nil {
				out.do = chirp
			}
			return &out
		}
	}
	panic("unable to find sprite")
	return &out
}

func (s *sprite) draw(screen *ebiten.Image) {
	if s.image == nil {
		fmt.Printf("%+v\n", s)
		panic("null sprite image")
	}

	if s.ani.in {
		s.aniStep(screen)
	} else {
		// draw still
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(s.scale, s.scale)
		op.GeoM.Translate(float64(s.loc.X), float64(s.loc.Y))
		screen.DrawImage(s.image, op)
	}
}

func (s *sprite) drawLabel(label string, screen *ebiten.Image) {
	space := float64(s.image.Bounds().Max.Y) * s.scale

	top := &text.DrawOptions{}
	top.GeoM.Translate(float64(s.loc.X), float64(s.loc.Y)+space)
	top.LineSpacing = controlfont.Size
	text.Draw(screen, label, controlfont, top)
}

func chirp(g *Game) {
	fmt.Println("play sprite chirp")
}

type spriteanimation struct {
	in   bool
	step int
}

func (s *sprite) aniStep(screen *ebiten.Image) {
	s.ani.step = s.ani.step + 1

	scale := spriteScale + scaleWibble(float64(s.ani.step))
	theta := thetaWibble(float64(s.ani.step))
	recenterx, recentery := locWibble(float64(s.loc.X), float64(s.loc.Y), float64(s.ani.step))

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Rotate(theta)
	op.GeoM.Translate(recenterx, recentery)
	screen.DrawImage(s.image, op)

	if s.ani.step > (hz / 4) { // quarter of a second
		s.ani.step = 0
		s.ani.in = false
	}
}

func scaleWibble(i float64) float64 {
	return math.Sin(i/4) / 3
}

func thetaWibble(i float64) float64 {
	return math.Sin(i/6) / 6
}

func locWibble(x, y, step float64) (float64, float64) {
	z := thetaWibble(step) * 25
	return x + z, y + z
}

func (s *sprite) startanimation() {
	if s.ani.step != 0 || s.ani.in {
		return
	}
	s.ani.step = 1
	s.ani.in = true
}

func (s *sprite) stopanimation() {
	s.ani.step = 0
	s.ani.in = false
}
