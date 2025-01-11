package main

import (
	"image/color"

	"github.com/cloudkucooland/AlarmClock/ledserver"
)

func (g *Game) ledAllOn() {
	if g.ledclient == nil {
		return
	}

	var res ledserver.Result

	cmd := &ledserver.Command{
		Command: ledserver.AllOn,
		Color:   color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	}

	if err := g.ledclient.Call("LED.Set", cmd, &res); err != nil {
		g.debug(err.Error())
	}
}

func (g *Game) ledRainbow() {
	if g.ledclient == nil {
		return
	}

	var res ledserver.Result

	cmd := &ledserver.Command{
		Command: ledserver.Rainbow,
	}

	if err := g.ledclient.Call("LED.Set", cmd, &res); err != nil {
		g.debug(err.Error())
	}
}

func (g *Game) ledAllOff() {
	if g.ledclient == nil {
		return
	}

	var res ledserver.Result

	cmd := &ledserver.Command{
		Command: ledserver.AllOn,
		Color:   color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
	}

	if err := g.ledclient.Call("LED.Set", cmd, &res); err != nil {
		g.debug(err.Error())
	}
}
