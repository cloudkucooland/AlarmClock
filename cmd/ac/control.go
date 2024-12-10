package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type control struct {
	sprite   *sprite
	label    string
	labelimg *ebiten.Image
	labelloc image.Point
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

		if controls[x].labelloc.X == 0 {
			b := controls[x].sprite.image.Bounds()
			spritecenterx := int(float64(controls[x].sprite.loc.X) + float64(b.Max.X)*spriteScale/2.0)
			lb := controls[x].labelimg.Bounds()
			labelcenterx := lb.Max.X / 2
			controls[x].labelloc.X = spritecenterx - labelcenterx
			controls[x].labelloc.Y = controls[x].sprite.loc.Y + int(float64(b.Max.Y)*spriteScale + 4.0)
		}

		// center label below sprite
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(controls[x].labelloc.X), float64(controls[x].labelloc.Y))
		screen.DrawImage(controls[x].labelimg, op)
	}
}
