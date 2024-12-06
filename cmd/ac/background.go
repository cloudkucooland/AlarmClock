package main

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	bgres "github.com/cloudkucooland/AlarmClock/resources/backgrounds"
)

type background struct {
	name  string
	raw   []byte
	image *ebiten.Image
}

var backgrounds = []background{
	{
		name: "One",
		raw:  bgres.Default,
	},
}

func getBackground(name string) *background {
	for x := range backgrounds {
		if backgrounds[x].name == name {
			return &backgrounds[x]
		}
	}
	return &backgrounds[0]
}

func randomBackground() *background {
	return &backgrounds[0]
}

func loadBackgrounds() error {
	for x := range backgrounds {
		img, _, err := image.Decode(bytes.NewReader(backgrounds[x].raw))
		if err != nil {
			return err
		}
		backgrounds[x].image = ebiten.NewImageFromImage(img)
	}
	return nil
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	alpha := 0.75
	if g.inScreenSaver() {
		alpha = 0.10
	}

	op := &colorm.DrawImageOptions{}
	op.Blend = ebiten.BlendCopy
	var cm colorm.ColorM
	cm.Scale(1.0, 1.0, 1.0, alpha)
	colorm.DrawImage(screen, g.background.image, cm, op)
}
