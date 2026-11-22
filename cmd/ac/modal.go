package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	modalMargin  = 24
	modalPadding = 16
	modalRadius  = float32(18)
)

type modalbutton struct {
	*sprite
}

var modalbuttons = []*modalbutton{
	{
		sprite: getSprite("Indignent", "Close", modalclose),
	},
}

func modalBounds() (x, y, w, h int) {
	return modalMargin, modalMargin, screensize.X - (modalMargin * 2), screensize.Y - (modalMargin * 2)
}

func (g *Game) drawModal(screen *ebiten.Image) {
	// Dim underlying wallpaper/clock so the modal has clear depth
	vector.FillRect(screen, 0, 0, float32(screensize.X), float32(screensize.Y), scrimBlack, false)

	mx, my, mw, mh := modalBounds()
	drawRoundedBox(screen, float32(mx), float32(my), float32(mw), float32(mh), modalRadius, modalSlate, rimLight, 1.5)

	// Anchor Close button dynamically in the top-right corner of the modal card
	closeBtn := modalbuttons[0]
	b := closeBtn.image.Bounds()
	scaledW := int(float64(b.Max.X) * closeBtn.scale)
	closeBtn.setLocation(mx+mw-scaledW-modalPadding, my+modalPadding)
	closeBtn.drawWithLabel(screen)
}

func modalclose(g *Game) {
	chirp(g)
	g.state = inNormal
}

