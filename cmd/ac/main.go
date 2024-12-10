package main

import (
	"context"
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const hz = 60

var screensize = image.Point{
	X: 800,
	Y: 480,
}

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
	state        gameState
	lastAct      time.Time
	clock        *clock
	background   *ebiten.Image
	weather      string
	audioContext *audio.Context
	radio        *audio.Player
}

func (g *Game) startScreenSaver() {
	g.state = inScreenSaver
	g.clock.clearCache()
	g.setBackground()
}

func (g *Game) leaveScreenSaver() {
	g.state = inNormal
	g.clock.clearCache()
	g.clock.X = defaultClockLocationX
	g.clock.Y = defaultClockLocationY
	g.setBackground()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	if err := loadfonts(); err != nil {
		log.Fatal(err)
	}

	g := &Game{
		state:   inAlarm,
		clock:   &clock{},
		weather: "Not loaded",
	}
	g.setBackground()
	g.clock.X = defaultClockLocationX
	g.clock.Y = defaultClockLocationY
	g.audioContext = audio.NewContext(44100)

	// setup clock
	now := time.Now()

	// attempt to get the minute-change correct...
	ms := now.Sub(now.Truncate(time.Second))
	g.clock.cyclesSinceTick = int(ms.Milliseconds() * hz / 1000)
	g.clock.timestring = now.Format("03:04")

	ebiten.SetWindowSize(screensize.X, screensize.Y)
	// ebiten.SetFullscreen(true)

	// ctx, cancel := context.WithCancel(context.Background())
	go g.runWeather(context.Background())

	ebiten.SetWindowTitle("Alarm Clock")
	ebiten.SetTPS(hz)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
