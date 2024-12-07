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
	name  string
	raw   []byte
	image *ebiten.Image
	// cache *ebiten.Image
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

func getBackground(name string) *background {
	for idx := range backgrounds {
		if backgrounds[idx].name == name {
			return &backgrounds[idx]
		}
	}
	return &backgrounds[0]
}

func randomBackground() *background {
	idx := rand.Intn(len(backgrounds))
	return &backgrounds[idx]
}

func loadBackgrounds() error {
	for idx := range backgrounds {
		img, _, err := image.Decode(bytes.NewReader(backgrounds[idx].raw))
		if err != nil {
			return err
		}
		backgrounds[idx].image = ebiten.NewImageFromImage(img)
	}
	return nil
}

// TODO implement cache so we don't do the alpha Blend every tick
func (g *Game) drawBackground(screen *ebiten.Image) {
	alpha := 0.75
	if g.state == inScreenSaver {
		alpha = 0.10
	}

	op := &colorm.DrawImageOptions{}
	op.Blend = ebiten.BlendCopy
	var cm colorm.ColorM
	cm.Scale(1.0, 1.0, 1.0, alpha)
	colorm.DrawImage(screen, g.background.image, cm, op)
}
