package main

import (
	"context"
	"image"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

	owm "github.com/briandowns/openweathermap"
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
	inRadio
)

const clockformat = "3:04"

type Game struct {
	state            gameState
	debugString      string
	lastAct          time.Time
	clock            *clock
	background       *ebiten.Image
	weather          *owm.CurrentWeatherData
	weathercache     *ebiten.Image
	audioContext     *audio.Context
	radio            *audio.Player
	enabledAlarmID   alarmid
	selectedStation  *radiobutton
	inSleepCountdown bool
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
	return screensize.X, screensize.Y
}

func main() {
	if err := loadfonts(); err != nil {
		log.Fatal(err)
	}

	g := &Game{
		state:          inNormal,
		clock:          &clock{},
		weather:        nil,
		enabledAlarmID: -1,
	}
	g.selectedStation = defaultStation()
	g.setBackground()
	g.clock.X = defaultClockLocationX
	g.clock.Y = defaultClockLocationY
	g.audioContext = audio.NewContext(44100)

	// setup clock
	now := time.Now()

	// attempt to get the minute-change correct...
	ms := now.Sub(now.Truncate(time.Second))
	g.clock.cyclesSinceTick = int(ms.Milliseconds() * hz / 1000)
	g.clock.timestring = now.Format(clockformat)

	ebiten.SetWindowSize(screensize.X, screensize.Y)
	if hostname, _ := os.Hostname(); strings.EqualFold(hostname, "birdhouse") {
		ebiten.SetFullscreen(true)
	}

	// ctx, cancel := context.WithCancel(context.Background())
	go g.runWeather(context.Background())

	ebiten.SetWindowTitle("Alarm Clock")
	ebiten.SetTPS(hz)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
