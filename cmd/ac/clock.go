package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	// "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
)

type clock struct {
	image.Point
	timestring      string
	cyclesSinceTick int
}

const defaultClockLocationX = 50
const defaultClockLocationY = 0

func (c *clock) screensaverClockLocation() {
	// get clock size, determine max range...
	c.X = rand.Int() % 150
	c.Y = rand.Int() % 250
}

func (g *Game) drawClock(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(g.clock.X), float64(g.clock.Y))
	if g.state == inScreenSaver {
		op.ColorScale.ScaleAlpha(0.25)
	}
	op.LineSpacing = clockfont.Size * 0.25
	text.Draw(screen, g.clock.timestring, clockfont, op)
}
