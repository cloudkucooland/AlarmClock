package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/cloudkucooland/AlarmClock/resources"
)

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
	state   gameState
	lastAct time.Time
	clock   *clock
}

func (g *Game) inScreenSaver() bool {
	return g.state == inScreenSaver
}

var (
	spaceMonoSource *text.GoTextFaceSource
	clockfont       *text.GoTextFace
	controlfont     *text.GoTextFace
)

func (g *Game) Update() error {
	g.clock.cyclesSinceTick = (g.clock.cyclesSinceTick + 1) % 3600 // assumes 60Hz
	if g.clock.cyclesSinceTick == 1 {
		now := time.Now()
		g.clock.timestring = now.Format("03:04")
		if g.inScreenSaver() {
			g.clock.screensaverClockLocation()
		}
		if !g.inScreenSaver() && now.After(g.lastAct.Add(5*time.Minute)) {
			// start screen saver
			g.state = inScreenSaver
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.inScreenSaver() {
			g.lastAct = time.Now()
			g.leaveScreenSaver()
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.lastAct = time.Now()
		if g.inScreenSaver() {
			g.leaveScreenSaver()
		}
		x, y := ebiten.CursorPosition()

		for _, s := range sprites {
			if s.in(x, y) {
				fmt.Printf("in sprite release %s\n", s.name)
				s.do(&s)
			}
		}

		for _, c := range controls {
			if c.in(x, y) {
				fmt.Printf("in control release %s\n", c.label)
				c.do(g)
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		log.Fatal("shutting down")
	}

	return nil
}

func (g *Game) leaveScreenSaver() {
	g.state = inNormal
	g.clock.clockLocationX = defaultClockLocationX
	g.clock.clockLocationY = defaultClockLocationY
}

func (g *Game) Draw(screen *ebiten.Image) {
	// black  := color.RGBA{0x00, 0x00, 0x00, 0x00}
	// white := color.RGBA{0xff, 0xff, 0xff, 0x00}

	if !g.inScreenSaver() {
		g.drawControls(screen)
	}

	g.drawClock(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
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

	g := &Game{
		state: inNormal,
		clock: &clock{
			clockLocationX: defaultClockLocationX,
			clockLocationY: defaultClockLocationY,
		},
	}

	// setup clock
	now := time.Now()
	g.lastAct = now
	// g.state = inNormal
	// attempt to get the minute-change correct...
	ms := now.Sub(now.Truncate(time.Second))
	g.clock.cyclesSinceTick = int(ms.Milliseconds() * 60 / 1000) // assumes 60Hz

	g.clock.timestring = now.Format("03:04")

	ebiten.SetWindowSize(800, 480)
	// ebiten.SetFullscreen(true)

	ebiten.SetWindowTitle("Alarm Clock")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
