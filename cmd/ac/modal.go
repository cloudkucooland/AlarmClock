package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type modalbutton struct {
	*spritelabel
	*sprite
}

var modalbuttons = []modalbutton{
	{
		sprite:      getSprite("Indignent", modalclose),
		spritelabel: &spritelabel{label: "Close"},
	},
}

func (g *Game) drawModal(screen *ebiten.Image) {
	grey := color.RGBA{0xaa, 0xaa, 0xaa, 0x99}
	border := color.RGBA{0x66, 0x66, 0x66, 0x00}

	borderwidth := 20

	vector.DrawFilledRect(screen, float32(borderwidth), float32(borderwidth), float32(screensize.X-(borderwidth*2)), float32(screensize.Y-(borderwidth*2)), grey, false)
	vector.StrokeRect(screen, float32(borderwidth), float32(borderwidth), float32(screensize.X-(borderwidth*2)), float32(screensize.Y-(borderwidth*2)), float32(4), border, false)
	vector.StrokeRect(screen, float32(borderwidth)*1.5, float32(borderwidth)*1.5, float32(screensize.X-(borderwidth*3)), float32(screensize.Y-(borderwidth*3)), float32(2), border, false)

	modalbuttons[0].setLocation(screensize.X-(borderwidth*4), screensize.Y-(borderwidth*5)) // should be based on image size not borderwidth
	modalbuttons[0].setScale(spriteScale)

	if modalbuttons[0].labelimg == nil {
		genlabel(&modalbuttons[0], color.RGBA{0x33, 0x33, 0x33, 0xee}, controlfont)
		modalbuttons[0].setScale(spriteScale)
		modalbuttons[0].draw(screen)
		b := modalbuttons[0].image.Bounds()
		spritecenterx := modalbuttons[0].loc.X + int(float64(b.Max.X)*spriteScale/2.0)

		lb := modalbuttons[0].labelimg.Bounds()
		labelcenterx := lb.Max.X / 2
		modalbuttons[0].labelloc.X = spritecenterx - labelcenterx
		modalbuttons[0].labelloc.Y = int(float64(modalbuttons[0].loc.Y) + float64(b.Max.Y)*spriteScale + 4.0)
	}

	// center label below sprite
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(modalbuttons[0].labelloc.X), float64(modalbuttons[0].labelloc.Y))

	modalbuttons[0].draw(screen)
	screen.DrawImage(modalbuttons[0].labelimg, op)
}

func modalclose(g *Game) {
	g.state = inNormal
}
