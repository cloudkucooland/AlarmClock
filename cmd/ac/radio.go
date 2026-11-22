package main

import (
	"context"
	"fmt"
	"io"
	"maps"
	"math"
	"net/http"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type stationName string

const defaultRadioStation = stationName(" KERA")

type radiobutton struct {
	*sprite
	url string
}

type readerOnly struct {
	io.Reader
}

func (g *Game) setupRadioButtons() {
	g.radiobuttons = map[stationName]*radiobutton{
		"WRR": {
			sprite: getSprite("Tea Time", "WRR", chirp),
			url:    "https://kera.streamguys1.com/wrrlive",
		},
		defaultRadioStation: {
			sprite: getSprite("Swan Mommy", "KERA", chirp),
			url:    "https://kera.streamguys1.com/keralive",
		},
		"90s90s Techno": {
			sprite: getSprite("Mad", "90s90s Techno", chirp),
			url:    "http://streams.90s90s.de/techno/mp3-192/",
		},
		"Dub Techno": {
			sprite: getSprite("Bathtime", "Dub Techno", chirp),
			url:    "http://ice1.somafm.com/groovesalad-128-mp3",
		},
		"Chillout": {
			sprite: getSprite("Baby", "Chillout", chirp),
			url:    "http://144.76.106.52:7000/chillout.mp3",
		},
		"Sleeping Pill": {
			sprite: getSprite("Artist", "Sleeping Pill", chirp),
			url:    "http://radio.stereoscenic.com/asp-s",
		},
		"a.m. ambient": {
			sprite: getSprite("Pinwheel", "a.m. ambient", chirp),
			url:    "http://radio.stereoscenic.com/ama-s",
		},
		"modern ambient": {
			sprite: getSprite("Indignent", "modern ambient", chirp),
			url:    "http://radio.stereoscenic.com/mod-s",
		},
		"SomaFM Drone Zone": {
			sprite: getSprite("Indignent", "Drone Zone", chirp),
			url:    "http://ice1.somafm.com/dronezone-128-mp3",
		},
		"Sleepbot Environmental": {
			sprite: getSprite("Pinwheel", "Sleepbot Environmental", chirp),
			url:    "http://ice1.somafm.com/sleepbot-128-mp3",
		},
		"Suburbs of Goa": {
			sprite: getSprite("Spring", "Suburbs of Goa", chirp),
			url:    "http://ice1.somafm.com/sf1033-128-mp3",
		},
	}
}

func radioDialog(g *Game) {
	g.state = inRadio
}

func (g *Game) drawRadioDialog(screen *ebiten.Image) {
	g.drawModal(screen)

	// TODO: base this on sprite size not hardcoded values
	paddedspritesize := 70 // spritesize * 1.5? modal border + 10?
	x := paddedspritesize
	y := paddedspritesize
	rowspacing := 142 // make dynamic

	for _, k := range slices.Sorted(maps.Keys(g.radiobuttons)) {
		rb, ok := g.radiobuttons[k]
		if !ok {
			continue
		}
		rb.setLocation(x, y)
		rb.drawWithLabel(screen)

		x = x + rowspacing
		if x > screensize.X-100 {
			x = paddedspritesize
			y = y + rowspacing
		}
	}
}

func (r *radiobutton) startPlayer(g *Game) {
	if r == nil {
		r = g.defaultStation()
	}

	g.selectedStation = r
	g.stopPlayer()

	if strings.Contains(r.url, "m3u") {
		g.playExternal(r.url)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	g.streamCancel = cancel

	go func(url string) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			g.debug(err.Error())
			chirp(g)
			return
		}

		client := &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: 8 * time.Second,
			},
		}
		req.Header.Set("User-Agent", "Birdhouse-Alarmclock/1.0")

		stream, err := client.Do(req)
		if err != nil {
			if ctx.Err() == nil {
				g.debug(err.Error())
				chirp(g)
			}
			return
		}

		if stream.StatusCode != http.StatusOK {
			stream.Body.Close()
			g.debug(fmt.Sprintf("stream HTTP %d: %s", stream.StatusCode, url))
			chirp(g)
			return
		}

		decoded, err := mp3.DecodeWithSampleRate(44100, stream.Body)
		if err != nil {
			stream.Body.Close()
			if ctx.Err() == nil {
				g.debug(err.Error())
				chirp(g)
			}
			return
		}

		if ctx.Err() != nil {
			stream.Body.Close()
			return
		}

		// player, err := g.audioContext.NewPlayer(decoded)
		player, err := g.audioContext.NewPlayer(readerOnly{decoded})
		if err != nil {
			stream.Body.Close()
			g.debug(err.Error())
			chirp(g)
			return
		}

		player.SetVolume(g.targetVol)
		g.audioPlayer = player
		g.streamBody = stream.Body
		g.audioPlayer.Play()
	}(r.url)
}

