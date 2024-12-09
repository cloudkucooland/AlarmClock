package main

import (
	"fmt"

	"github.com/bluenviron/gohlslib/v2"
	"github.com/bluenviron/gohlslib/v2/pkg/codecs"
)

func (g *Game) playhls(url string) {
	c := &gohlslib.Client{
		URI: url,
	}

	// called when tracks are parsed
	c.OnTracks = func(tracks []*gohlslib.Track) error {
		for _, track := range tracks {
			ttrack := track

			fmt.Printf("detected track with codec %T\n", track.Codec)

			// set a callback that is called when data is received
			switch track.Codec.(type) {
			case *codecs.AV1:
				c.OnDataAV1(track, func(pts int64, tu [][]byte) {
					fmt.Printf("received data from track %T, pts = %v\n", ttrack.Codec, pts)
				})

			case *codecs.H264, *codecs.H265:
				c.OnDataH26x(track, func(pts int64, dts int64, au [][]byte) {
					fmt.Printf("received data from track %T, pts = %v\n", ttrack.Codec, pts)
				})

			case *codecs.MPEG4Audio:
				c.OnDataMPEG4Audio(track, func(pts int64, aus [][]byte) {
					fmt.Printf("received data from track %T, pts = %v\n", ttrack.Codec, pts)
				})

			case *codecs.Opus:
				c.OnDataOpus(track, func(pts int64, packets [][]byte) {
					fmt.Printf("received data from track %T, pts = %v\n", ttrack.Codec, pts)
				})
			default:
				fmt.Println("something")
			}
		}
		return nil
	}

	// start reading
	if err := c.Start(); err != nil {
		fmt.Println("unable to start", err.Error())
		return
	}
	defer c.Close()

	// wait for a fatal error
	err := <-c.Wait()
	fmt.Println(err.Error())
}
