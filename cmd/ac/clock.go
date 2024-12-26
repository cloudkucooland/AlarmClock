package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type clock struct {
	image.Point
	timestring      string
	cyclesSinceTick int
	cache           *ebiten.Image
}

const defaultClockLocationX = 50
const defaultClockLocationY = 0

func (c *clock) screensaverClockLocation() {
	// get clock size, determine max range...
	// #nosec G404 - we don't need strong randomness
	c.X = rand.Int() % 150
	// #nosec G404 - we don't need strong randomness
	c.Y = rand.Int() % 250
}

func (c *clock) clearCache() {
	if c.cache == nil {
		return
	}

	c.cache.Deallocate()
	c.cache = nil
}

func (g *Game) drawClock(screen *ebiten.Image) {
	if g.clock.cache == nil {
		g.clock.cache = ebiten.NewImage(screensize.X, screensize.Y)
		op := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				LineSpacing: 16,
			},
		}
		op.Blend = ebiten.Blend{
			BlendFactorSourceRGB: ebiten.BlendFactorSourceColor,
		}
		op.GeoM.Translate(float64(g.clock.X), float64(g.clock.Y))
		if g.state == inScreenSaver {
			op.ColorScale.ScaleAlpha(0.25)
		}
		text.Draw(g.clock.cache, g.clock.timestring, clockfont, op)
	}
	screen.DrawImage(g.clock.cache, &ebiten.DrawImageOptions{})
}
