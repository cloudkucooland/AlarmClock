package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type vp struct {
	x float32
	y float32
}

func roundbuttontext(screen *ebiten.Image, x int, y int, height int, width int, text string, bgcolor color.RGBA) {
	roundbutton(screen, x, y, height, width, bgcolor)
	// draw text centered in button
}

func roundbutton(screen *ebiten.Image, x int, y int, height int, width int, bgcolor color.RGBA) {
	var p vector.Path

	radius := float32(15)

	f := vp{x: float32(x), y: float32(y)}
	p0 := vp{x: f.x - radius, y: f.y + radius}
	p1 := vp{x: p0.x + float32(width) - 2*radius, y: p0.y}
	p2 := vp{x: p1.x + radius, y: p1.y + float32(width)}

	p.ArcTo(f.x, f.y, p0.x, p0.y, radius)
	p.LineTo(p1.x, p1.y)
	p.ArcTo(p1.x, p1.y, p2.x, p2.y, radius)
	// TODO draw bottom half
	p.Close()

	// why no work?
	// vector.DrawFilledPath(screen, &p, bgcolor, true, vector.FillRuleNonZero)
}
