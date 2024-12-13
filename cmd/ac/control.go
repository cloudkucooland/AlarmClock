package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type control struct {
	*sprite
}

var controls = []control{
	{
		sprite: getSprite("Mad", "Alarms", alarmConfigDialog),
	},
	{
		sprite: getSprite("Happy", "Radio", radioDialog),
	},
}

func (g *Game) drawControls(screen *ebiten.Image) {
	x := screensize.X - 100
	y := 30
	spacing := 100

	for idx := range controls {
		controls[idx].setLocation(x, y)
		controls[idx].drawWithLabel(screen)
		y = y + spacing
	}
}
