package ledserver

import (
	"context"
	"image/color"
	"time"
)

var colors = map[string]color.RGBA{
	"red":    {0xff, 0x00, 0x00, 0x00},
	"orange": {0xff, 0xa5, 0x00, 0x00},
	"yellow": {0xff, 0xff, 0x00, 0x00},
	"green":  {0x00, 0xff, 0x00, 0x00},
	"blue":   {0x00, 0x00, 0xff, 0x00},
	"indigo": {0x4b, 0x00, 0x82, 0x00},
	"violet": {0x7f, 0x00, 0xff, 0x00},
}

func (l *LED) rainbow(ctx context.Context) {
	l.off(true)
	inorder := []string{"red", "orange", "yellow", "green", "blue", "indigo", "violet", "indigo", "blue", "green", "yellow", "orange"}

	// spin in go task until ctx.Done()
	go func(ctx context.Context) {
		for {
			for _, color := range inorder {
				for i := 0; i < l.bufsize; i += channels {
					cc := colors[color]
					l.buf[i] = cc.R
					l.buf[i+1] = cc.G
					l.buf[i+2] = cc.B
					select {
					case <-ctx.Done():
						l.off(true)
						return
					default:
						l.leds.Write(l.buf)
						time.Sleep(25 * time.Millisecond)
					}
				}
			}
		}
	}(ctx)
}

func (l *LED) startup_test(ctx context.Context) {
	l.off(true)

	go func(ctx context.Context) {
		// test each individual pixel, all three channels
		for i := 0; i < l.bufsize; i += channels {
			for j := 0; j < channels; j++ {
				select {
				case <-ctx.Done():
					return
				default:
					l.buf[i+j] = 0x1f
					l.leds.Write(l.buf)
					time.Sleep(25 * time.Millisecond)
					l.buf[i+j] = 0x00
					l.leds.Write(l.buf)
					time.Sleep(25 * time.Millisecond)
				}
			}
		}
	}(ctx)
}
