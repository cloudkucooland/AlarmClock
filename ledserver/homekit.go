package ledserver

import (
	"image/color"
	"math"

	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
	// "https://github.com/gerow/go-color"
)

type LedServer struct {
	*accessory.A

	Lightbulb *LedServerSvc
}

func NewLedServer(info accessory.Info, l *LED) *LedServer {
	a := LedServer{}

	a.A = accessory.New(info, accessory.TypeLightbulb)
	a.Lightbulb = NewLedServerSvc(l)
	a.AddS(a.Lightbulb.S)

	return &a
}

type LedServerSvc struct {
	*service.S

	led *LED

	On         *characteristic.On
	Brightness *characteristic.Brightness
	Saturation *characteristic.Saturation
	Hue        *characteristic.Hue
}

func NewLedServerSvc(l *LED) *LedServerSvc {
	s := LedServerSvc{}
	s.S = service.New(service.TypeLightbulb)

	s.led = l

	s.On = characteristic.NewOn()
	s.AddC(s.On.C)
	s.On.OnValueRemoteUpdate(func(newstate bool) {
		if !newstate {
			s.led.white(0x00)
			return
		}

		hsb := HSB{
			H: s.Hue.Value(),
			S: s.Saturation.Value(),
			B: s.Brightness.Value(),
		}
		s.led.staticColor(hsb.ToRGB())
	})

	s.Brightness = characteristic.NewBrightness()
	s.AddC(s.Brightness.C)
	s.Brightness.OnValueRemoteUpdate(func(newstate int) {
		hsb := HSB{
			H: s.Hue.Value(),
			S: s.Saturation.Value(),
			B: newstate,
		}
		s.led.staticColor(hsb.ToRGB())
	})

	s.Saturation = characteristic.NewSaturation()
	s.AddC(s.Saturation.C)
	s.Saturation.OnValueRemoteUpdate(func(newstate float64) {
		hsb := HSB{
			H: s.Hue.Value(),
			S: newstate,
			B: s.Brightness.Value(),
		}
		s.led.staticColor(hsb.ToRGB())
	})

	s.Hue = characteristic.NewHue()
	s.AddC(s.Hue.C)
	s.Hue.OnValueRemoteUpdate(func(newstate float64) {
		hsb := HSB{
			H: newstate,
			S: s.Saturation.Value(),
			B: s.Brightness.Value(),
		}
		s.led.staticColor(hsb.ToRGB())
	})

	return &s
}

type HSB struct {
	H float64
	S float64
	B int
}

func (hsb HSB) ToRGB() color.RGBA {
	rgb := color.RGBA{}

	// Normalize values to [0, 1]
	h := hsb.H / 360
	s := hsb.S / 100
	b := float64(hsb.B) / 100

	if s == 0 {
		rgb.R = uint8(b * 255)
		rgb.G = uint8(b * 255)
		rgb.B = uint8(b * 255)
		return rgb
	}

	// var f func(n float64) float64
	f := func(n float64) float64 {
		k := math.Mod(n+h/60, 6)
		return b - b*s*math.Max(0, math.Min(math.Min(k, 4-k), 1))
	}

	rgb.R = uint8(f(5) * 255)
	rgb.G = uint8(f(3) * 255)
	rgb.B = uint8(f(1) * 255)
	return rgb
}
