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
		{
			sprite: getSprite("Artist", "Music", musicDialog),
		},
	}
}

func (g *Game) drawControls(screen *ebiten.Image) {
	x := screensize.X - 100
	y := 22
	spacing := 125

	for _, c := range g.controls {
		c.setLocation(x, y)
		c.drawWithLabel(screen)
		y = y + spacing
	}
}
