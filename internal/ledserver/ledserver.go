package ledserver

import (
	"context"
	"fmt"
	"image/color"
	"sync"
	"sync/atomic"
	"time"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/nrzled"
	"periph.io/x/host/v3"
)

const (
	Pipefile  = "/tmp/ledserver.sock"
	pin       = "SPI0.0"
	channels  = 3
	numpixels = 16
)

var runningHot atomic.Bool

type LED struct {
	mu      sync.Mutex
	buf     []byte
	leds    *nrzled.Dev
	bufsize int
	reg     spi.PortCloser
	hk      *LedServer
	cancel  func()
	done    chan struct{}
}

type CommandCode int

const (
	AllOn CommandCode = iota
	Startup
	Rainbow
	Off
	Sunrise
	Candle
	Fireflies
	AlertPulse
)

type Command struct {
	Command  CommandCode
	Color    color.RGBA
	Duration time.Duration // optional; used by Sunrise (defaults to 15m if <= 0)
}

type Result bool

func (l *LED) Set(cmd *Command, res *Result) error {
	*res = true
	l.stopRunning()

	switch cmd.Command {
	case AllOn:
		l.staticColor(cmd.Color, true)
	case Startup:
		ctx := l.startAnim()
		l.startup_test(ctx, l.done)
	case Rainbow:
		ctx := l.startAnim()
		l.rainbow(ctx, l.done)
	case Sunrise:
		ctx := l.startAnim()
		dur := cmd.Duration
		if dur <= 0 {
			dur = 15 * time.Minute
		}
		l.sunrise(ctx, l.done, dur)
	case Candle:
		ctx := l.startAnim()
		l.candle(ctx, l.done, 50)
	case Fireflies:
		ctx := l.startAnim()
		l.fireflies(ctx, l.done)
	case AlertPulse:
		ctx := l.startAnim()
		l.alertPulse(ctx, l.done, cmd.Color, cmd.Duration)
	case Off:
		l.off(true)
	}
	return nil
}

func (l *LED) startAnim() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel
	l.done = make(chan struct{})
	return ctx
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

	ctx := l.startAnim()
	l.startup_test(ctx, l.done)
	return nil
}

func (l *LED) Shutdown() {
	l.stopRunning()
	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.leds.Halt()
	_ = l.reg.Close()
}

func (l *LED) staticColor(c color.RGBA, updateHomekit bool) {
	l.mu.Lock()
	hot := runningHot.Load()
	for i := 0; i < l.bufsize; i += channels {
		// if we are running hot only light up every 4th led
		if hot && i%(4*channels) != 0 {
			l.buf[i] = 0
			l.buf[i+1] = 0
			l.buf[i+2] = 0
			continue
		}
		l.buf[i] = c.R
		l.buf[i+1] = c.G
		l.buf[i+2] = c.B
	}
	_, _ = l.leds.Write(l.buf)
	l.mu.Unlock()

	if updateHomekit {
		l.updateHomeKit(c)
	}
}

func (l *LED) white(brightness byte, updateHomekit bool) {
	l.mu.Lock()
	for i := range l.buf {
		l.buf[i] = brightness
	}
	_, _ = l.leds.Write(l.buf)
	l.mu.Unlock()

	if updateHomekit {
		l.updateHomeKit(color.RGBA{brightness, brightness, brightness, 0x00})
	}
}

func (l *LED) off(updateHomekit bool) {
	l.white(0x00, updateHomekit)
	l.mu.Lock()
	_ = l.leds.Halt()
	l.mu.Unlock()
}

func (l *LED) stopRunning() {
	if l.cancel != nil {
		l.cancel()
		l.cancel = nil
		if l.done != nil {
			select {
			case <-l.done:
			case <-time.After(500 * time.Millisecond):
			}
			l.done = nil
		}
	}
}

func ThermalHigh() {
	if !runningHot.Swap(true) {
		fmt.Println("running hot")
	}
}

func ThermalNormal() {
	if runningHot.Swap(false) {
		fmt.Println("cooled off")
	}
}