func (g *Game) stopPlayer() {
	g.inSleepCountdown = false
	if g.sleepCancel != nil {
		g.sleepCancel()
		g.sleepCancel = nil
	}

	if g.streamCancel != nil {
		g.streamCancel()
		g.streamCancel = nil
	}

	if g.hasExternalAudio() {
		g.stopExternalPlayer()
	}

	if g.audioPlayer != nil {
		if g.audioPlayer.IsPlaying() {
			g.audioPlayer.Pause()
		}
		if err := g.audioPlayer.Close(); err != nil {
			g.debug(err.Error())
		}
		g.audioPlayer = nil
	}

	if g.streamBody != nil {
		if err := g.streamBody.Close(); err != nil {
			g.debug(err.Error())
		}
		g.streamBody = nil
	}
}

func (g *Game) sleepStopPlayer() {
	if g.sleepCancel != nil {
		g.sleepCancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	g.sleepCancel = cancel
	g.inSleepCountdown = true

	// kick the screensaver on in 5
	g.lastAct = time.Now().Add(-(screensaverTimeout - (5 * time.Second)))

	go func(ctx context.Context) {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for i := 0; i < 6; i++ {
			if g.hasExternalAudio() {
				_ = exec.Command("wpctl", "set-volume", "@DEFAULT_AUDIO_SINK@", "5%-").Run()
				g.refreshExternalVolume()
			} else if g.audioPlayer != nil {
				vol := math.Max(g.audioPlayer.Volume()-0.05, 0.05)
				g.audioPlayer.SetVolume(vol)
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
		g.stopPlayer()
	}(ctx)
}

func (g *Game) defaultStation() *radiobutton {
	s, ok := g.radiobuttons[defaultRadioStation]
	if !ok {
		g.debug("unable to get default station, getting a random one")
		for name, rando := range g.radiobuttons {
			// whatever happens to come first, which is nondeterministic
			g.debug(fmt.Sprintf("Got %s", name))
			return rando
		}
	}
	return s
}

func (g *Game) isAudioPlaying() bool {
	if g.hasExternalAudio() {
		return true
	}
	if g.audioPlayer == nil {
		return false
	}

	return g.audioPlayer.IsPlaying()
}

func (g *Game) setStreamVolume(vol float64) {
	vol = math.Max(0.05, math.Min(1.0, vol))
	g.targetVol = vol

	if g.hasExternalAudio() {
		pct := fmt.Sprintf("%d%%", int(math.Round(vol*100)))
		if err := exec.Command("wpctl", "set-volume", "-l", "1.0", "@DEFAULT_AUDIO_SINK@", pct).Run(); err != nil {
			g.debug(err.Error())
		}
		g.refreshExternalVolume()
		return
	}

	if g.audioPlayer != nil {
		g.audioPlayer.SetVolume(vol)
	}
	if dn, ok := g.radiocontrols["VolDn"]; ok {
		dn.setLabel(fmt.Sprintf("%d", int(math.Round(vol*100))))
	}
}
