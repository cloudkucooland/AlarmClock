package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type control struct {
	sprite   *sprite
	label    string
	labelimg *ebiten.Image
	do       func(*Game)
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

	for x := range controls {
		if controls[x].labelimg == nil {
			controls[x].genlabel(color.RGBA{0x33, 0x33, 0x33, 0xee}, controlfont)
		}
		controls[x].sprite.setScale(spriteScale)
		controls[x].sprite.draw(screen)
		b := controls[x].sprite.image.Bounds()
		spritecenterx := float64(controls[x].sprite.loc.X) + float64(b.Max.X)*spriteScale/2.0

		lb := controls[x].labelimg.Bounds()
		labelcenterx := float64(lb.Max.X / 2)
		labelx := spritecenterx - labelcenterx
		labely := float64(controls[x].sprite.loc.Y) + float64(b.Max.Y)*spriteScale + 4.0

		// center label below sprite
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(labelx, labely)
		screen.DrawImage(controls[x].labelimg, op)
	}
}
