package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type radiocontrol struct {
	*sprite
}

func (g *Game) setupRadioControls() {
	g.radiocontrols = map[string]*radiocontrol{
		"Stop": {
			sprite: getSprite("Spring", "Stop", func(g *Game) {
				g.stopPlayer()
			}),
		},
		"SleepCountdown": {
			sprite: getSprite("Tea Time", "Sleeptimer", func(g *Game) {
				g.sleepStopPlayer()
			}),
		},
		"VolUp": {
			sprite: getSprite("Up", "", func(g *Game) {
				g.volumeUp()
			}),
		},
		"VolDn": {
			sprite: getSprite("Dn", "", func(g *Game) {
				g.volumeDn()
			}),
		},
	}
}

func (g *Game) drawRadioControls(screen *ebiten.Image) {
	if !g.isAudioPlaying() {
		return
	}

	boxwidth := 240
	if !g.inSleepCountdown {
		boxwidth += 100
	}

	boxheight := 150
	borderwidth := 20
	x := (screensize.X / 2) - (boxwidth / 2)
	y := 230
	ypadding := 16
	xpadding := 10

	// TODO: base this on sprite size not hardcoded values
	vector.FillRect(screen, float32(x), float32(y), float32(boxwidth), float32(boxheight), modalgrey, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(boxwidth), float32(boxheight), float32(4), bordergrey, false)
	vector.StrokeRect(screen, float32(x+xpadding), float32(y+10), float32(boxwidth-borderwidth), float32(boxheight-borderwidth), float32(2), bordergrey, false)

	// move from box corner to initial location of icons
	y = y + ypadding
	x = x + 2*xpadding

	up := g.radiocontrols["VolUp"]
	up.scale = 1.0
	bounds := up.sprite.image.Bounds()
	up.setLocation(x, y)
	up.draw(screen)

	dn := g.radiocontrols["VolDn"]
	dn.scale = 1.0
	dn.setLocation(x, y+bounds.Max.Y+ypadding)
	dn.setLabel(fmt.Sprintf("%d", g.currentVolumePercent()))
	dn.drawWithLabel(screen)

	x = x + 100

	stop := g.radiocontrols["Stop"]
	stop.setLocation(x, y)
	stop.drawWithLabel(screen)

	x = x + 100

	if !g.inSleepCountdown {
		sleep := g.radiocontrols["SleepCountdown"]
		sleep.setLocation(x, y)
		sleep.drawWithLabel(screen)
	}
}

func (g *Game) volumeUp() {
	if g.hasExternalAudio() {
		volumeUpExternal(g)
		tick(g)
		return
	}

	g.setStreamVolume(g.targetVol + 0.10)
	tick(g)
}

func (g *Game) volumeDn() {
	if g.hasExternalAudio() {
		volumeDnExternal(g)
		tick(g)
		return
	}

	g.setStreamVolume(g.targetVol - 0.10)
	tick(g)
}

func (g *Game) currentVolumePercent() int {
	if g.hasExternalAudio() {
		return int(g.getExternalVolume())
	}
	return int(math.Round(g.targetVol * 100))
}
