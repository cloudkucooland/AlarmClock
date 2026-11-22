package ac

import (
	"bytes"
	"context"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

type externalProc struct {
	cancel context.CancelFunc
}

func (g *Game) playExternal(url string) {
	g.playExternalWithCallback(url, nil)
}

func (g *Game) playExternalWithCallback(target string, onComplete func()) {
	g.stopExternalPlayer()
	g.refreshExternalVolume()

	ctx, cancel := context.WithCancel(context.Background())
	proc := &externalProc{cancel: cancel}

	g.externalMu.Lock()
	g.externalAudio = proc
	g.externalMu.Unlock()

	go func() {
		completedNaturally := false
		defer func() {
			cancel()
			g.externalMu.Lock()
			if g.externalAudio == proc {
				g.externalAudio = nil
			}
			g.externalMu.Unlock()

			if completedNaturally && onComplete != nil {
				onComplete()
			}
		}()

		var out bytes.Buffer
		args := []string{"-nodisp", "-autoexit", "-ac", "2", "-loglevel", "error", "-vn", target}
		cmd := exec.CommandContext(ctx, "ffplay", args...)
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Start(); err != nil {
			g.debug(err.Error())
			return
		}

		err := cmd.Wait()
		if out.Len() > 0 {
			g.debug(out.String())
		}
		if err != nil {
			if ctx.Err() == nil {
				g.debug(err.Error())
			}
			return
		}
		if ctx.Err() == nil {
			completedNaturally = true
		}
	}()
}

func (g *Game) stopExternalPlayer() {
	g.externalMu.Lock()
	proc := g.externalAudio
	g.externalAudio = nil
	g.externalMu.Unlock()

	if proc != nil {
		proc.cancel()
	}
}

func (g *Game) hasExternalAudio() bool {
	g.externalMu.Lock()
	defer g.externalMu.Unlock()
	return g.externalAudio != nil
}

func volumeUpExternal(g *Game) {
	if err := exec.Command("wpctl", "set-volume", "-l", "1.0", "@DEFAULT_AUDIO_SINK@", "10%+").Run(); err != nil {
		g.debug(err.Error())
		return
	}
	g.refreshExternalVolume()
}

func volumeDnExternal(g *Game) {
	cur := g.getExternalVolume()
	if cur <= 5 {
		return
	}
	step := "10%-"
	if cur < 15 {
		step = "5%"
	}
	if err := exec.Command("wpctl", "set-volume", "@DEFAULT_AUDIO_SINK@", step).Run(); err != nil {
		g.debug(err.Error())
		return
	}
	g.refreshExternalVolume()
}

func (g *Game) refreshExternalVolume() {
	out, err := exec.Command("wpctl", "get-volume", "@DEFAULT_AUDIO_SINK@").Output()
	if err != nil {
		g.debug(err.Error())
		return
	}
	fields := strings.Fields(string(out))
	if len(fields) >= 2 {
		if v, err := strconv.ParseFloat(fields[1], 64); err == nil {
			g.externalMu.Lock()
			g.externalVol = uint8(math.Round(v * 100))
			g.externalMu.Unlock()
		}
	}
}

func (g *Game) getExternalVolume() uint8 {
	g.externalMu.Lock()
	defer g.externalMu.Unlock()
	return g.externalVol
}
