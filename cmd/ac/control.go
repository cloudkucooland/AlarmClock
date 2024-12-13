package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type control struct {
	*spritelabel
	*sprite
}

var controls = []control{
	{
		sprite:      getSprite("Mad", alarmConfigDialog),
		spritelabel: &spritelabel{label: "Alarms"},
	},
	{
		sprite:      getSprite("Happy", radioDialog),
		spritelabel: &spritelabel{label: "Radio"},
	},
	{
		sprite:      getSprite("Pinwheel", weatherDialog),
		spritelabel: &spritelabel{label: "Weather"},
	},
}

func (g *Game) drawControls(screen *ebiten.Image) {
	controls[0].setLocation(screensize.X-100, 20)
	controls[1].setLocation(screensize.X-100, 120)
	controls[2].setLocation(screensize.X-100, 220)

	for x := range controls {
		if controls[x].labelimg == nil {
			genlabel(&(controls[x]), color.RGBA{0x33, 0x33, 0x33, 0xee}, controlfont)
		}
		controls[x].setScale(spriteScale)
		controls[x].draw(screen)

		if controls[x].labelloc.X == 0 {
			b := controls[x].image.Bounds()
			spritecenterx := int(float64(controls[x].loc.X) + float64(b.Max.X)*spriteScale/2.0)
			lb := controls[x].labelimg.Bounds()
			labelcenterx := lb.Max.X / 2
			controls[x].labelloc.X = spritecenterx - labelcenterx
			controls[x].labelloc.Y = controls[x].loc.Y + int(float64(b.Max.Y)*spriteScale+4.0)
		}

		// center label below sprite
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(controls[x].labelloc.X), float64(controls[x].labelloc.Y))
		screen.DrawImage(controls[x].labelimg, op)
	}
}
