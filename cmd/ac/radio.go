package main

import (
	"fmt"
	"image"
	"image/color"
	"net/http"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type radiobutton struct {
	sprite   *sprite
	label    string
	labelimg *ebiten.Image
	labelloc image.Point
	url      string
}

func (r *radiobutton) getlabel() string {
	return r.label
}

func (r *radiobutton) setlabelimg(i *ebiten.Image) {
	r.labelimg = i
}

var radiobuttons = []radiobutton{
	{
		sprite: getSprite("Tea Time"),
		label:  "WRR",
		url:    "https://kera.streamguys1.com/wrrlive",
	},
	{
		sprite: getSprite("Swan Mommy"),
		label:  "KERA",
		url:    "https://kera.streamguys1.com/keralive",
	},
	{
		sprite: getSprite("Indignent"),
		label:  "BBC 6 Music",
		url:    "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/bbc_6music/bbc_6music.isml/bbc_6music-audio%3d96000.norewind.m3u8",
	},
	{
		sprite: getSprite("Pinwheel"),
		label:  "BBC 4",
		url:    "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/bbc_radio_fourfm/bbc_radio_fourfm.isml/bbc_radio_fourfm-audio%3d96000.norewind.m3u8",
	},
	{
		sprite: getSprite("Spring"),
		label:  "Radio 1 Dance",
		url:    "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/bbc_radio_one_dance/bbc_radio_one_dance.isml/bbc_radio_one_dance-audio%3d96000.norewind.m3u8",
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
		radiobuttons[idx].sprite.setLocation(x, y)
		radiobuttons[idx].sprite.setScale(spriteScale)
		radiobuttons[idx].sprite.draw(screen)

		if radiobuttons[idx].labelimg == nil {
			genlabel(&(radiobuttons[idx]), color.RGBA{0x33, 0x33, 0x33, 0xee}, controlfont)
		}

		if radiobuttons[idx].labelloc.X == 0 {
			b := radiobuttons[idx].sprite.image.Bounds()
			spritecenterx := int(float64(radiobuttons[idx].sprite.loc.X) + float64(b.Max.X)*spriteScale/2.0)
			lb := radiobuttons[idx].labelimg.Bounds()
			labelcenterx := lb.Max.X / 2
			radiobuttons[idx].labelloc.X = spritecenterx - labelcenterx
			radiobuttons[idx].labelloc.Y = radiobuttons[idx].sprite.loc.Y + int(float64(b.Max.Y)*spriteScale+4.0)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(radiobuttons[idx].labelloc.X), float64(radiobuttons[idx].labelloc.Y))
		screen.DrawImage(radiobuttons[idx].labelimg, op)

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
	if strings.HasSuffix(r.url, ".m3u8") || strings.HasSuffix(r.url, ".m3u") {
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
	if g.radio == nil {
		return
	}
	if !g.radio.IsPlaying() {
		g.radio.Pause()
	}
	fmt.Println("stopping current stream")
	if err := g.radio.Close(); err != nil {
		fmt.Println(err.Error())
		return
	}
	g.radio = nil
}
