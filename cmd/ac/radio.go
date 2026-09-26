package main

import (
	"fmt"
	"maps"
	// "math"
	// "net/http"
	"slices"
	// "strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type stationName string

const defaultRadioStation = stationName(" KERA")

type radiobutton struct {
	*sprite
	url string
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
			url:    "http://94.130.113.214:8000/dubtechno",
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
		"BBC 6 Music": {
			sprite: getSprite("Indignent", "BBC 6 Music", chirp),
			url:    "http://as-hls-ww-live.akamaized.net/pool_81827798/live/ww/bbc_6music/bbc_6music.isml/bbc_6music-audio%3d96000.norewind.m3u8",
		},
		"BBC 4": {
			sprite: getSprite("Pinwheel", "BBC 4", chirp),
			url:    "http://as-hls-ww-live.akamaized.net/pool_55057080/live/ww/bbc_radio_fourfm/bbc_radio_fourfm.isml/bbc_radio_fourfm-audio%3d96000.norewind.m3u8",
		},
		"BBC World Service": {
			sprite: getSprite("Spring", "BBC World Service", chirp),
			url:    "http://as-hls-ww-live.akamaized.net/pool_87948813/live/ww/bbc_world_service/bbc_world_service.isml/bbc_world_service-audio%3d96000.norewind.m3u8",
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

	// All streams now routed through external player
	go g.playExternal(r.url)
}

func (g *Game) stopPlayer() {
	g.inSleepCountdown = false

	if g.externalAudio != nil {
		g.stopExternalPlayer()
	}

	if g.audioPlayer == nil {
		return
	}

	// do not use g.isAudioPlaying since external audio is covered above
	if g.audioPlayer.IsPlaying() {
		g.audioPlayer.Pause()
	}
	if err := g.audioPlayer.Close(); err != nil {
		g.debug(err.Error())
		return
	}
	g.audioPlayer = nil
}

func (g *Game) sleepStopPlayer() {
	g.inSleepCountdown = true

	// kick the screensaver on in 5
	g.lastAct = time.Now().Add(-(screensaverTimeout - (5 * time.Second)))

	go func(g *Game) {
		c := time.Tick(5 * time.Minute)

		i := 0
		for i < 6 {
			i = i + 1
			// Reduce volume
			volumeDnExternal(g)
			<-c
		}
		g.stopPlayer()
	}(g)
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
	if g.externalAudio != nil {
		return true
	}

	return g.audioPlayer.IsPlaying()
}
