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
const vectorlabelpadding = float32(1.0)

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
	totalheight := textheight + padding*2
	diameter := totalheight
	radius := diameter / 2
	totalwidth := radius + padding + textwidth + padding + radius

	img := ebiten.NewImage(int(math.Ceil(float64(totalwidth))), int(math.Ceil(float64(totalheight))))
	vector.DrawFilledCircle(img, radius, radius, radius, bgcolor, true)
	vector.DrawFilledCircle(img, totalwidth-radius, radius, radius, bgcolor, true)
	vector.DrawFilledRect(img, radius, 0, totalwidth-diameter, totalheight, bgcolor, true)
	return img
}

func (c *control) genlabel(bgcolor color.RGBA, font *text.GoTextFace) {
	textwidth, textheight := text.Measure(c.label, font, 0)
	b := buttonbase(float32(textwidth), float32(textheight), bgcolor, vectorlabelpadding)
	bsize := b.Bounds()

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bsize.Max.X/2)-float64(textwidth/2), float64(bsize.Max.Y/2)-float64(textheight/2))
	text.Draw(b, c.label, font, op)

	c.labelimg = b
}

func (m *modalbutton) genlabel(bgcolor color.RGBA, font *text.GoTextFace) {
	textwidth, textheight := text.Measure(m.label, font, 0)
	b := buttonbase(float32(textwidth), float32(textheight), bgcolor, vectorlabelpadding)
	bsize := b.Bounds()

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bsize.Max.X/2)-float64(textwidth/2), float64(bsize.Max.Y/2)-float64(textheight/2))
	text.Draw(b, m.label, font, op)

	m.labelimg = b
}
