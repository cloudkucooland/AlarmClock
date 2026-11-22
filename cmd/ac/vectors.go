package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	// "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const vectorbuttonpadding = float32(10.0)

func button(label string, bgcolor color.RGBA, font *text.GoTextFace) *ebiten.Image {
	textwidth, textheight := text.Measure(label, font, 0)
	b := buttonbase(float32(textwidth), float32(textheight), bgcolor, vectorbuttonpadding)
	bsize := b.Bounds()

	// draw text centered in button
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bsize.Max.X/2)-float64(textwidth/2), float64(bsize.Max.Y/2)-float64(textheight/2))
	text.Draw(b, label, font, op)
	return b
}

func buttonbase(textwidth float32, textheight float32, bgcolor color.RGBA, padding float32) *ebiten.Image {
	border := float32(2.0)
	totalheight := textheight + padding*2
	diameter := totalheight
	radius := diameter / 2
	totalwidth := radius + padding + textwidth + padding + radius

	w := int(math.Ceil(float64(totalwidth)))
	h := int(math.Ceil(float64(totalheight)))
	img := ebiten.NewImage(w, h)

	// Subtle dark border pill
	borderCol := color.RGBA{R: 0, G: 0, B: 0, A: 90}
	vector.FillCircle(img, radius, radius, radius, borderCol, true)
	vector.FillCircle(img, totalwidth-radius, radius, radius, borderCol, true)
	vector.FillRect(img, radius, 0, totalwidth-diameter, totalheight, borderCol, true)

	// Inner fill pill inset by `border`
	innerR := radius - border
	vector.FillCircle(img, radius, radius, innerR, bgcolor, true)
	vector.FillCircle(img, totalwidth-radius, radius, innerR, bgcolor, true)
	vector.FillRect(img, radius, border, totalwidth-diameter, totalheight-border*2, bgcolor, true)

	return img
}
