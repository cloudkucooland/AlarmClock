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
			sprite: getSprite("Spring", "Sleeptimer", sleepStopPlayer),
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

	boxwidth := 440
	boxheight := 150
	borderwidth := 20
	x := (screensize.X / 2) - (boxwidth / 2)
	y := 230
	ypadding := 16
	xpadding := 10

	// TODO: base this on sprite size not hardcoded values
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(boxwidth), float32(boxheight), modalgrey, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(boxwidth), float32(boxheight), float32(4), bordergrey, false)
	vector.StrokeRect(screen, float32(x+xpadding), float32(y+10), float32(boxwidth-borderwidth), float32(boxheight-borderwidth), float32(2), bordergrey, false)

	// move from box corner to initial location fo icons
	y = y + ypadding
	x = x + 2*xpadding

	if !g.radio.IsPlaying() {
		play := g.radiocontrols["Play"]
		play.setLocation(x, y)
		play.drawWithLabel(screen)
		x = x + 100
	}

	if g.radio.IsPlaying() {
		up := g.radiocontrols["VolUp"]
		up.scale = 1.0
		bounds := up.sprite.image.Bounds()
		up.setLocation(x, y)
		up.draw(screen)

		dn := g.radiocontrols["VolDn"]
		dn.scale = 1.0
		dn.setLocation(x, y+bounds.Max.Y+ypadding)
		dn.setLabel(fmt.Sprintf("%d", int(g.radio.Volume()*100.0)))
		dn.drawWithLabel(screen)

		x = x + 100

		pause := g.radiocontrols["Pause"]
		pause.setLocation(x, y)
		pause.drawWithLabel(screen)

		x = x + 100

		stop := g.radiocontrols["Stop"]
		stop.setLocation(x, y)
		stop.drawWithLabel(screen)

		x = x + 100

		if !g.inSleepCountdown {
			stop := g.radiocontrols["SleepCountdown"]
			stop.setLocation(x, y)
			stop.drawWithLabel(screen)
		}
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
