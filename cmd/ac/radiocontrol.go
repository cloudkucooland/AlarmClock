package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
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

	boxwidth := 230
	if !g.inSleepCountdown {
		boxwidth += 105
	}

	boxheight := 142
	x := (screensize.X / 2) - (boxwidth / 2)
	y := 228
	padX := 22
	padY := 14

	// Frosted rounded playback dock matching the weather pill & modal style
	drawRoundedBox(screen, float32(x), float32(y), float32(boxwidth), float32(boxheight), 18, modalSlate, rimLight, 1.5)

	cx := x + padX
	cy := y + padY

	up := g.radiocontrols["VolUp"]
	up.scale = 1.0
	bounds := up.sprite.image.Bounds()
	up.setLocation(cx, cy)
	up.draw(screen)

	dn := g.radiocontrols["VolDn"]
	dn.scale = 1.0
	dn.setLocation(cx, cy+bounds.Max.Y+12)
	dn.setLabel(fmt.Sprintf("%d%%", g.currentVolumePercent()))
	dn.drawWithLabel(screen)

	cx += 96

	stop := g.radiocontrols["Stop"]
	stop.setLocation(cx, cy)
	stop.drawWithLabel(screen)

	cx += 104

	if !g.inSleepCountdown {
		sleep := g.radiocontrols["SleepCountdown"]
		sleep.setLocation(cx, cy)
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

