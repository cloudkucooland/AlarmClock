package main

import (
	"context"
	"fmt"
	// "io"
	"os/exec"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) playExternal(url string) {
	ctx, cancel := context.WithCancel(context.Background())
	session := &audioSession{cancel: cancel}
	g.externalAudio = session

	args := []string{"-ac", "1", "-loglevel", "error", "-nodisp", "-volume", "50", "-vn", url}
	cmd := exec.CommandContext(ctx, "ffplay", args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		g.debug(fmt.Sprintf("Failed to get stdin pipe: %v", err))
		g.externalAudio = nil
		cancel()
		return
	}
	g.externalAudioStdin = stdin

	if err := cmd.Start(); err != nil {
		g.debug(fmt.Sprintf("Failed to start ffplay: %v", err))
		g.externalAudio = nil
		cancel()
		return
	}

	// Goroutine to monitor the process and clean up state
	go func(s *audioSession) {
		if err := cmd.Wait(); err != nil {
			g.debug(fmt.Sprintf("ffplay exited: %v", err))
		}
		if g.externalAudio == s {
			g.externalAudio = nil
			g.externalAudioStdin = nil
		}
	}(session)
}

func (g *Game) stopExternalPlayer() {
	if g.externalAudio == nil {
		return
	}

	if g.externalAudioStdin != nil {
		g.externalAudioStdin.Close()
		g.externalAudioStdin = nil
	}

	g.externalAudio.cancel()
	g.externalAudio = nil
}

func (g *Game) drawExternalControls(screen *ebiten.Image) {
	if g.externalAudio == nil {
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
	dn.setLabel(fmt.Sprintf("%d", getExternalVolume()))
	dn.drawWithLabel(screen)

	x = x + 100
	stop := g.radiocontrols["Stop"]
	stop.setLocation(x, y)
	stop.drawWithLabel(screen)

	/* not yet
	x = x + 100
	if !g.inSleepCountdown {
		stop := g.radiocontrols["SleepCountdown"]
		stop.setLocation(x, y)
		stop.drawWithLabel(screen)
	} */
}

func volumeUpExternal(g *Game) {
	if g.externalAudioStdin != nil {
		g.externalAudioStdin.Write([]byte("9"))
	}
}

func volumeDnExternal(g *Game) {
	if g.externalAudioStdin != nil {
		g.externalAudioStdin.Write([]byte("0"))
	}
}

func getExternalVolume() uint8 {
	// ffplay doesn't provide easy volume feedback.
	// Returning a dummy value for the UI label.
	return 50
}
