package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	// update clock every 15 seconds
	g.clock.cyclesSinceTick = (g.clock.cyclesSinceTick + 1) % (15 * hz)
	if g.clock.cyclesSinceTick == 1 {
		g.clock.clearCache()
		now := time.Now()
		g.clock.timestring = now.Format("03:04")

		g.checkAlarms(now.Hour(), now.Minute())

		if g.state == inScreenSaver {
			g.clock.screensaverClockLocation()
		}
		if g.state == inNormal && now.After(g.lastAct.Add(2*time.Minute)) {
			g.startScreenSaver()
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
				if controls[idx].sprite.in(x, y) && !controls[idx].sprite.ani.in {
					controls[idx].sprite.startanimation()
				}
			}
		case inRadio, inAlarmConfig, inWeather, inAlarm, inSnooze:
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
			for idx := range controls {
				if controls[idx].sprite.in(x, y) {
					fmt.Printf("in control release %s\n", controls[idx].label)
					if controls[idx].sprite.ani.in {
						controls[idx].sprite.stopanimation()
					}
					controls[idx].do(g)
				}
			}
		case inRadio:
			for idx := range radiobuttons {
				if radiobuttons[idx].sprite.in(x, y) {
					fmt.Println("in radiobutton mouseup", radiobuttons[idx].label)
					radiobuttons[idx].startPlayer(g)
				}
			}
			for idx := range modalbuttons {
				if modalbuttons[idx].sprite.in(x, y) {
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
				if modalbuttons[idx].sprite.in(x, y) {
					fmt.Println("in modalbutton mouseup")
					modalbuttons[idx].modaldo(g)
				}
			}
		case inWeather:
			for idx := range modalbuttons {
				if modalbuttons[idx].sprite.in(x, y) {
					fmt.Println("in modalbutton mouseup")
					modalbuttons[idx].modaldo(g)
				}
			}
		case inAlarm:
			for idx := range alarmbuttons {
				if alarmbuttons[idx].in(x, y) {
					fmt.Println("in alarmbutton mouseup")
					alarmbuttons[idx].do(g)
				}
			}
		case inSnooze:
			if b, ok := alarmbuttons["Stop"]; ok {
				if b.in(x, y) {
					b.do(g)
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
