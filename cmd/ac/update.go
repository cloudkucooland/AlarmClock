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
			g.alarmStateIcon.location()
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
			for _, c := range g.controls {
				if c.in(x, y) && !c.ani.in {
					c.startanimation()
				}
			}
			if g.audioPlayer != nil || g.externalAudio != nil {
				for _, r := range g.radiocontrols {
					if r.in(x, y) && !r.ani.in {
						r.startanimation()
					}
				}
			}
		case inRadio:
			for _, rb := range g.radiobuttons {
				if rb.in(x, y) && !rb.ani.in {
					rb.startanimation()
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
			for _, c := range g.controls {
				if c.in(x, y) {
					if c.ani.in {
						c.stopanimation()
					}
					c.do(g)
				}
			}

			if g.audioPlayer != nil || g.externalAudio != nil {
				for _, r := range g.radiocontrols {
					if r.in(x, y) {
						if r.ani.in {
							r.stopanimation()
						}
						r.do(g)
					}
				}
			}
		case inRadio:
			for _, rb := range g.radiobuttons {
				if rb.in(x, y) {
					rb.startPlayer(g)
				}
			}
			for _, m := range modalbuttons {
				if m.in(x, y) {
					m.do(g)
				}
			}
		case inAlarmConfig:
			processed := false
			for key, a := range g.config.Alarms {
				if a.dialogButton.hourUp.in(x, y) {
					a.dialogButton.hourUp.do(g)
					processed = true
				}
				if a.dialogButton.hourDn.in(x, y) {
					a.dialogButton.hourDn.do(g)
					processed = true
				}
				if a.dialogButton.minUp.in(x, y) {
					a.dialogButton.minUp.do(g)
					processed = true
				}
				if a.dialogButton.minDn.in(x, y) {
					a.dialogButton.minDn.do(g)
					processed = true
				}

				// check if we are in the body of the list-item, to enable/disable it
				if !processed && a.in(x, y) {
					if g.config.EnabledAlarmID == key {
						chirp(g)
						g.config.EnabledAlarmID = disabledAlarmID
						g.clearAlarmStateIcon()
					} else {
						chirp(g)
						g.config.EnabledAlarmID = key
						a.station = g.selectedStation
						g.setAlarmStateIcon()
					}
					processed = true
				}

				if processed {
					_ = g.storeconfig()
				}
			}
			for _, m := range modalbuttons {
				if m.in(x, y) {
					m.do(g)
				}
			}
		case inAlarm:
			for _, a := range g.alarmbuttons {
				if a.in(x, y) {
					g.debug("in alarmbutton mouseup")
					a.do(g)
				}
			}
		case inSnooze:
			g.debug("mouseup in snooze")
			if b, ok := g.alarmbuttons["Stop"]; ok {
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
