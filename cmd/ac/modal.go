package main

import (
	// "fmt"
	// "image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type modalbutton struct {
	sprite *sprite
	label  string
	x      int
	y      int
	do     func(*Game)
}

var modalbuttons = []modalbutton{
	{
		sprite: getSprite("Indignent"),
		label:  "Close",
		do:     modalclose,
	},
	{
		sprite: getSprite("Happy"),
		label:  "OK",
		do:     modalok,
	},
}

func (g *Game) drawModal(screen *ebiten.Image) {
	grey := color.RGBA{0xaa, 0xaa, 0xaa, 0x99}
	border := color.RGBA{0x66, 0x66, 0x66, 0x00}
	vector.DrawFilledRect(screen, float32(20), float32(20), float32(760), float32(440), grey, false)
	vector.StrokeRect(screen, float32(20), float32(20), float32(760), float32(440), float32(4), border, false)
	vector.StrokeRect(screen, float32(30), float32(30), float32(740), float32(420), float32(2), border, false)

	modalbuttons[0].x = 720
	modalbuttons[0].y = 380

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(controlScale, controlScale)
	op.GeoM.Translate(float64(modalbuttons[0].x), float64(modalbuttons[0].y))
	screen.DrawImage(modalbuttons[0].sprite.image, op)

	top := &text.DrawOptions{}
	top.GeoM.Translate(float64(modalbuttons[0].x), float64(modalbuttons[0].y+controlYspace))
	top.LineSpacing = controlfont.Size
	text.Draw(screen, modalbuttons[0].label, controlfont, top)
}

func (r modalbutton) in(x, y int) bool {
	return (x >= r.x && x <= r.x+controlIconY*controlScale) && (y >= r.y && y <= r.y+controlIconY*controlScale)
	// return r.sprite.in(x, y)
}

func (m *modalbutton) modaldo(g *Game) {
	m.do(g)
}

func modalclose(g *Game) {
	g.state = inNormal
}

func modalok(g *Game) {
	g.state = inNormal
}
