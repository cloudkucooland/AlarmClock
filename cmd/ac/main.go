package main

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/cloudkucooland/AlarmClock/resources"
)

const hz = 60

type gameState int

const (
	inNormal gameState = iota
	inAlarm
	inSnooze
	inScreenSaver
	inAlarmConfig
	inWeather
	inRadio
)

type Game struct {
	state      gameState
	lastAct    time.Time
	clock      *clock
	background *background
	weather    string
}

var (
	spaceMonoSource *text.GoTextFaceSource
	clockfont       *text.GoTextFace
	controlfont     *text.GoTextFace
)

func (g *Game) leaveScreenSaver() {
	g.state = inNormal
	g.clock.clockLocationX = defaultClockLocationX
	g.clock.clockLocationY = defaultClockLocationY
	g.background = randomBackground()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(resources.SpaceMonoBold_ttf))
	if err != nil {
		log.Fatal(err)
	}
	spaceMonoSource = s

	clockfont = &text.GoTextFace{
		Source: spaceMonoSource,
		Size:   192,
	}
	controlfont = &text.GoTextFace{
		Source: spaceMonoSource,
		Size:   12,
	}

	if err = loadSprites(); err != nil {
		log.Fatal(err)
	}
	if err = buildControls(); err != nil {
		log.Fatal(err)
	}
	if err = loadBackgrounds(); err != nil {
		log.Fatal(err)
	}

	g := &Game{
		state: inNormal,
		clock: &clock{
			clockLocationX: defaultClockLocationX,
			clockLocationY: defaultClockLocationY,
		},
		weather: "Not loaded",
	}
	g.background = randomBackground()

	// setup clock
	now := time.Now()

	// g.lastAct = now
	// g.state = inNormal
	g.state = inScreenSaver

	// attempt to get the minute-change correct...
	ms := now.Sub(now.Truncate(time.Second))
	g.clock.cyclesSinceTick = int(ms.Milliseconds() * hz / 1000)
	g.clock.timestring = now.Format("03:04")

	ebiten.SetWindowSize(800, 480)
	// ebiten.SetFullscreen(true)

	// ctx, cancel := context.WithCancel(context.Background())
	go g.runWeather(context.Background())

	ebiten.SetWindowTitle("Alarm Clock")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
