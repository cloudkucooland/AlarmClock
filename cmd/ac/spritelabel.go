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

const vectorlabelpadding = float32(3.0)

func (s *spritelabel) genlabelimage(bgcolor color.RGBA, font *text.GoTextFace) {
	textwidth, textheight := text.Measure(s.label, font, 0)
	s.labelimg = buttonbase(float32(textwidth), float32(textheight), bgcolor, vectorlabelpadding)
	bsize := s.labelimg.Bounds()

	tx := float64(bsize.Max.X)/2.0 - textwidth/2.0
	ty := float64(bsize.Max.Y)/2.0 - textheight/2.0

	shadowOp := &text.DrawOptions{}
	shadowOp.GeoM.Translate(tx, ty+1)
	shadowOp.ColorScale.ScaleWithColor(textShadow)
	text.Draw(s.labelimg, s.label, font, shadowOp)

	op := &text.DrawOptions{}
	op.GeoM.Translate(tx, ty)
	op.ColorScale.ScaleWithColor(textPrimary)
	text.Draw(s.labelimg, s.label, font, op)
}

func (s *sprite) setLabel(l string) {
	if s.spritelabel.label == l && s.spritelabel.labelimg != nil {
		return
	}
	s.spritelabel.label = l
	if s.spritelabel.labelimg != nil {
		s.spritelabel.labelimg.Deallocate()
		s.spritelabel.labelimg = nil
	}
}
