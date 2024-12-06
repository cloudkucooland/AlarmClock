package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	// update clock every minute
	g.clock.cyclesSinceTick = (g.clock.cyclesSinceTick + 1) % (60 * hz)
	if g.clock.cyclesSinceTick == 1 {
		now := time.Now()
		g.clock.timestring = now.Format("03:04")
		if g.state == inScreenSaver {
			g.clock.screensaverClockLocation()
		}
		if g.state != inScreenSaver && now.After(g.lastAct.Add(5*time.Minute)) {
			// start screen saver
			g.state = inScreenSaver
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		switch g.state {
		case inScreenSaver:
			g.lastAct = time.Now()
			g.leaveScreenSaver()
		case inNormal:
			for idx := range controls {
				if controls[idx].in(x, y) {
					controls[idx].startanimation()
				}
			}
		case inRadio:
			for idx := range radiobuttons {
				if radiobuttons[idx].in(x, y) {
					fmt.Println("in radiobutton mousedown")
				}
			}
			for idx := range modalbuttons {
				if modalbuttons[idx].in(x, y) {
					fmt.Println("in modal mousedown")
				}
			}
		case inAlarmConfig:
			/* for idx := range radiobuttons {
				if radiobuttons[idx].in(x, y) {
					fmt.Println("in radiobutton mousedown")
				}
			} */
			for idx := range modalbuttons {
				if modalbuttons[idx].in(x, y) {
					fmt.Println("in modal mousedown")
				}
			}
		case inWeather:
			// nothing
		default:
			fmt.Println("mousedown in unknown state")
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.lastAct = time.Now()
		x, y := ebiten.CursorPosition()

		switch g.state {
		case inScreenSaver:
			g.leaveScreenSaver()
		case inNormal:
		for idx := range sprites {
			if sprites[idx].in(x, y) {
				fmt.Printf("in sprite release %s\n", sprites[idx].name)
				sprites[idx].do(&sprites[idx])
			}
		}

		for idx := range controls {
			if controls[idx].in(x, y) {
				fmt.Printf("in control release %s\n", controls[idx].label)
				controls[idx].do(g)
			}
		}
		case inRadio:
			for idx := range radiobuttons {
				if radiobuttons[idx].in(x, y) {
					fmt.Println("in radiobutton mouseup")
				}
			}
			for idx := range modalbuttons {
				if modalbuttons[idx].in(x, y) {
					fmt.Println("in modalbutton mouseup")
					modalbuttons[idx].modaldo(g)
				}
			}
		case inAlarmConfig:
			/* for idx := range radiobuttons {
				if radiobuttons[idx].in(x, y) {
					fmt.Println("in radiobutton mouseup")
				}
			} */
			for idx := range modalbuttons {
				if modalbuttons[idx].in(x, y) {
					fmt.Println("in modalbutton mouseup")
					modalbuttons[idx].modaldo(g)
				}
			}
		case inWeather:
			for idx := range modalbuttons {
				if modalbuttons[idx].in(x, y) {
					fmt.Println("in modalbutton mouseup")
					modalbuttons[idx].modaldo(g)
				}
			}
		default:
			fmt.Println("mouseup in unknown state")
		}

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		log.Fatal("shutting down")
	}

	return nil
}
