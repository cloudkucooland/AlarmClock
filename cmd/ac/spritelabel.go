package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type spritelabel struct {
	label    string
	labelimg *ebiten.Image
	labelloc image.Point
}

func (s *spritelabel) getlabel() string {
	return s.label
}

func (s *spritelabel) setlabelimg(i *ebiten.Image) {
	s.labelimg = i
}

type labelhaver interface {
	getlabel() string
	setlabelimg(*ebiten.Image)
}

const vectorlabelpadding = float32(1.0)

func genlabel(c labelhaver, bgcolor color.RGBA, font *text.GoTextFace) {
	textwidth, textheight := text.Measure(c.getlabel(), font, 0)
	b := buttonbase(float32(textwidth), float32(textheight), bgcolor, vectorlabelpadding)
	bsize := b.Bounds()

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bsize.Max.X/2)-float64(textwidth/2), float64(bsize.Max.Y/2)-float64(textheight/2))
	text.Draw(b, c.getlabel(), font, op)

	c.setlabelimg(b)
}
