package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type radiocontrol struct {
	*sprite
}

func (g *Game) setupRadioControls() {
	g.radiocontrols = map[string]*radiocontrol{
		"Pause": {
			sprite: getSprite("Tea Time", "Pause", pausePlayer),
		},
		"Play": {
			sprite: getSprite("Swan Mommy", "Play", resumePlayer),
		},
		"Stop": {
			sprite: getSprite("Spring", "Stop", stopPlayer),
		},
		"SleepCountdown": {
			sprite: getSprite("Spring", "Stop in 30 min", sleepStopPlayer),
		},
		"VolUp": {
			sprite: getSprite("Up", "", volumeUp),
		},
		"VolDn": {
			sprite: getSprite("Dn", "", volumeDn),
		},
	}
}

func (g *Game) drawRadioControls(screen *ebiten.Image) {
	if g.radio == nil {
		return
	}

	boxwidth := 160
	boxheight := 120
	borderwidth := 20
	y := 270

	// TODO: base this on sprite size not hardcoded values
	vector.DrawFilledRect(screen, float32(borderwidth), float32(240), float32(screensize.X-(boxwidth)), float32(boxheight), modalgrey, false)
	vector.StrokeRect(screen, float32(borderwidth), float32(240), float32(screensize.X-(boxwidth)), float32(boxheight), float32(4), bordergrey, false)
	vector.StrokeRect(screen, float32(borderwidth)*1.5, float32(250), float32(screensize.X-(boxwidth+borderwidth)), float32(boxheight-borderwidth), float32(2), bordergrey, false)

	x := borderwidth * 2

	if g.radio.IsPlaying() {
		up := g.radiocontrols["VolUp"]
		up.setLocation(x, y-2) // make dynamic
		up.draw(screen)
		dn := g.radiocontrols["VolDn"]
		dn.setLocation(x, y+24) // make dynamic
		dn.setLabel(fmt.Sprintf("%d", int(g.radio.Volume()*100.0)))
		dn.drawWithLabel(screen)
		x = x + 100
	}
	if !g.radio.IsPlaying() {
		play := g.radiocontrols["Play"]
		play.setLocation(x, y)
		play.drawWithLabel(screen)
		x = x + 100
	}
	if g.radio.IsPlaying() {
		pause := g.radiocontrols["Pause"]
		pause.setLocation(x, y)
		pause.drawWithLabel(screen)
		x = x + 100
	}
	if g.radio.IsPlaying() {
		stop := g.radiocontrols["Stop"]
		stop.setLocation(x, y)
		stop.drawWithLabel(screen)
		x = x + 100
	}
	if !g.inSleepCountdown {
		stop := g.radiocontrols["SleepCountdown"]
		stop.setLocation(x, y)
		stop.drawWithLabel(screen)
		// x = x + 100
	}
}

func volumeUp(g *Game) {
	vol := g.radio.Volume()
	if vol > 0.89 {
		return
	}
	vol = (vol + 0.10)
	g.radio.SetVolume(vol)
	dn := g.radiocontrols["VolDn"]
	dn.setLabel(fmt.Sprintf("%d", int(vol*100.0)))
}

func volumeDn(g *Game) {
	vol := g.radio.Volume()
	if vol < 0.11 {
		return
	}
	vol = (vol - 0.10)
	g.radio.SetVolume(vol)
	dn := g.radiocontrols["VolDn"]
	dn.setLabel(fmt.Sprintf("%d", int(vol*100.0)))
}
