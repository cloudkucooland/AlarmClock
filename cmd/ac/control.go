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
	{
		sprite: getSprite("Pinwheel", "Weather", weatherDialog),
	},
}

func (g *Game) drawControls(screen *ebiten.Image) {
	controls[0].setLocation(screensize.X-100, 20)
	controls[1].setLocation(screensize.X-100, 120)
	controls[2].setLocation(screensize.X-100, 220)

	for x := range controls {
		controls[x].drawWithLabel(screen)
	}
}
