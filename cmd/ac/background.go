package main

import (
	"bytes"
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	bgres "github.com/cloudkucooland/AlarmClock/resources/backgrounds"
)

type background struct {
	raw  []byte
	tags []string
}

func (g *Game) setupBackgrounds() {
	g.backgrounds = map[uint8]*background{
		0: {raw: bgres.Default, tags: []string{"Day", "Sunny"}},
		1: {raw: bgres.Owlmoon, tags: []string{"Night"}},
		2: {raw: bgres.Owleyes, tags: []string{"Night"}},
		3: {raw: bgres.Hummingbird, tags: []string{"Day"}},
	}
}

func (g *Game) setBackground() {
	if g.background != nil {
		g.background.Deallocate()
		g.background = nil
	}
	g.background = ebiten.NewImage(screensize.X, screensize.Y)

	// random for now, later we can do by season/weather/time-of-date, etc
	// #nosec G404 G115
	key := uint8(rand.Intn(len(g.backgrounds)))
	bg, ok := g.backgrounds[key]
	if !ok {
		g.debug("backgrounds maps not keyed correctly")
		return
	}
	decoded, _, err := image.Decode(bytes.NewReader(bg.raw))
	if err != nil {
		g.debug(err.Error())
		return
	}
	img := ebiten.NewImageFromImage(decoded)

	// cache dim for screensaver
	if g.state == inScreenSaver {
		op := &colorm.DrawImageOptions{}
		op.Blend = ebiten.BlendCopy
		var cm colorm.ColorM
		cm.Scale(1.0, 1.0, 1.0, 0.10)
		colorm.DrawImage(g.background, img, cm, op)
		return
	}

	// cache normal
	g.background.DrawImage(img, &ebiten.DrawImageOptions{})
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	if g.background != nil {
		screen.DrawImage(g.background, &ebiten.DrawImageOptions{})
		return
	}
}
