package main

import (
	// "fmt"
	"image"
	// "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
)

type alarm struct {
	alarmTime  alarmTime
	enabled    bool
	station    *radiobutton
	triggered  bool
	sleep      bool
	sleepCount uint8
	sleepTime  alarmTime
}

type alarmTime struct {
	hour   int
	minute int
}

type alarmbutton struct {
	image.Rectangle
}

func alarmConfigDialog(g *Game) {
	g.state = inAlarmConfig
}

func (g *Game) drawAlarmConfig(screen *ebiten.Image) {
	g.drawModal(screen)
}
