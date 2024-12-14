package main

import (
	"fmt"
	"image"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
)

type alarmid int

var enabledAlarmID alarmid = alarmid(-1)

// station     *radiobutton
const snoozeduration = 9

type alarm struct {
	alarmTime    alarmTime
	triggered    bool
	snooze       bool
	snoozeCount  int
	dialogButton image.Rectangle
}

var alarms = map[alarmid]*alarm{
	0: {alarmTime: alarmTime{7, 00}},
	1: {alarmTime: alarmTime{8, 00}},
	2: {alarmTime: alarmTime{5, 00}},
	3: {alarmTime: alarmTime{6, 30}},
	4: {alarmTime: alarmTime{4, 30}},
}

type alarmTime struct {
	hour   int
	minute int
}

func (g *Game) checkAlarms(hour int, minute int) {
	a, ok := alarms[enabledAlarmID]
	if !ok {
		// fmt.Println("no alarm enabled")
		return
	}

	if !a.snooze && a.alarmTime.hour == hour && a.alarmTime.minute == minute {
		g.startAlarm(enabledAlarmID)
		return
	}

	if a.snooze {
		snoozehour := a.alarmTime.hour
		snoozemin := a.alarmTime.minute + (snoozeduration * a.snoozeCount)
		if snoozemin >= 60 {
			snoozemin = snoozemin - 60
			snoozehour = snoozehour + 1
		}
		fmt.Println("snoozing until %d:%d (%d)", snoozehour, snoozemin, a.snoozeCount)
		if snoozehour == hour && snoozemin == minute {
			g.wakeFromSnooze(enabledAlarmID)
		}
	}
}

func (g *Game) startAlarm(a alarmid) {
	g.lastAct = time.Now()
	if g.state == inScreenSaver {
		g.leaveScreenSaver()
	}
	g.state = inAlarm

	aa, ok := alarms[a]
	if !ok {
		fmt.Println("cannot start unknown alarm?")
		return
	}
	aa.triggered = true
	fmt.Println("Starting", aa)

	g.startAlarmPlayer()
}

func snooze(g *Game) {
	fmt.Println("snoozing")

	g.lastAct = time.Now()

	a, ok := alarms[enabledAlarmID]
	if !ok {
		fmt.Println("no alarms enabled, bailing")
		// no alarms enabled
		// stop playing?
		g.state = inScreenSaver
		return
	}

	if a.snoozeCount >= 3 {
		fmt.Println("snoozed too many times... disable")
		// turn the volume up a click
		// play "nope" sound
		return
	}

	a.snoozeCount = a.snoozeCount + 1
	g.state = inSnooze

	a.triggered = false
	a.snooze = true
	a.snoozeCount = a.snoozeCount + 1
	fmt.Println("snoozing", a)

	g.stopAlarmPlayer()
}

func stop(g *Game) {
	fmt.Println("stopping triggered alarms")

	a, ok := alarms[enabledAlarmID]
	if !ok {
		fmt.Println("alarm not enabled, nothing to stop")
		return
	}

	g.lastAct = time.Now()
	g.state = inNormal

	if !a.triggered {
		fmt.Println("enabled alarm not triggerd, nothing to stop")
		return
	}
	g.stopAlarmPlayer()
	a.triggered = false
	a.snooze = false
	a.snoozeCount = 0
}

func (g *Game) wakeFromSnooze(a alarmid) {
	g.lastAct = time.Now()
	if g.state == inScreenSaver {
		g.leaveScreenSaver()
	}
	g.state = inAlarm

	aa, ok := alarms[a]
	if !ok {
		fmt.Println("unable to wake from snooze for unknown alarm")
		return
	}
	aa.triggered = true
	fmt.Println("unsnoozing", aa)

	g.startAlarmPlayer()
}

func (g *Game) drawAlarm(screen *ebiten.Image) {
	if len(alarmbuttons) == 0 {
		setupAlarmButtons()
	}

	if slp, ok := alarmbuttons["Stop"]; ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(slp.loc.Min.X), float64(slp.loc.Min.Y))
		screen.DrawImage(slp.img, op)
	}

	if snz, ok := alarmbuttons["Snooze"]; ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(snz.loc.Min.X), float64(snz.loc.Min.Y))
		screen.DrawImage(snz.img, op)
	}

	// draw clock in middle
	textwidth, textheight := text.Measure(g.clock.timestring, clockfont, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(screensize.X/2)-float64(textwidth/2), float64(screensize.Y/2)-float64(textheight/2))
	text.Draw(screen, g.clock.timestring, clockfont, op)
}

type alarmbutton struct {
	img *ebiten.Image
	loc image.Rectangle
	do  func(g *Game)
}

func (a alarmbutton) in(x int, y int) bool {
	return (x >= a.loc.Min.X && x <= a.loc.Max.X) && (y >= a.loc.Min.Y && y <= a.loc.Max.Y)
}

var alarmbuttons map[string]alarmbutton

func setupAlarmButtons() {
	padding := float64(10)

	alarmbuttons = make(map[string]alarmbutton)

	{
		btn := button("SNOOZE", green, bigbuttonfont)
		btnsize := btn.Bounds()
		x := int(math.Ceil(float64(screensize.X/2) - float64(btnsize.Max.X/2))) // centered
		y := int(math.Ceil(float64(screensize.Y-btnsize.Max.Y) - padding))      // at bottom

		q := alarmbutton{
			img: btn,
			loc: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + btnsize.Max.X, Y: y + btnsize.Max.Y},
			},
			do: snooze,
		}
		alarmbuttons["Snooze"] = q
	}

	{
		btn := button("STOP", pink, bigbuttonfont)
		btnsize := btn.Bounds()
		x := int(math.Ceil(float64(screensize.X/2) - float64(btnsize.Max.X/2))) // centered
		y := int(padding)                                                       // at top

		q := alarmbutton{
			img: btn,
			loc: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + btnsize.Max.X, Y: y + btnsize.Max.Y},
			},
			do: stop,
		}
		alarmbuttons["Stop"] = q
	}
}

func (g *Game) drawSnooze(screen *ebiten.Image) {
	if slp, ok := alarmbuttons["Stop"]; ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(slp.loc.Min.X), float64(slp.loc.Min.Y))
		screen.DrawImage(slp.img, op)
	}

	if g.clock.cache == nil {
		g.clock.cache = ebiten.NewImage(screensize.X, screensize.Y)
		textwidth, textheight := text.Measure(g.clock.timestring, clockfont, 0)
		op := &text.DrawOptions{}
		op.ColorScale.ScaleWithColor(green)
		op.GeoM.Translate(float64(screensize.X/2)-float64(textwidth/2), float64(screensize.Y/2)-float64(textheight/2))
		text.Draw(g.clock.cache, g.clock.timestring, clockfont, op)
	}
	screen.DrawImage(g.clock.cache, &ebiten.DrawImageOptions{})
}

func (g *Game) startAlarmPlayer() {
	fmt.Println("starting alarm player")
}

func (g *Game) stopAlarmPlayer() {
	fmt.Println("stopping alarm player")
}
