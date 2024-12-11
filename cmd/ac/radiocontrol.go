package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type radiocontrol struct {
	sprite   *sprite
	label    string
	labelimg *ebiten.Image
	labelloc image.Point
	do       func(*Game)
}

func (r *radiocontrol) getlabel() string {
	return r.label
}

func (r *radiocontrol) setlabelimg(i *ebiten.Image) {
	r.labelimg = i
}

var radiocontrols = []radiocontrol{
	{
		sprite: getSprite("Tea Time"),
		label:  "Pause",
		do:     pausePlayer,
	},
	{
		sprite: getSprite("Swan Mommy"),
		label:  "Play",
		do:     resumePlayer,
	},
	{
		sprite: getSprite("Spring"),
		label:  "Stop",
		do:     stopPlayer,
	},
}

func (g *Game) drawRadioControls(screen *ebiten.Image) {
	if g.radio == nil {
		return
	}

	boxwidth := 160
	boxheight := 120
	borderwidth := 20

	// TODO: base this on sprite size not hardcoded values
	grey := color.RGBA{0xaa, 0xaa, 0xaa, 0x99}
	border := color.RGBA{0x66, 0x66, 0x66, 0x00}

	vector.DrawFilledRect(screen, float32(borderwidth), float32(240), float32(screensize.X-(boxwidth)), float32(boxheight), grey, false)
	vector.StrokeRect(screen, float32(borderwidth), float32(240), float32(screensize.X-(boxwidth)), float32(boxheight), float32(4), border, false)
	vector.StrokeRect(screen, float32(borderwidth)*1.5, float32(250), float32(screensize.X-(boxwidth+borderwidth)), float32(boxheight-borderwidth), float32(2), border, false)

	x := (screensize.X - boxwidth) / 2

	for idx := range radiocontrols {
		if radiocontrols[idx].label == "Play" && g.radio.IsPlaying() {
			continue
		}
		if radiocontrols[idx].label == "Pause" && !g.radio.IsPlaying() {
			continue
		}

		radiocontrols[idx].sprite.setLocation(x, 270)
		radiocontrols[idx].sprite.setScale(spriteScale)
		radiocontrols[idx].sprite.draw(screen)

		if radiocontrols[idx].labelimg == nil {
			genlabel(&(radiocontrols[idx]), color.RGBA{0x33, 0x33, 0x33, 0xee}, controlfont)
		}

		if radiocontrols[idx].labelloc.X == 0 {
			b := radiocontrols[idx].sprite.image.Bounds()
			spritecenterx := int(float64(radiocontrols[idx].sprite.loc.X) + float64(b.Max.X)*spriteScale/2.0)
			lb := radiocontrols[idx].labelimg.Bounds()
			labelcenterx := lb.Max.X / 2
			radiocontrols[idx].labelloc.X = spritecenterx - labelcenterx
			radiocontrols[idx].labelloc.Y = radiocontrols[idx].sprite.loc.Y + int(float64(b.Max.Y)*spriteScale+4.0)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(radiocontrols[idx].labelloc.X), float64(radiocontrols[idx].labelloc.Y))
		screen.DrawImage(radiocontrols[idx].labelimg, op)
		x = x + 100
	}
}
