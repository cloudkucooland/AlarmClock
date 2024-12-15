package main

import (
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
		g.clock.timestring = now.Format(g.config.ClockFormat)

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
		case inNormal, inSnooze:
			for idx := range controls {
				if controls[idx].in(x, y) && !controls[idx].ani.in {
					controls[idx].startanimation()
				}
			}
			if g.radio != nil {
				for idx := range radiocontrols {
					if radiocontrols[idx].in(x, y) && !radiocontrols[idx].ani.in {
						radiocontrols[idx].startanimation()
					}
				}
			}
		case inRadio:
			for idx := range radiobuttons {
				if radiobuttons[idx].in(x, y) && !radiobuttons[idx].ani.in {
					radiobuttons[idx].startanimation()
				}
			}
		case inAlarmConfig, inAlarm:
		default:
			g.debug("mousedown in unknown state")
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
				if controls[idx].in(x, y) {
					if controls[idx].ani.in {
						controls[idx].stopanimation()
					}
					controls[idx].do(g)
				}
			}

			if g.radio != nil {
				for idx := range radiocontrols {
					if radiocontrols[idx].in(x, y) {
						if radiocontrols[idx].ani.in {
							radiocontrols[idx].stopanimation()
						}
						radiocontrols[idx].do(g)
					}
				}
			}
		case inRadio:
			for idx := range radiobuttons {
				if radiobuttons[idx].in(x, y) {
					radiobuttons[idx].startPlayer(g)
				}
			}
			for idx := range modalbuttons {
				if modalbuttons[idx].in(x, y) {
					modalbuttons[idx].do(g)
				}
			}
		case inAlarmConfig:
			for key := range g.config.Alarms {
				if g.config.Alarms[key].in(x, y) {
					if g.config.EnabledAlarmID == key {
						g.config.EnabledAlarmID = disabledAlarmID
					} else {
						g.config.EnabledAlarmID = key
					}
					_ = g.storeconfig()
				}
			}
			for idx := range modalbuttons {
				if modalbuttons[idx].in(x, y) {
					modalbuttons[idx].do(g)
				}
			}
		case inAlarm:
			for idx := range alarmbuttons {
				if alarmbuttons[idx].in(x, y) {
					g.debug("in alarmbutton mouseup")
					alarmbuttons[idx].do(g)
				}
			}
		case inSnooze:
			g.debug("mouseup in snooze")
			if b, ok := alarmbuttons["Stop"]; ok {
				if b.in(x, y) {
					b.do(g)
				}
			}
		default:
			g.debug("mouseup in unknown state")
		}

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		log.Fatal("shutting down")
	}

	return nil
}
