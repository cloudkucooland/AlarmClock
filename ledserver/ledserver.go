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

type LED struct {
	buf     []byte
	leds    *nrzled.Dev
	bufsize int
	reg     spi.PortCloser
}

type CommandCode int

const (
	AllOn CommandCode = iota
	Startup
	Rainbow
	Off
)

type Command struct {
	Command CommandCode
	Color   color.RGBA
}

type Result bool

func (l *LED) Set(cmd *Command, res *Result) error {
	switch cmd.Command {
	case AllOn:
		l.stopRunning()
		l.staticColor(cmd.Color)
		*res = true
	case Startup:
		l.stopRunning()
		l.startup_test()
		*res = true
	case Rainbow:
		l.stopRunning()
		l.rainbow()
		*res = true
	case Off:
		l.stopRunning()
		l.off()
	}
	return nil
}

func (l *LED) Init() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	fmt.Printf("Using: %s\n", pin)

	var err error
	l.reg, err = spireg.Open(pin)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	pins, ok := l.reg.(spi.Pins)
	if !ok {
		err := fmt.Errorf("unable to open SPI")
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("Using pins CLK: %s  MOSI: %s  MISO: %s\n", pins.CLK(), pins.MOSI(), pins.MISO())

	l.leds, err = nrzled.NewSPI(l.reg, &nrzled.Opts{
		NumPixels: numpixels,
		Freq:      2500 * physic.KiloHertz,
		Channels:  channels,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	l.bufsize = numpixels * channels
	l.buf = make([]byte, l.bufsize)

	// go l.startup_test()
	return nil
}

func (l *LED) Shutdown() {
	l.leds.Halt()
	l.reg.Close()
}

func (l *LED) startup_test() {
	l.white(0x00)

	// test each individual pixel, all three channels
	for i := 0; i < l.bufsize; i += channels {
		for j := 0; j < 3; j++ {
			l.buf[i+j] = 0x0f
			l.leds.Write(l.buf)
			time.Sleep(100 * time.Millisecond)
			l.buf[i+j] = 0x00
			l.leds.Write(l.buf)
		}
	}

	// step down from full to dim
	steps := []byte{0xff, 0xdd, 0xbb, 0x99, 0x77, 0x55, 0x33, 0x11, 0x00}
	for _, v := range steps {
		l.white(v)
		time.Sleep(100 * time.Millisecond)
	}
	l.leds.Halt()
}

func (l *LED) staticColor(c color.RGBA) {
	l.white(0x00)
	for i := 0; i < l.bufsize; i += channels {
		l.buf[i] = c.R
		l.buf[i+1] = c.G
		l.buf[i+2] = c.B
	}
	l.leds.Write(l.buf)
}

func (l *LED) white(brightness byte) {
	// Google's AI hallucinated a slices.Fill function to do this... but alas it does not exist
	for i := 0; i < l.bufsize; i++ {
		l.buf[i] = brightness
	}
	l.leds.Write(l.buf)
}

func (l *LED) off() {
	l.white(0x00)
	l.leds.Halt()
}

func (l *LED) stopRunning() {
	// stop any running sequence
}
