package main

import (
	// "fmt"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackground(screen)

	switch g.state {
	case inNormal:
		g.drawClock(screen)
		g.drawControls(screen)
		g.drawRadioControls(screen)
	case inAlarm:
		g.drawAlarm(screen)
	case inSnooze:
		g.drawSnooze(screen)
	case inScreenSaver:
		g.drawClock(screen)
	case inAlarmConfig:
		g.drawAlarmConfig(screen)
	case inWeather:
		g.drawWeather(screen)
	case inRadio:
		g.drawRadioDialog(screen)
	}

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}
