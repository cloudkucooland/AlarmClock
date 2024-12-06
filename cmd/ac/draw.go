package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackground(screen)

	switch g.state {
	case inNormal:
		g.drawControls(screen)
		g.drawClock(screen)
	case inAlarm:
		g.drawClock(screen)
	case inSnooze:
		g.drawClock(screen)
	case inScreenSaver:
		g.drawClock(screen)
	case inAlarmConfig:
		g.drawAlarmConfig(screen)
	case inWeather:
		g.drawWeather(screen)
	case inRadio:
		g.drawRadio(screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}
