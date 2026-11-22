package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
		"Next": {
			sprite: getSprite("Pinwheel", "Next", func(g *Game) {
				chirp(g)
				g.nextMusicTrack()
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

	hasMusic := g.hasActiveMusicQueue()

	boxwidth := 230
	if hasMusic {
		boxwidth += 104
	}
	if !g.inSleepCountdown {
		boxwidth += 104
	}

	boxheight := 142
	x := (screensize.X / 2) - (boxwidth / 2)
	y := 232
	padX := 22
	padY := 14

	// Now-playing pill above the playback dock when playing from the music library
	if hasMusic {
		if title := g.currentNowPlaying(); title != "" {
			banner := truncateLabel(title, 44)
			tw, th := text.Measure(banner, controlfont, 0)
			bw := float32(tw) + 28
			bh := float32(th) + 12
			bx := float32(screensize.X/2) - bw/2
			by := float32(y) - bh - 8
			drawRoundedBox(screen, bx, by, bw, bh, bh/2, modalSlate, rimLight, 1.0)
			drawShadowedText(screen, banner, controlfont, float64(bx)+14, float64(by)+6, textPrimary)
		}
	}

	// Frosted rounded playback dock
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

	if hasMusic {
		next := g.radiocontrols["Next"]
		next.setLocation(cx, cy)
		next.drawWithLabel(screen)
		cx += 104
	}

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
