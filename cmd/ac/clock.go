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

const defaultClockLocationX = 50 
const defaultClockLocationY = 0

func (c *clock) screensaverClockLocation() {
	// get clock size, determine max range...
	c.clockLocationX = rand.Int() % 100
	c.clockLocationY = rand.Int() % 200
}

func (g *Game) drawClock(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(g.clock.clockLocationX), float64(g.clock.clockLocationY))
	if g.state == inScreenSaver {
		op.ColorScale.ScaleAlpha(0.25)
	}
	op.LineSpacing = clockfont.Size * 0.25
	text.Draw(screen, g.clock.timestring, clockfont, op)
}
