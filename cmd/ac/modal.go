package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type modalbutton struct {
	sprite *sprite
	label  string
	do     func(*Game)
}

var modalbuttons = []modalbutton{
	{
		sprite: getSprite("Indignent"),
		label:  "Close",
		do:     modalclose,
	},
}

func (g *Game) drawModal(screen *ebiten.Image) {
	grey := color.RGBA{0xaa, 0xaa, 0xaa, 0x99}
	border := color.RGBA{0x66, 0x66, 0x66, 0x00}

	borderwidth := 20

	vector.DrawFilledRect(screen, float32(borderwidth), float32(borderwidth), float32(screensize.X-(borderwidth*2)), float32(screensize.Y-(borderwidth*2)), grey, false)
	vector.StrokeRect(screen, float32(borderwidth), float32(borderwidth), float32(screensize.X-(borderwidth*2)), float32(screensize.Y-(borderwidth*2)), float32(4), border, false)
	vector.StrokeRect(screen, float32(borderwidth)*1.5, float32(borderwidth)*1.5, float32(screensize.X-(borderwidth*3)), float32(screensize.Y-(borderwidth*3)), float32(2), border, false)

	modalbuttons[0].sprite.setLocation(screensize.X-(borderwidth*4), screensize.Y-(borderwidth*5)) // should be based on image size not borderwidth
	modalbuttons[0].sprite.setScale(spriteScale)

	modalbuttons[0].sprite.draw(screen)
	modalbuttons[0].sprite.drawLabel(modalbuttons[0].label, screen)
}

func (m *modalbutton) modaldo(g *Game) {
	m.do(g)
}

func modalclose(g *Game) {
	g.state = inNormal
}
