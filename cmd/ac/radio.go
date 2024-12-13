package main

import (
	"fmt"
	// "image/color"
	"net/http"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type radiobutton struct {
	*sprite
	url string
}

var radiobuttons = []radiobutton{
	{
		sprite: getSprite("Tea Time", "WRR", chirp),
		url:    "https://kera.streamguys1.com/wrrlive",
	},
	{
		sprite: getSprite("Swan Mommy", "KERA", chirp),
		url:    "https://kera.streamguys1.com/keralive",
	},
	{
		sprite: getSprite("Spring", "90s90s Dance", chirp),
		url:    "https://streams.90s90s.de/danceradio/mp3-192/",
	},
	{
		sprite: getSprite("Spring", "90s90s Techno", chirp),
		url:    "http://streams.90s90s.de/techno/mp3-192/",
	},
	{
		sprite: getSprite("Spring", "Sunshine Live", chirp),
		url:    "http://stream.sunshine-live.de/techno/mp3-192/play.m3u",
	},
	{
		sprite: getSprite("Spring", "Dub Techno", chirp),
		url:    "http://94.130.113.214:8000/dubtechno",
	},
	{
		sprite: getSprite("Spring", "Chillout", chirp),
		url:    "http://144.76.106.52:7000/chillout.mp3",
	},
	{
		sprite: getSprite("Spring", "Ambient Sleeping Pill", chirp),
		url:    "http://radio.stereoscenic.com:80/asp-l.mp3",
	},
	{
		sprite: getSprite("Swan Mommy", "Radio Frisky", chirp),
		url:    "https://stream.friskyradio.com",
	},
	{
		sprite: getSprite("Indignent", "BBC 6 Music", chirp),
		// url:    "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/bbc_6music/bbc_6music.isml/bbc_6music-audio%3d96000.norewind.m3u8",
		url: "http://lstn.lv/bbc.m3u8?station=bbc_6music&bitrate=320000",
	},
	{
		sprite: getSprite("Pinwheel", "BBC 4", chirp),
		url:    "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/bbc_radio_fourfm/bbc_radio_fourfm.isml/bbc_radio_fourfm-audio%3d96000.norewind.m3u8",
	},
	{
		sprite: getSprite("Spring", "Radio 1 Dance", chirp),
		url:    "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/bbc_radio_one_dance/bbc_radio_one_dance.isml/bbc_radio_one_dance-audio%3d96000.norewind.m3u8",
	},
	{
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

	for idx := range radiobuttons {
		radiobuttons[idx].setLocation(x, y)
		radiobuttons[idx].drawWithLabel(screen)

		x = x + rowspacing
		if x > screensize.X-60 {
			x = paddedspritesize
			y = y + rowspacing
		}
	}
}

func (r radiobutton) startPlayer(g *Game) {
	r.stopPlayer(g)

	// if a playlist is requested, do that in a new goprocess
	if strings.Contains(r.url, "m3u") {
		go g.playhls(r.url)
		return
	}

	fmt.Println("Starting stream", r.url)
	stream, err := http.Get(r.url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	decoded, err := mp3.DecodeWithSampleRate(44100, stream.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	g.radio, err = g.audioContext.NewPlayer(decoded)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	g.radio.Play()
}

func (r radiobutton) stopPlayer(g *Game) {
	stopPlayer(g)
}

func stopPlayer(g *Game) {
	if g.radio == nil {
		return
	}
	if g.radio.IsPlaying() {
		g.radio.Pause()
	}
	fmt.Println("stopping current stream")
	if err := g.radio.Close(); err != nil {
		fmt.Println(err.Error())
		return
	}
	g.radio = nil
}

func pausePlayer(g *Game) {
	if g.radio == nil {
		return
	}
	if g.radio.IsPlaying() {
		g.radio.Pause()
	}
	fmt.Println("pausing current stream")
}

func resumePlayer(g *Game) {
	if g.radio == nil {
		return
	}
	if !g.radio.IsPlaying() {
		g.radio.Play()
	}
	fmt.Println("resuming current stream")
}
