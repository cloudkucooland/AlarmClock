package main

import (
	"context"
	"image"
	"log"
	"net/rpc"
	"os"
	// "os/exec"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

	owm "github.com/briandowns/openweathermap"

	"github.com/cloudkucooland/AlarmClock/ledserver"
)

const screenSaverHz = 5
const normalHz = 60

var hz = normalHz

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

type Game struct {
	state               gameState
	debugString         string
	lastAct             time.Time
	clock               *clock
	background          *ebiten.Image
	backgrounds         map[uint8]*background
	weather             *owm.CurrentWeatherData
	weathercache        *ebiten.Image
	audioContext        *audio.Context
	audioPlayer         *audio.Player
	externalAudioCancel context.CancelFunc
	radiobuttons        map[string]*radiobutton
	selectedStation     *radiobutton
	inSleepCountdown    bool
	config              *Config
	controls            []*control
	radiocontrols       map[string]*radiocontrol
	alarmbuttons        map[string]*alarmbutton
	ledclient           *rpc.Client
}

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
	g.ledFrontOn()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screensize.X, screensize.Y
}

func main() {
	if err := loadfonts(); err != nil {
		log.Fatal(err)
	}

	g := &Game{
		state:   inNormal,
		clock:   &clock{},
		weather: nil,
	}
	g.clock.X = defaultClockLocationX
	g.clock.Y = defaultClockLocationY
	g.audioContext = audio.NewContext(44100)
	g.loadconfig()
	g.setupBackgrounds()
	g.setBackground()
	g.setupControls()
	g.setupRadioButtons()
	g.selectedStation = g.defaultStation()
	g.setupAlarmButtons()
	g.setupRadioControls()

	if client, err := rpc.DialHTTP("unix", ledserver.Pipefile); err != nil {
		log.Printf("led server not connected: %s", err.Error())
	} else {
		g.ledclient = client
	}
	g.ledAllOn()

	// setup clock
	now := time.Now()

	// attempt to get the minute-change correct...
	// ms := now.Sub(now.Truncate(time.Second))
	// g.clock.cyclesSinceTick = 1000 - int(ms.Milliseconds()*hz/1000)
	g.clock.timestring = now.Format(g.config.ClockFormat)

	ebiten.SetWindowSize(screensize.X, screensize.Y)
	if hostname, _ := os.Hostname(); strings.EqualFold(hostname, "birdhouse") {
		ebiten.SetFullscreen(true)
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	}

	// ctx, cancel := context.WithCancel(context.Background())
	go g.runWeather(context.Background())

	ebiten.SetWindowTitle("Alarm Clock")
	ebiten.SetTPS(hz)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
