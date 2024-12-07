package main

import (
	"bytes"
	"fmt"
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	bgres "github.com/cloudkucooland/AlarmClock/resources/backgrounds"
)

type background struct {
	name string
	raw  []byte
}

var backgrounds = []background{
	{
		name: "Flower Splash",
		raw:  bgres.Default,
	},
	{
		name: "Owl Moon",
		raw:  bgres.Owlmoon,
	},
	{
		name: "Owl Eyes",
		raw:  bgres.Owleyes,
	},
	{
		name: "Hummingbird",
		raw:  bgres.Hummingbird,
	},
}

func (g *Game) setBackground() error {
	if g.background != nil {
		g.background.Deallocate()
		g.background = nil
	}
	g.background = ebiten.NewImage(screensize.X, screensize.Y)

	// random for now, later we can do by season/weather/time-of-date, etc
	idx := rand.Intn(len(backgrounds))
	raw, _, err := image.Decode(bytes.NewReader(backgrounds[idx].raw))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	img := ebiten.NewImageFromImage(raw)

	// cache dim for screensaver
	if g.state == inScreenSaver {
		op := &colorm.DrawImageOptions{}
		op.Blend = ebiten.BlendCopy
		var cm colorm.ColorM
		cm.Scale(1.0, 1.0, 1.0, 0.10)
		colorm.DrawImage(g.background, img, cm, op)
		return nil
	}

	// cache normal
	g.background.DrawImage(img, &ebiten.DrawImageOptions{})
	return nil
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	if g.background != nil {
		screen.DrawImage(g.background, &ebiten.DrawImageOptions{})
		return
	}
}
