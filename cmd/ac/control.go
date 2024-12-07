package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type control struct {
	sprite *sprite
	label  string
	do     func(*Game)
}

var controls = []control{
	{
		sprite: getSprite("Mad"),
		label:  "Alarms",
		do:     alarmConfigDialog,
	},
	{
		sprite: getSprite("Happy"),
		label:  "Radio",
		do:     radioDialog,
	},
	{
		sprite: getSprite("Pinwheel"),
		label:  "Weather",
		do:     weatherDialog,
	},
}

func (g *Game) drawControls(screen *ebiten.Image) {
	controls[0].sprite.setLocation(screensize.X-100, 20)
	controls[1].sprite.setLocation(screensize.X-100, 120)
	controls[2].sprite.setLocation(screensize.X-100, 220)
	controls[0].sprite.setScale(spriteScale)
	controls[1].sprite.setScale(spriteScale)
	controls[2].sprite.setScale(spriteScale)

	for x := range controls {
		controls[x].sprite.draw(screen)
		controls[x].sprite.drawLabel(controls[x].label, screen) // move animation logic to sprite.draw
	}
}
