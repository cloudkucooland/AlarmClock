package ledserver

import (
	"fmt"
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

	l.hk = &a

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
		s.led.stopRunning()
		if !newstate {
			s.led.off(false)
			return
		}

		hsb := HSB{
			H: s.Hue.Value(),
			S: s.Saturation.Value(),
			B: s.Brightness.Value(),
		}
		s.led.staticColor(hsb.ToRGB(), false)
	})

	s.Brightness = characteristic.NewBrightness()
	s.AddC(s.Brightness.C)
	s.Brightness.OnValueRemoteUpdate(func(newstate int) {
		s.led.stopRunning()
		hsb := HSB{
			H: s.Hue.Value(),
			S: s.Saturation.Value(),
			B: newstate,
		}
		s.led.staticColor(hsb.ToRGB(), false)
	})

	s.Saturation = characteristic.NewSaturation()
	s.AddC(s.Saturation.C)
	s.Saturation.OnValueRemoteUpdate(func(newstate float64) {
		s.led.stopRunning()
		hsb := HSB{
			H: s.Hue.Value(),
			S: newstate,
			B: s.Brightness.Value(),
		}
		s.led.staticColor(hsb.ToRGB(), false)
	})

	s.Hue = characteristic.NewHue()
	s.AddC(s.Hue.C)
	s.Hue.OnValueRemoteUpdate(func(newstate float64) {
		s.led.stopRunning()
		hsb := HSB{
			H: newstate,
			S: s.Saturation.Value(),
			B: s.Brightness.Value(),
		}
		s.led.staticColor(hsb.ToRGB(), false)
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
	rgb.R, rgb.G, rgb.B = c.RGB255()
	return rgb
}

func (l *LED) updateHomeKit(c color.RGBA) {
	if l.hk == nil {
		return
	}

	cf := colorful.LinearRgb(float64(c.R)/255, float64(c.G)/255, float64(c.B)/255)
	h, s, b := cf.Hsl()
	fmt.Printf("H: %f S: %f B: %f\n", h, s*100, b*100)
	l.hk.Lightbulb.Hue.SetValue(h * 360)
	l.hk.Lightbulb.Saturation.SetValue(s * 100)
	l.hk.Lightbulb.Brightness.SetValue(int(b * 100))
}
