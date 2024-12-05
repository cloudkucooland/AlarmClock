package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	// "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
)

type clock struct {
	clockLocationX  int
	clockLocationY  int
	timestring      string
	cyclesSinceTick int
}

const defaultClockLocationX = 5
const defaultClockLocationY = 20

func (c *clock) screensaverClockLocation() {
	// get clock size, determine max range...
	c.clockLocationX = rand.Int() % 100
	c.clockLocationY = rand.Int() % 200
}

func (g *Game) drawClock(screen *ebiten.Image) {
	c := g.clock

	// w, h := text.Measure(g.clock.timestring, clockfont, clockfont.Size*1.2)
	// vector.DrawFilledRect(screen, g.clock.clockLocationX, g.clock.clockLocationY, float32(w), float32(h), black, false)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(c.clockLocationX), float64(c.clockLocationY))
	if g.inScreenSaver() {
		op.ColorScale.ScaleAlpha(0.25)
	}
	op.LineSpacing = clockfont.Size * 1
	text.Draw(screen, c.timestring, clockfont, op)
}
