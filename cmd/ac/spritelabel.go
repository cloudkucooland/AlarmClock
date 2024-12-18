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

const vectorlabelpadding = float32(1.0)

func (s *spritelabel) genlabelimage(bgcolor color.RGBA, font *text.GoTextFace) {
	textwidth, textheight := text.Measure(s.label, font, 0)
	s.labelimg = buttonbase(float32(textwidth), float32(textheight), bgcolor, vectorlabelpadding)
	bsize := s.labelimg.Bounds()

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bsize.Max.X/2)-float64(textwidth/2), float64(bsize.Max.Y/2)-float64(textheight/2))
	text.Draw(s.labelimg, s.label, font, op)
}

func (s *sprite) setLabel(l string) {
	s.spritelabel.label = l
	if s.spritelabel.labelimg != nil {
		s.spritelabel.labelimg.Deallocate()
		s.spritelabel.labelimg = nil
	}
}
