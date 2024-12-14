package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type radiocontrol struct {
	*sprite
}

var radiocontrols = []radiocontrol{
	{
		sprite: getSprite("Tea Time", "Pause", pausePlayer),
	},
	{
		sprite: getSprite("Swan Mommy", "Play", resumePlayer),
	},
	{
		sprite: getSprite("Spring", "Stop", stopPlayer),
	},
	{
		sprite: getSprite("Spring", "Stop in 30 min", sleepStopPlayer),
	},
}

func (g *Game) drawRadioControls(screen *ebiten.Image) {
	if g.radio == nil {
		return
	}

	boxwidth := 160
	boxheight := 120
	borderwidth := 20

	// TODO: base this on sprite size not hardcoded values
	vector.DrawFilledRect(screen, float32(borderwidth), float32(240), float32(screensize.X-(boxwidth)), float32(boxheight), modalgrey, false)
	vector.StrokeRect(screen, float32(borderwidth), float32(240), float32(screensize.X-(boxwidth)), float32(boxheight), float32(4), bordergrey, false)
	vector.StrokeRect(screen, float32(borderwidth)*1.5, float32(250), float32(screensize.X-(boxwidth+borderwidth)), float32(boxheight-borderwidth), float32(2), bordergrey, false)

	x := (screensize.X - boxwidth - (32 * spriteScale)) / 2

	for idx := range radiocontrols {
		if radiocontrols[idx].label == "Play" && g.radio.IsPlaying() {
			continue
		}
		if radiocontrols[idx].label == "Pause" && !g.radio.IsPlaying() {
			continue
		}

		radiocontrols[idx].setLocation(x, 270)
		radiocontrols[idx].drawWithLabel(screen)

		x = x + 100
	}
}
