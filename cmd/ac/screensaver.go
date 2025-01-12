package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const screenSaverHz = 5
const normalHz = 60

var screensaverTimeout = 2 * time.Minute

func (g *Game) startScreenSaver() {
	hz = screenSaverHz
	ebiten.SetTPS(hz)
	g.debugString = ""
	g.state = inScreenSaver
	g.clock.clearCache()
	g.setBackground()
	g.ledAllOff()
}

func (g *Game) leaveScreenSaver() {
	hz = normalHz
	ebiten.SetTPS(hz)
	g.debugString = ""
	g.state = inNormal
	g.clock.clearCache()
	g.clock.X = defaultClockLocationX
	g.clock.Y = defaultClockLocationY
	g.setBackground()
	g.ledAllOn()
}
