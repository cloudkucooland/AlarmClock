package ledserver

import (
	"image/color"

	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"

	"github.com/lucasb-eyer/go-colorful"
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
	c := colorful.Hsl(hsb.H, hsb.S/100, float64(hsb.B)/100)
	rgb.R, rgb.G, rgb.B= c.RGB255()
	return rgb
}
