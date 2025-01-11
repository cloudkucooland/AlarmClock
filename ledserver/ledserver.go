package ledserver

import (
	"fmt"
	"image/color"
	"time"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/nrzled"
	"periph.io/x/host/v3"
)

const Pipefile = "/tmp/ledserver.sock"
const pin = "SPI0.0"
const channels = 3
const numpixels = 16
var bufsize = channels * numpixels
var buf []byte
var leds *nrzled.Dev

type LED struct {}

type CommandCode int
const (
	AllOn CommandCode = iota
	Rainbow
)

type Command struct {
	Command CommandCode
	Color   color.RGBA
}

type Result bool

func (l *LED) Set(cmd *Command, res *Result) error {
	switch cmd.Command {
	case AllOn:
		// fmt.Printf("All on color: %+v\n", cmd.Color)
		staticColor(cmd.Color)
		*res = true
	case Rainbow:
		fmt.Println("rainbow")
		*res = true
		// do a rainbow
	}
	return nil
}

func (l *LED) Init() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	fmt.Printf("Using: %s\n", pin)

	s, err := spireg.Open(pin)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if p, ok := s.(spi.Pins); ok {
		fmt.Printf("Using pins CLK: %s  MOSI: %s  MISO: %s\n", p.CLK(), p.MOSI(), p.MISO())
	} else {
		err := fmt.Errorf("unable to open SPI")
		fmt.Println(err.Error())
		return err
	}

	leds, err = nrzled.NewSPI(s, &nrzled.Opts{
		NumPixels: numpixels,
		Freq:      2500 * physic.KiloHertz,
		Channels:  channels,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	buf = make([]byte, bufsize)

	go startup()
	return nil
}

func startup() {
	for i := 0; i < bufsize; i += channels {
		for j := 0; j < 3; j++ {
			buf[i + j] = 0xff
			leds.Write(buf)
			time.Sleep(100 * time.Millisecond)
			buf[i + j] = 0x00
			leds.Write(buf)
		}
	}

	white(0x88)
	time.Sleep(100 * time.Millisecond)
	white(0xff)
	time.Sleep(100 * time.Millisecond)
	white(0x22)
	time.Sleep(100 * time.Millisecond)
	white(0x00)
	leds.Halt()
}

func staticColor(c color.RGBA) {
	for i := 0; i < bufsize; i += channels {
		buf[i] = c.R
		buf[i+1] = c.G
		buf[i+2] = c.B
	}
	leds.Write(buf)
}

func white(brightness byte) {
	for i := 0; i < bufsize; i++ {
		buf[i] = brightness
	}
	leds.Write(buf)
}
