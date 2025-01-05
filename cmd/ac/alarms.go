package main

import (
	"bytes"
	"image"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/cloudkucooland/AlarmClock/resources/sounds"
)

type alarmid int

const disabledAlarmID = -1

type Alarm struct {
	AlarmTime    AlarmTime
	triggered    bool
	snooze       bool
	snoozeCount  int
	station      *radiobutton
	dialogButton alarmDialogButton
}

type AlarmTime struct {
	Hour   int
	Minute int
}

func (g *Game) checkAlarms(hour int, minute int) {
	a, ok := g.config.Alarms[g.config.EnabledAlarmID]
	if !ok {
		return
	}

	if !a.snooze && a.AlarmTime.Hour == hour && a.AlarmTime.Minute == minute {
		g.startAlarm(g.config.EnabledAlarmID)
		return
	}

	if a.snooze {
		snoozehour := a.AlarmTime.Hour
		snoozemin := a.AlarmTime.Minute + (g.config.SnoozeDuration * a.snoozeCount)
		if snoozemin >= 60 {
			snoozemin = snoozemin - 60
			snoozehour = snoozehour + 1
		}
		if snoozehour == hour && snoozemin == minute {
			g.wakeFromSnooze(g.config.EnabledAlarmID)
		}
	}
}

func (g *Game) startAlarm(id alarmid) {
	g.lastAct = time.Now()
	if g.state == inScreenSaver {
		g.leaveScreenSaver()
	}

	a, ok := g.config.Alarms[id]
	if !ok {
		g.debug("cannot start unknown alarm?")
		return
	}

	g.state = inAlarm
	a.triggered = true
	g.startAlarmPlayer(a)
	g.audioPlayer.SetVolume(0.25)
}

func snoozeAlarm(g *Game) {
	g.lastAct = time.Now()

	a, ok := g.config.Alarms[g.config.EnabledAlarmID]
	if !ok {
		g.debug("alarm not enabled, bailing")
		g.state = inScreenSaver
		return
	}

	if a.snoozeCount >= 3 {
		g.debug("snoozed too many times... disable")
		// turn the volume up a click
		// play "nope" sound
		return
	}

	g.state = inSnooze
	// a.triggered = false
	a.snooze = true
	a.snoozeCount = a.snoozeCount + 1
	g.stopAlarmPlayer()
}

func stopAlarm(g *Game) {
	a, ok := g.config.Alarms[g.config.EnabledAlarmID]
	if !ok {
		g.debug("no alarm enabled, nothing to stop")
		g.stopAlarmPlayer()
		return
	}

	g.lastAct = time.Now()
	g.state = inNormal

	if !a.triggered {
		g.debug("enabled alarm not triggered, nothing to stop")
		// do it anyways?
		// return
	}
	g.stopAlarmPlayer()
	a.triggered = false
	a.snooze = false
	a.snoozeCount = 0
	g.config.EnabledAlarmID = disabledAlarmID
}

func (g *Game) wakeFromSnooze(id alarmid) {
	g.lastAct = time.Now()
	if g.state == inScreenSaver {
		g.leaveScreenSaver()
	}
	g.state = inAlarm

	a, ok := g.config.Alarms[id]
	if !ok {
		g.debug("unable to wake from snooze for unknown alarm")
		return
	}
	a.triggered = true
	g.startAlarmPlayer(a)
	// vol := float64((60 + (aa.snoozeCount * 10)) / 100)
	g.audioPlayer.SetVolume(0.50)
}

func (g *Game) drawAlarm(screen *ebiten.Image) {
	if stp, ok := g.alarmbuttons["Stop"]; ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(stp.loc.Min.X), float64(stp.loc.Min.Y))
		screen.DrawImage(stp.img, op)
	}

	if snz, ok := g.alarmbuttons["Snooze"]; ok {
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

func (g *Game) setupAlarmButtons() {
	padding := float64(10)

	g.alarmbuttons = make(map[string]*alarmbutton)

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
			do: snoozeAlarm,
		}
		g.alarmbuttons["Snooze"] = &q
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
			do: stopAlarm,
		}
		g.alarmbuttons["Stop"] = &q
	}
}

func (g *Game) drawSnooze(screen *ebiten.Image) {
	if stp, ok := g.alarmbuttons["Stop"]; ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(stp.loc.Min.X), float64(stp.loc.Min.Y))
		screen.DrawImage(stp.img, op)
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

func (g *Game) startAlarmPlayer(a *Alarm) {
	// use the station playing when the alarm was set
	r := a.station
	// if the alarm doesn't have a station set, use the last played station
	if r == nil {
		g.debug("alarm enabled w/o station set?")
		r = g.selectedStation
	}
	if r == nil {
		g.debug("alarm enabled, no station set, using default")
		r = g.defaultStation()
	}
	g.selectedStation = r
	r.startPlayer(g)

	// backup alarm if internet is down
	go func(g *Game) {
		time.Sleep(2 * time.Second)
		// if the radio isn't playing, but is triggered and not already snoozed
		// !triggered means the stop button has been pushed
		// snooze means the snooze button has been pushed
		a, ok := g.config.Alarms[g.config.EnabledAlarmID]
		if !ok {
			g.debug("no enabled alarms in backup check")
			return
		}
		if !g.audioPlayer.IsPlaying() && a.triggered && !a.snooze {
			g.stopAlarmPlayer()

			backup, ok := sounds.Sounds["BackupAlarm"]
			if !ok {
				// I guess we are sleeping in today...
				g.debug("no backup alarm found")
				return
			}
			backupAlarm, err := mp3.DecodeWithoutResampling(bytes.NewReader(backup))
			if err != nil {
				// I guess we are sleeping in today...
				g.debug(err.Error())
				return
			}

			loop := audio.NewInfiniteLoop(backupAlarm, backupAlarm.Length())
			loopplayer, err := g.audioContext.NewPlayer(loop)
			if err != nil {
				// I guess we are sleeping in today...
				g.debug(err.Error())
				return
			}
			loopplayer.SetVolume(0.33)
			loopplayer.Play()
		}
	}(g)
}

func (g *Game) stopAlarmPlayer() {
	stopPlayer(g)
}
