package main

import (
	"bufio"
	"fmt"
	"image/color"
	"net/http"
	"net/url"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/grafov/m3u8"
)

type radiobutton struct {
	sprite *sprite
	label  string
	url    string
}

var radiobuttons = []radiobutton{
	{
		sprite: getSprite("Indignent"),
		label:  "BBC 6Music",
		url:    "http://as-hls-ww-live.akamaized.net/pool_904/live/ww/bbc_6music/bbc_6music.isml/bbc_6music-audio%3d96000.norewind.m3u8",
		// url:    "http://lstn.lv/bbc.m3u8?station=bbc_6music&bitrate=96000",
	},
	{
		sprite: getSprite("Love"),
		label:  "WRR",
		url:    "https://kera.streamguys1.com/wrrlive",
	},
	{
		sprite: getSprite("Love"),
		label:  "KERA",
		url:    "https://kera.streamguys1.com/keralive",
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

	paddedspritesize := 40 // spritesize * 1.5? modal border + 10?
	x := paddedspritesize
	y := paddedspritesize

	for idx := range radiobuttons {
		radiobuttons[idx].sprite.setLocation(x, y)
		radiobuttons[idx].sprite.setScale(spriteScale)
		radiobuttons[idx].sprite.draw(screen)
		radiobuttons[idx].sprite.drawLabel(radiobuttons[idx].label, screen)

		// TODO: base this on sprite size not hardcoded values
		perrow := 5       // make dynamic
		rowspacing := 112 // make dynamic
		x = x + rowspacing
		if (idx % perrow) == (perrow - 1) {
			x = paddedspritesize
			y = y + rowspacing
		}
	}
}

// for the main screen
func (g *Game) drawRadioControls(screen *ebiten.Image) {
	if g.radio == nil {
		return
	}

	grey := color.RGBA{0xaa, 0xaa, 0xaa, 0x99}
	border := color.RGBA{0x66, 0x66, 0x66, 0x00}

	borderwidth := 20

	vector.DrawFilledRect(screen, float32(borderwidth), float32(240), float32(screensize.X-(140)), float32(160), grey, false)
	vector.StrokeRect(screen, float32(borderwidth), float32(240), float32(screensize.X-(140)), float32(160), float32(4), border, false)
	vector.StrokeRect(screen, float32(borderwidth)*1.5, float32(260), float32(screensize.X-(160)), float32(120), float32(2), border, false)
}

func (r radiobutton) startPlayer(g *Game) {
	if g.radio != nil {
		g.radio.Pause()
		fmt.Println("stopping current stream")
		if err := g.radio.Close(); err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	// if a playlist is requested, do that in a new goprocess
	if strings.HasSuffix(r.url, ".m3u8") {
		go g.playlist(r.url)
		return
	}
	g.fetchandplay(r.url)
}

func (g *Game) fetchandplay(url string) {
	fmt.Println("Starting stream", url)
	stream, err := http.Get(url)
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

// TODO: https://www.youtube.com/watch?v=uhOAt-aao64
func (g *Game) playlist(playlisturl string) {
	f, err := http.Get(playlisturl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	split := strings.Split(playlisturl, "/")
	base := strings.Join(split[:len(split)-1], "/")
	fmt.Printf("playlist: %s\nbase: %s\n", playlisturl, base)

	parsed, listType, err := m3u8.DecodeFrom(bufio.NewReader(f.Body), true)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch listType {
	case m3u8.MEDIA:
		mediapl := parsed.(*m3u8.MediaPlaylist)
		for _, segment := range mediapl.Segments {
			if segment != nil {
				mediaurl, err := abs(segment.URI, base)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				g.fetchandplayts(mediaurl)
			}
		}
	case m3u8.MASTER:
		fmt.Println("found master playlist: not implemented")
	}
}

func abs(mediaurl string, base string) (string, error) {
	parsedmediaurl, err := url.Parse(mediaurl)
	if err != nil {
		return "", err
	}

	if parsedmediaurl.IsAbs() {
		return parsedmediaurl.String(), nil
	}

	return fmt.Sprintf("%s/%s", base, mediaurl), nil

}

func (g *Game) fetchandplayts(url string) {
	// TODO dunno how to do these yet
}
