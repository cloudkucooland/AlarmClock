package ledserver

import (
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

func (l *LED) rainbow() {
	l.white(0x00)

	inorder := []string{"red", "orange", "yellow", "green", "blue", "indigo", "violet", "indigo", "blue", "green", "yellow", "orange"}

	l.mu.Lock()
	for i := 0; i < 4; i++ {
		for _, color := range inorder {
			for i := 0; i < l.bufsize; i += channels {
				l.buf[i] = colors[color].R
				l.buf[i+1] = colors[color].G
				l.buf[i+2] = colors[color].B
				l.leds.Write(l.buf)
				time.Sleep(25 * time.Millisecond)
			}
		}
	}
	l.mu.Unlock()

	l.leds.Halt()
}
