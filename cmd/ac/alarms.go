package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
		alarmTime:   alarmTime{21, 05}, // when to wake up
		enabled:     true,
		station:     &radiobuttons[1],
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
	g.lastAct = time.Now()
	if g.state == inScreenSaver {
		g.leaveScreenSaver()
	}
	g.state = inAlarm

	alarms[alarmID].triggered = true
	fmt.Println("Starting", alarms[alarmID])

	alarms[alarmID].station.startPlayer(g)
}

func snooze(g *Game) {
	fmt.Println("snoozing")

	g.lastAct = time.Now()

	alarmID := -1
	for idx := range alarms {
		if alarms[idx].triggered {
			alarmID = idx
			break
		}
	}
	if alarmID == -1 {
		// no alarms in triggered state?
		fmt.Println("no alarms triggered, bailing")
		// stop playing?
		g.state = inScreenSaver
		return
	}

	if alarms[alarmID].snoozeCount >= 3 {
		fmt.Println("snoozed too many times... disable")
		// turn the volume up a click
		// play "nope" sound
		return
	}

	alarms[alarmID].snoozeCount = alarms[alarmID].snoozeCount + 1
	g.state = inSnooze

	alarms[alarmID].triggered = false
	alarms[alarmID].snooze = true
	alarms[alarmID].snoozeCount = alarms[alarmID].snoozeCount + 1
	fmt.Println("snoozing", alarms[alarmID])

	alarms[alarmID].station.stopPlayer(g)
}

func stop(g *Game) {
	fmt.Println("stopping triggered alarms")
	g.lastAct = time.Now()
	g.state = inNormal

	// this is almost certainly wrong... brute force all triggered alarms to off
	for idx := range alarms {
		if alarms[idx].triggered {
			alarms[idx].station.stopPlayer(g)
			alarms[idx].triggered = false
			alarms[idx].snooze = false
			alarms[idx].snoozeCount = 0
		}
	}
}

func (g *Game) wakeFromSnooze(alarmID int) {
	g.lastAct = time.Now()
	if g.state == inScreenSaver {
		g.leaveScreenSaver()
	}
	g.state = inAlarm

	alarms[alarmID].triggered = true
	fmt.Println("unsnoozing", alarms[alarmID])
	alarms[alarmID].station.startPlayer(g)
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
	pink := color.RGBA{0xff, 0x33, 0x77, 0xff}
	green := color.RGBA{0x33, 0xff, 0x33, 0xff}
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

func alarmConfigDialog(g *Game) {
	g.state = inAlarmConfig
}

func (g *Game) drawAlarmConfig(screen *ebiten.Image) {
	g.drawModal(screen)
}

func (g *Game) drawSnooze(screen *ebiten.Image) {
	green := color.RGBA{0x33, 0xff, 0x33, 0xff}

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
