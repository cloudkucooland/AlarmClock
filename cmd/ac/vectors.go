package main

import (
	// "fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	// "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type vp struct {
	x float32
	y float32
}

const vectorbuttonpadding = float32(10.0)

func button(label string, bgcolor color.RGBA, font *text.GoTextFace) *ebiten.Image {
	textwidth, textheight := text.Measure(label, font, 5)
	b := buttonbase(float32(textwidth), float32(textheight), bgcolor)
	bsize := b.Bounds()

	// draw text centered in button
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bsize.Max.X/2)-float64(textwidth/2), float64(bsize.Max.Y/2)-float64(textheight/2))
	text.Draw(b, label, font, op)
	return b
}

func buttonbase(textwidth float32, textheight float32, bgcolor color.RGBA) *ebiten.Image {
	totalheight := textheight + vectorbuttonpadding*2
	diameter := totalheight
	radius := diameter / 2
	totalwidth := radius + vectorbuttonpadding + textwidth + vectorbuttonpadding + radius

	img := ebiten.NewImage(int(math.Ceil(float64(totalwidth))), int(math.Ceil(float64(totalheight))))
	vector.DrawFilledCircle(img, radius, radius, radius, bgcolor, true)
	vector.DrawFilledCircle(img, totalwidth-radius, radius, radius, bgcolor, true)
	vector.DrawFilledRect(img, radius, 0, totalwidth-diameter, totalheight, bgcolor, true)
	return img
}

func roundbutton(screen *ebiten.Image, x int, y int, height int, width int, bgcolor color.RGBA) {
	var p vector.Path

	radius := float32(15)

	// left-most
	f := vp{x: float32(x), y: float32(y)}
	p.MoveTo(f.x, f.y)
	// top-left corner
	p0 := vp{x: f.x - radius, y: f.y + radius}
	p.ArcTo(f.x, f.y, p0.x, p0.y, radius)
	// top right
	p1 := vp{x: p0.x + float32(width) - 2*radius, y: p0.y}
	p.LineTo(p1.x, p1.y)
	// right-most
	p2 := vp{x: p1.x + radius, y: p1.y + float32(height/2)}
	// bottom right
	p.ArcTo(p1.x, p1.y, p2.x, p2.y, radius)
	p3 := vp{x: p1.x, y: p1.y + float32(height)}
	// bottom left
	p.ArcTo(p2.x, p2.y, p3.x, p3.y, radius)
	// back to (almost) start
	p.ArcTo(p3.x, p3.y, f.x-0.01, f.y-0.01, radius)
	// finish
	p.Close()

	// vector.DrawFilledPath(screen, &p, bgcolor, true, vector.FillRuleNonZero)
}
