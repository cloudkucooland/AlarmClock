package main

import (
	"bytes"
	"fmt"

	"github.com/bluenviron/gohlslib/v2"
	"github.com/bluenviron/gohlslib/v2/pkg/codecs"
)

func (g *Game) playhls(url string) {
	c := &gohlslib.Client{
		URI: url,
	}

	var b bytes.Buffer

	// called when tracks are parsed
	c.OnTracks = func(tracks []*gohlslib.Track) error {
		for _, track := range tracks {
			ttrack := track

			g.debug(fmt.Sprintf("detected track with codec %T\n", track.Codec))

			// set a callback that is called when data is received
			switch track.Codec.(type) {
			case *codecs.AV1:
				c.OnDataAV1(track, func(pts int64, tu [][]byte) {
					fmt.Printf("received data from track %T, pts = %v\n", ttrack.Codec, pts)
					for i := range tu {
						b.Write(tu[i])
					}
				})

			case *codecs.H264, *codecs.H265:
				c.OnDataH26x(track, func(pts int64, dts int64, au [][]byte) {
					fmt.Printf("received data from track %T, pts = %v\n", ttrack.Codec, pts)
					for i := range au {
						b.Write(au[i])
					}
				})

			case *codecs.MPEG4Audio:
				c.OnDataMPEG4Audio(track, func(pts int64, aus [][]byte) {
					fmt.Printf("received data from track %T, pts = %v\n", ttrack.Codec, pts)
					for i := range aus {
						b.Write(aus[i])
					}
				})

			case *codecs.Opus:
				c.OnDataOpus(track, func(pts int64, packets [][]byte) {
					fmt.Printf("received data from track %T, pts = %v\n", ttrack.Codec, pts)
					for i := range packets {
						b.Write(packets[i])
					}
				})
			}
		}
		return nil
	}

	// start reading
	err := c.Start()
	if err != nil {
		g.debug(err.Error())
		return
	}
	defer c.Close()
	// r.stopPlayer(g)

	g.radio, err = g.audioContext.NewPlayer(&b)
	if err != nil {
		g.debug(err.Error())
		return
	}
	g.radio.Play()

	// wait for a fatal error
	err = <-c.Wait()
	g.debug(err.Error())
}
