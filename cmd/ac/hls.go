package main

import (
	"context"
	// "fmt"
	//"os"
	"os/exec"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ffplay -ac 1 -loglevel error -vn URL

func (g *Game) playhls(url string) {
	ctx, cancel := context.WithCancel(context.Background())

	args := []string{"-ac", "1", "-loglevel", "error", "-vn", url}
	cmd := exec.CommandContext(ctx, "ffplay", args...)

	g.debug("starting ffplay")
	if err := cmd.Start(); err != nil {
		g.debug(err.Error())
	}

	g.externalAudioCancel = cancel
	cmd.Wait()
	g.debug("ffplay exited")
	g.externalAudioCancel = nil
}

func (g *Game) stopPlayerHLS() {
	if g.externalAudioCancel == nil {
		return
	}

	g.externalAudioCancel()
	g.externalAudioCancel = nil
}

func (g *Game) hlsDrawRadioControls(screen *ebiten.Image) {
	if g.externalAudioCancel == nil {
		return
	}

	boxwidth := 340
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

	if g.audioPlayer.IsPlaying() {
		/*
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
		*/
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

func hlsVolumeUp(g *Game) {}

func hlsVolumeDn(g *Game) {}
