package main

import (
	"fmt"
	"image"
	// "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
)

type alarm struct {
	alarmTime   alarmTime
	enabled     bool
	station     *radiobutton
	triggered   bool
	snooze      bool
	snoozeCount int
}

const snoozeduration = 9

var alarms = []alarm{
	{
		alarmTime:   alarmTime{15, 00}, // when to wake up
		enabled:     true,
		station:     &radiobuttons[3],
		triggered:   false,
		snooze:      false,
		snoozeCount: 0,
	},
}

type alarmTime struct {
	hour   int
	minute int
}

func (g *Game) checkAlarms(hour int, minute int) {
	for idx, a := range alarms {
		if a.enabled && !a.snooze && a.alarmTime.hour == hour && a.alarmTime.minute == minute {
			g.startAlarm(idx)
		}
		if a.enabled && a.snooze {
			snoozehour := a.alarmTime.hour
			snoozemin := a.alarmTime.minute + (snoozeduration * a.snoozeCount)
			if snoozemin >= 60 { // XXX assumes you don't mash on it more than a dozen times
				snoozemin = snoozemin - 60
				snoozehour = snoozehour + 1
			}
			if snoozehour == hour && snoozemin == minute {
				g.wakeFromSnooze(idx)
			}
		}
	}
}

func (g *Game) startAlarm(alarmID int) {
	alarm := alarms[alarmID]
	alarm.triggered = true
	g.state = inAlarm
	fmt.Println("Starting", alarms[alarmID])
	// start playing radio
}

func (g *Game) snooze(alarmID int) {
	alarm := alarms[alarmID]
	alarm.snooze = true
	alarm.snoozeCount = alarm.snoozeCount + 1
	// alarm.snoozeUntil =
	panic("finish snoozeuntil logic")
}

func (g *Game) wakeFromSnooze(alarmID int) {
	alarm := alarms[alarmID]
	alarm.triggered = true
	g.state = inAlarm
	fmt.Println("unsnoozing", alarms[alarmID])
	// start playing radio
}

func (g *Game) drawAlarm(screen *ebiten.Image) {
	// draw stop / snooze buttons
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
