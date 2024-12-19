package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type control struct {
	*sprite
}

func (g *Game) setupControls() {
	g.controls = []*control{
		{
			sprite: getSprite("Mad", "Alarms", alarmConfigDialog),
		},
		{
			sprite: getSprite("Happy", "Radio", radioDialog),
		},
	}
}

func (g *Game) drawControls(screen *ebiten.Image) {
	x := screensize.X - 100
	y := 30
	spacing := 150

	for _, c := range g.controls {
		c.setLocation(x, y)
		c.drawWithLabel(screen)
		y = y + spacing
	}
}
