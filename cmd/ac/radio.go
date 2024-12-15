package main

import (
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type radiobutton struct {
	*sprite
	url   string
	works bool
}

var radiobuttons = map[string]*radiobutton{
	"WRR": {
		sprite: getSprite("Tea Time", "WRR", chirp),
		url:    "https://kera.streamguys1.com/wrrlive",
		works:  true,
	},
	"KERA": {
		sprite: getSprite("Swan Mommy", "KERA", chirp),
		url:    "https://kera.streamguys1.com/keralive",
		works:  true,
	},
	"90s90s Dance": {
		sprite: getSprite("Spring", "90s90s Dance", chirp),
		url:    "https://streams.90s90s.de/danceradio/mp3-192/",
		works:  true,
	},
	"90s90s Techno": {
		sprite: getSprite("Spring", "90s90s Techno", chirp),
		url:    "http://streams.90s90s.de/techno/mp3-192/",
		works:  true,
	},
	"Sunshine Live": {
		sprite: getSprite("Spring", "Sunshine Live", chirp),
		// url:    "http://stream.sunshine-live.de/techno/mp3-192/play.m3u",
		url:   "http://stream.sunshine-live.de/techno/mp3-192/",
		works: true,
	},
	"Dub Techno": {
		sprite: getSprite("Spring", "Dub Techno", chirp),
		url:    "http://94.130.113.214:8000/dubtechno",
		works:  true,
	},
	"Chillout": {
		sprite: getSprite("Spring", "Chillout", chirp),
		url:    "http://144.76.106.52:7000/chillout.mp3",
		works:  true,
	},
	"Sleeping Pill": {
		sprite: getSprite("Spring", "Sleeping Pill", chirp),
		url:    "http://radio.stereoscenic.com:80/asp-l.mp3",
		works:  true,
	},
	"Radio Frisky": {
		sprite: getSprite("Swan Mommy", "Radio Frisky", chirp),
		url:    "https://stream.friskyradio.com",
	},
	"BBC 6 Music": {
		sprite: getSprite("Indignent", "BBC 6 Music", chirp),
		// url:    "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/bbc_6music/bbc_6music.isml/bbc_6music-audio%3d96000.norewind.m3u8",
		url: "http://lstn.lv/bbc.m3u8?station=bbc_6music&bitrate=320000",
	},
	"BBC 4": {
		sprite: getSprite("Pinwheel", "BBC 4", chirp),
		url:    "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/bbc_radio_fourfm/bbc_radio_fourfm.isml/bbc_radio_fourfm-audio%3d96000.norewind.m3u8",
	},
	"Radio 1 Dance": {
		sprite: getSprite("Spring", "Radio 1 Dance", chirp),
		url:    "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/bbc_radio_one_dance/bbc_radio_one_dance.isml/bbc_radio_one_dance-audio%3d96000.norewind.m3u8",
	},
	"BBC World Service": {
		sprite: getSprite("Spring", "BBC World Service", chirp),
		// url:    "http://stream.live.vc.bbcmedia.co.uk/bbc_world_service",
		url: "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/audio_pop_up_01/audio_pop_up_01.isml/audio_pop_up_01-audio=96000.norewind.m3u8",
	},
}

func radioDialog(g *Game) {
	g.state = inRadio
}

func (g *Game) drawRadioDialog(screen *ebiten.Image) {
	g.drawModal(screen)

	// TODO: base this on sprite size not hardcoded values
	paddedspritesize := 60 // spritesize * 1.5? modal border + 10?
	x := paddedspritesize
	y := paddedspritesize
	rowspacing := 112 // make dynamic

	for _, k := range slices.Sorted(maps.Keys(radiobuttons)) {
		rb, ok := radiobuttons[k]
		if !ok || !rb.works {
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
	g.selectedStation = r

	r.stopPlayer(g)

	// if a playlist is requested, do that in a new goprocess
	if strings.Contains(r.url, "m3u") {
		g.debug("starting HLS player logic")
		go g.playhls(r.url)
		return
	}

	g.debug(r.url)
	stream, err := http.Get(r.url)
	if err != nil {
		g.debug(err.Error())
		return
	}

	decoded, err := mp3.DecodeWithSampleRate(44100, stream.Body)
	if err != nil {
		g.debug(err.Error())
		return
	}

	g.radio, err = g.audioContext.NewPlayer(decoded)
	if err != nil {
		g.debug(err.Error())
		return
	}
	g.radio.Play()
}

func (r *radiobutton) stopPlayer(g *Game) {
	stopPlayer(g)
}

func stopPlayer(g *Game) {
	if g.radio == nil {
		return
	}
	if g.radio.IsPlaying() {
		g.radio.Pause()
	}
	g.debug("stopping current stream")
	if err := g.radio.Close(); err != nil {
		g.debug(err.Error())
		return
	}
	g.radio = nil
}

func sleepStopPlayer(g *Game) {
	g.inSleepCountdown = true
	go func(g *Game) {
		g.debug("sleep countdown")
		time.Sleep(30 * time.Minute)
		g.debug("stopping player from sleep countdown")
		stopPlayer(g)
	}(g)
}

func pausePlayer(g *Game) {
	if g.radio == nil {
		return
	}
	if g.radio.IsPlaying() {
		g.radio.Pause()
	}
	g.debug("pausing current stream")
}

func resumePlayer(g *Game) {
	if g.radio == nil {
		return
	}
	if !g.radio.IsPlaying() {
		g.radio.Play()
	}
	g.debug("resuming current stream")
}

func defaultStation() *radiobutton {
	return radiobuttons["KERA"]
}
