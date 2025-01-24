package main

import (
	"fmt"
	"image/color"
	"net/rpc"

	"github.com/cloudkucooland/AlarmClock/ledserver"
)

func (g *Game) ledAllOn() {
	if g.ledclient == nil && !g.ledConnect() {
		return
	}

	var res ledserver.Result

	cmd := &ledserver.Command{
		Command: ledserver.AllOn,
		Color:   color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	}

	if err := g.ledclient.Call("LED.Set", cmd, &res); err != nil {
		g.debug(err.Error())
		g.ledDisconnect()
	}
}

func (g *Game) ledRainbow() {
	if g.ledclient == nil && !g.ledConnect() {
		return
	}

	var res ledserver.Result

	cmd := &ledserver.Command{
		Command: ledserver.Rainbow,
	}

	if err := g.ledclient.Call("LED.Set", cmd, &res); err != nil {
		g.debug(err.Error())
		g.ledDisconnect()
	}
}

func (g *Game) ledAllOff() {
	if g.ledclient == nil && !g.ledConnect() {
		return
	}

	var res ledserver.Result

	cmd := &ledserver.Command{
		Command: ledserver.AllOn,
		Color:   color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
	}

	if err := g.ledclient.Call("LED.Set", cmd, &res); err != nil {
		g.debug(err.Error())
		g.ledDisconnect()
	}
}

func (g *Game) ledConnect() bool {
	if client, err := rpc.DialHTTP("unix", ledserver.Pipefile); err != nil {
		err := fmt.Errorf("led server not connected: %s", err.Error())
		g.debug(err.Error())
		return false
	} else {
		g.ledclient = client
	}
	return true
}

func (g *Game) ledDisconnect() {
	g.ledclient.Close()
	g.ledclient = nil
}
