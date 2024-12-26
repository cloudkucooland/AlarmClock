package main

import (
	"github.com/cloudkucooland/AlarmClock/ledserver"
)

func (g *Game) ledAllOn() {
	if g.ledclient == nil {
		return
	}

	var res ledserver.Result

	cmd := &ledserver.Command{
		Command:    ledserver.AllOn,
		Brightness: 255,
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
		Command:    ledserver.Rainbow,
		Brightness: 255,
	}

	if err := g.ledclient.Call("LED.Set", cmd, &res); err != nil {
		g.debug(err.Error())
	}
}

func (g *Game) ledFrontOn() {
	if g.ledclient == nil {
		return
	}

	var res ledserver.Result

	cmd := &ledserver.Command{
		Command:    ledserver.FrontOn,
		Brightness: 255,
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
		Command:    ledserver.AllOn,
		Brightness: 0,
	}

	if err := g.ledclient.Call("LED.Set", cmd, &res); err != nil {
		g.debug(err.Error())
	}
}
