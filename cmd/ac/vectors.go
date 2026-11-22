package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const vectorbuttonpadding = float32(10.0)

func button(label string, bgcolor color.RGBA, font *text.GoTextFace) *ebiten.Image {
	textwidth, textheight := text.Measure(label, font, 0)
	b := buttonbase(float32(textwidth), float32(textheight), bgcolor, vectorbuttonpadding)
	bsize := b.Bounds()

	tx := float64(bsize.Max.X)/2.0 - textwidth/2.0
	ty := float64(bsize.Max.Y)/2.0 - textheight/2.0

	// Subtle 1px drop shadow for crisp legibility
	shadowOp := &text.DrawOptions{}
	shadowOp.GeoM.Translate(tx, ty+1)
	shadowOp.ColorScale.ScaleWithColor(textShadow)
	text.Draw(b, label, font, shadowOp)

	// Main foreground label
	op := &text.DrawOptions{}
	op.GeoM.Translate(tx, ty)
	op.ColorScale.ScaleWithColor(textPrimary)
	text.Draw(b, label, font, op)
	return b
}

func buttonbase(textwidth float32, textheight float32, bgcolor color.RGBA, padding float32) *ebiten.Image {
	border := float32(1.5)
	totalheight := textheight + padding*2
	radius := totalheight / 2
	totalwidth := radius + padding + textwidth + padding + radius

	w := int(math.Ceil(float64(totalwidth)))
	h := int(math.Ceil(float64(totalheight)))
	img := ebiten.NewImage(w, h)

	drawRoundedBox(img, 0, 0, totalwidth, totalheight, radius, bgcolor, rimLight, border)
	return img
}

// drawRoundedBox draws an anti-aliased rounded rectangle with an optional border.
func drawRoundedBox(dst *ebiten.Image, x, y, w, h, radius float32, fillCol, borderCol color.RGBA, borderWidth float32) {
	if radius > w/2 {
		radius = w / 2
	}
	if radius > h/2 {
		radius = h / 2
	}

	if borderWidth > 0 && borderCol.A > 0 {
		fillRoundedRect(dst, x, y, w, h, radius, borderCol)
		innerR := radius - borderWidth
		if innerR < 0 {
			innerR = 0
		}
		fillRoundedRect(dst, x+borderWidth, y+borderWidth, w-borderWidth*2, h-borderWidth*2, innerR, fillCol)
		return
	}

	fillRoundedRect(dst, x, y, w, h, radius, fillCol)
}

func fillRoundedRect(dst *ebiten.Image, x, y, w, h, r float32, col color.RGBA) {
	if w <= 0 || h <= 0 {
		return
	}
	if r <= 0 {
		vector.FillRect(dst, x, y, w, h, col, true)
		return
	}

	// Four corner circles
	vector.FillCircle(dst, x+r, y+r, r, col, true)
	vector.FillCircle(dst, x+w-r, y+r, r, col, true)
	vector.FillCircle(dst, x+r, y+h-r, r, col, true)
	vector.FillCircle(dst, x+w-r, y+h-r, r, col, true)

	// Center cross rectangles connecting the corners
	vector.FillRect(dst, x+r, y, w-r*2, h, col, true)
	vector.FillRect(dst, x, y+r, w, h-r*2, col, true)
}

