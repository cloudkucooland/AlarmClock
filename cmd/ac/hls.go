package main

import (
	"io"
	"fmt"

	"github.com/bluenviron/gohlslib/v2"
	"github.com/bluenviron/gohlslib/v2/pkg/codecs"
	// "github.com/bluenviron/mediacommon/pkg/formats/mpegts"

	"github.com/winlinvip/go-fdkaac/fdkaac"

	"github.com/hajimehoshi/ebiten/v2/audio/wav"	
)

func (g *Game) playhls(url string) {
	c := &gohlslib.Client{
		URI: "http://qthttp.apple.com.edgesuite.net/1010qwoeiuryfg/sl.m3u8",
		// URI: url,
	}

	r, w := io.Pipe()

	decoder := fdkaac.NewAacDecoder()
	asc := []byte{0x12, 0x10}
	if err := decoder.InitRaw(asc); err != nil {
		g.debug(err.Error())
		return
	}
	defer decoder.Close()
	
	c.OnTracks = func(tracks []*gohlslib.Track) error {
		track := findMPEG4AudioTrack(tracks)
		if track == nil {
			err := fmt.Errorf("no MPEG-4 audio track found")
			g.debug(err.Error())
			return err
		}

		c.OnDataMPEG4Audio(track, func(pts int64, aus [][]byte) {
			fmt.Printf("%+v\n", aus[0])
			var pcm []byte
			for i := range aus {
				if pcm, err = decoder.Decode(aus[i]); err != nil  {
					g.debug(err.Error())
					continue
				}
				w.Write(pcm)
			}
		})

		return nil
	}

	if err := c.Start(); err != nil {
		g.debug(err.Error())
		return
	}
	defer c.Close()


	// var err error
	g.audioPlayer, err = g.audioContext.NewPlayerFrom(r)
	if err != nil {
		g.debug(err.Error())
		return
	}
	g.audioPlayer.Play()
	defer stopPlayer(g)

	err = <-c.Wait()
	w.Close()
	// r.Close()
	g.debug(err.Error())
}

func findMPEG4AudioTrack(tracks []*gohlslib.Track) *gohlslib.Track {
	for _, track := range tracks {
		if _, ok := track.Codec.(*codecs.MPEG4Audio); ok {
			return track
		}
	}
	return nil
}
