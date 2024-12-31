package ledserver

import (
	"errors"
	"fmt"
	"image/color"

	// "periph.io/x/conn/v3/display"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/gpio/gpiostream"
	"periph.io/x/conn/v3/physic"
	// "periph.io/x/conn/v3/spi"
	// "periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/nrzled"
	"periph.io/x/host/v3"
)

const Pipefile = "/tmp/ledserver.sock"
const pin = "GPIO12"
const channels = 4
const numpixels = 5

type LED struct {
	stream    *nrzled.Dev
	numpixels int
	buf       []byte
}

type CommandCode int

const (
	AllOn CommandCode = iota
	BackOn
	FrontOn
	Rainbow
)

type Command struct {
	Command CommandCode
	Color   color.RGBA
}

type Result bool

func (l *LED) Set(cmd *Command, res *Result) error {
	fmt.Printf("%+v", cmd)

	switch cmd.Command {
	case AllOn:
		l.White(255)
	case BackOn:
		fmt.Println("back on")
		// just the back
	case FrontOn:
		fmt.Println("front on")
		// just the front
	case Rainbow:
		fmt.Println("rainbow")
		// do a rainbow
	}
	return nil
}

func (l *LED) Init() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	p := gpioreg.ByName(pin)
	if p == nil {
		return errors.New("specify a valid pin")
	}
	if rp, ok := p.(gpio.RealPin); ok {
		p = rp.Real()
	}
	s, ok := p.(gpiostream.PinOut)
	if !ok {
		return fmt.Errorf("pin %s doesn't support arbitrary bit stream", p)
	}

	opts := nrzled.DefaultOpts
	opts.NumPixels = numpixels
	opts.Freq = 2500 * physic.KiloHertz
	opts.Channels = channels

	l.numpixels = opts.NumPixels
	l.buf = make([]byte, l.numpixels*(channels))

	var err error
	if l.stream, err = nrzled.NewStream(s, &opts); err != nil {
		return err
	}

	l.StaticColor(color.RGBA{
		R: 0x20,
		G: 0x00,
		B: 0x22,
		A: 0x00,
	})
	return nil
}

func (l *LED) StaticColor(c color.RGBA) error {
	for i := 0; i < len(l.buf); i += channels {
		l.buf[i] = c.R
		l.buf[i+1] = c.G
		l.buf[i+2] = c.B
		l.buf[i+3] = c.A
	}

	if _, err := l.stream.Write(l.buf); err != nil {
		return err
	}
	return nil
}

func (l *LED) White(brightness byte) error {
	for i := 0; i < len(l.buf); i += channels {
		l.buf[i] = 0x00   // r
		l.buf[i+1] = 0x00 // g
		l.buf[i+2] = 0x00 // b
		l.buf[i+3] = brightness
	}
	if _, err := l.stream.Write(l.buf); err != nil {
		return err
	}
	return nil
}
