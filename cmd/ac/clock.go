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

		// Soft drop shadow pass for readability over bright wallpapers
		shadowOp := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				LineSpacing: 16,
			},
		}
		shadowOp.GeoM.Translate(float64(g.clock.X)+3, float64(g.clock.Y)+4)
		shadowOp.ColorScale.ScaleWithColor(textShadow)
		text.Draw(g.clock.cache, g.clock.timestring, clockfont, shadowOp)

		// Crisp foreground pass
		op := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				LineSpacing: 16,
			},
		}
		op.GeoM.Translate(float64(g.clock.X), float64(g.clock.Y))
		op.ColorScale.ScaleWithColor(textPrimary)
		text.Draw(g.clock.cache, g.clock.timestring, clockfont, op)
	}

	drawOp := &ebiten.DrawImageOptions{}
	if g.state == inScreenSaver {
		drawOp.ColorScale.ScaleAlpha(0.50)
	}
	screen.DrawImage(g.clock.cache, drawOp)
}

