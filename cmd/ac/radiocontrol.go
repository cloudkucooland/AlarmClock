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
			sprite: getSprite("Spring", "Stop", stopPlayer),
		},
		"SleepCountdown": {
			sprite: getSprite("Tea Time", "Sleeptimer", sleepStopPlayer),
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
	// if an external audio program is running (ffplay) draw controls for that
	if g.externalAudio != nil {
		g.drawExternalControls(screen)
		return
	}

	// if no internal player is running, don't draw anything
	if g.audioPlayer == nil {
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
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(boxwidth), float32(boxheight), modalgrey, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(boxwidth), float32(boxheight), float32(4), bordergrey, false)
	vector.StrokeRect(screen, float32(x+xpadding), float32(y+10), float32(boxwidth-borderwidth), float32(boxheight-borderwidth), float32(2), bordergrey, false)

	// move from box corner to initial location of icons
	y = y + ypadding
	x = x + 2*xpadding

	if g.audioPlayer.IsPlaying() {
		up := g.radiocontrols["VolUp"]
		up.scale = 1.0
		bounds := up.sprite.image.Bounds()
		up.setLocation(x, y)
		up.draw(screen)

		dn := g.radiocontrols["VolDn"]
		dn.scale = 1.0
		dn.setLocation(x, y+bounds.Max.Y+ypadding)
		dn.setLabel(fmt.Sprintf("%d", int(g.audioPlayer.Volume()*100.0)))
		dn.drawWithLabel(screen)

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
	if g.externalAudio != nil {
		volumeUpExternal(g)
		tick(g)
		return
	}

	// internal audio player
	vol := g.audioPlayer.Volume()
	vol = math.Min(vol+0.10, 1.0)
	g.audioPlayer.SetVolume(vol)
	dn := g.radiocontrols["VolDn"]
	dn.setLabel(fmt.Sprintf("%d", int(vol*100.0)))
	tick(g)
}

func volumeDn(g *Game) {
	if g.externalAudio != nil {
		volumeDnExternal(g)
		tick(g)
		return
	}

	// internal audio player
	vol := g.audioPlayer.Volume()
	vol = math.Max(vol-0.10, 0.05)
	g.audioPlayer.SetVolume(vol)
	dn := g.radiocontrols["VolDn"]
	dn.setLabel(fmt.Sprintf("%d", int(vol*100.0)))
	tick(g)
}
