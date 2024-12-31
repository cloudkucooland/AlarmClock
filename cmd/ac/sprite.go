package main

import (
	"bytes"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/cloudkucooland/AlarmClock/resources/sounds"
	spriteres "github.com/cloudkucooland/AlarmClock/resources/sprites"
)

const (
	spriteScale = 2.5 // the default scale for the icons
	drawhitbox  = false
)

type sprite struct {
	name  string
	image *ebiten.Image
	image.Point
	scale float64
	do    func(*Game)
	ani   *spriteanimation
	*spritelabel
}

func (s *sprite) in(x, y int) bool {
	b := s.image.Bounds()
	w := int((float64(b.Max.X) * s.scale))
	h := int((float64(b.Max.Y) * s.scale))
	return (x >= s.X && x <= s.X+w && y >= s.Y && y <= s.Y+h)
}

func (s *sprite) setLocation(x, y int) {
	s.X = x
	s.Y = y
}

func getSprite(name string, label string, do func(*Game)) *sprite {
	out := sprite{
		name:        name,
		ani:         &spriteanimation{},
		do:          do,
		scale:       spriteScale,
		spritelabel: &spritelabel{label: label},
	}

	raw, ok := spriteres.RawSprites[name]
	if !ok {
		panic("unable to find sprite")
	}

	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		panic(err.Error())
	}
	out.image = ebiten.NewImageFromImage(img)

	if out.do == nil {
		out.do = chirp
	}
	out.label = label
	if out.label != "" && controlfont != nil {
		// controlfont doesn't seem to be loaded early in the startup sequence
		out.genlabelimage(spritelabelgrey, controlfont)
	}
	return &out
}

func (s *sprite) draw(screen *ebiten.Image) {
	if s.ani.in {
		s.aniStep(screen)
	} else {
		// draw still
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(s.scale, s.scale)
		op.GeoM.Translate(float64(s.X), float64(s.Y))
		screen.DrawImage(s.image, op)

		// draw hitbox
		if drawhitbox {
			b := s.image.Bounds()
			w := int((float64(b.Max.X) * s.scale))
			h := int((float64(b.Max.Y) * s.scale))
			vector.StrokeRect(screen, float32(s.X), float32(s.Y), float32(w), float32(h), 1, color.Black, false)
		}
	}
}

func (s *sprite) drawlabel(screen *ebiten.Image) {
	if s.label == "" {
		return
	}

	if s.labelimg == nil {
		s.genlabelimage(spritelabelgrey, controlfont)
	}

	if s.labelloc.X == 0 {
		s.setlabelloc()
	}

	// center label below sprite
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.labelloc.X), float64(s.labelloc.Y))
	screen.DrawImage(s.labelimg, op)
}

func (s *sprite) setlabelloc() {
	if s.label == "" {
		return
	}

	if s.image == nil {
		return
	}

	b := s.image.Bounds()
	spritecenterx := int(float64(s.X) + float64(b.Max.X)*s.scale/2.0)
	lb := s.labelimg.Bounds()
	labelcenterx := lb.Max.X / 2
	s.labelloc.X = spritecenterx - labelcenterx
	s.labelloc.Y = s.Y + int(float64(b.Max.Y)*s.scale) + 4
}

func (s *sprite) drawWithLabel(screen *ebiten.Image) {
	s.draw(screen)
	s.drawlabel(screen)
}

func chirp(g *Game) {
	chirpraw, ok := sounds.Sounds["Khew"]
	if !ok {
		return
	}
	chirp, err := mp3.DecodeWithoutResampling(bytes.NewReader(chirpraw))
	if err != nil {
		g.debug(err.Error())
		return
	}

	p, err := g.audioContext.NewPlayer(chirp)
	if err != nil {
		g.debug(err.Error())
		return
	}
	p.SetVolume(0.20)
	p.Play()
	/* if err := p.Close(); err != nil {
		g.debug(err.Error())
	} */
}


func tick(g *Game) {
	raw, ok := sounds.Sounds["Tick"]
	if !ok {
		return
	}
	tick, err := mp3.DecodeWithoutResampling(bytes.NewReader(raw))
	if err != nil {
		g.debug(err.Error())
		return
	}

	p, err := g.audioContext.NewPlayer(tick)
	if err != nil {
		g.debug(err.Error())
		return
	}
	p.SetVolume(0.60)
	p.Play()
}

type spriteanimation struct {
	in   bool
	step int
}

func (s *sprite) aniStep(screen *ebiten.Image) {
	s.ani.step = s.ani.step + 1

	scale := s.scale + scaleWibble(float64(s.ani.step))
	theta := thetaWibble(float64(s.ani.step))
	recenterx, recentery := locWibble(float64(s.X), float64(s.Y), float64(s.ani.step))

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
