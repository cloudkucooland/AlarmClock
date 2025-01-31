package ledserver

import (
	"context"
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
const numpixels = 32

var runningHot bool

type LED struct {
	buf     []byte
	leds    *nrzled.Dev
	bufsize int
	reg     spi.PortCloser
	hk      *LedServer
	cancel  func()
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
	*res = true
	l.stopRunning()

	switch cmd.Command {
	case AllOn:
		l.staticColor(cmd.Color, true)
	case Startup:
		var ctx context.Context
		ctx, l.cancel = context.WithCancel(context.Background())
		l.startup_test(ctx)
	case Rainbow:
		var ctx context.Context
		ctx, l.cancel = context.WithCancel(context.Background())
		l.rainbow(ctx)
	case Off:
		l.off(true)
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

	var ctx context.Context
	ctx, l.cancel = context.WithCancel(context.Background())
	l.startup_test(ctx)
	return nil
}

func (l *LED) Shutdown() {
	l.leds.Halt()
	l.reg.Close()
}

func (l *LED) staticColor(c color.RGBA, updateHomekit bool) {
	for i := 0; i < l.bufsize; i += channels {
		// if we are running hot only light up every 4th led
		if runningHot && i%(4*channels) != 0 {
			continue
		}
		l.buf[i] = c.R
		l.buf[i+1] = c.G
		l.buf[i+2] = c.B
	}
	l.leds.Write(l.buf)
	if updateHomekit {
		l.updateHomeKit(c)
	}
}

func (l *LED) white(brightness byte, updateHomekit bool) {
	// Google's AI hallucinated a slices.Fill function to do this... but alas it does not exist
	for i := 0; i < l.bufsize; i++ {
		l.buf[i] = brightness
	}
	l.leds.Write(l.buf)
	if updateHomekit {
		l.updateHomeKit(color.RGBA{brightness, brightness, brightness, 0x00})
	}
}

func (l *LED) off(updateHomekit bool) {
	l.white(0x00, updateHomekit)
	if updateHomekit {
		l.updateHomeKit(color.RGBA{0x00, 0x00, 0x00, 0x00})
	}
	l.leds.Halt()
}

func (l *LED) stopRunning() {
	if l.cancel != nil {
		l.cancel()
		l.cancel = nil
		// there is a race here, this is a naïve way of ensuring the context cancel finished before we return
		time.Sleep(200 * time.Millisecond)
		l.leds.Halt()
	}
}

func ThermalHigh() {
	if !runningHot {
		fmt.Println("running hot")
		runningHot = true
	}
}

func ThermalNormal() {
	if runningHot {
		fmt.Println("cooled off")
		runningHot = false
	}
}
