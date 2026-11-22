package ac

import (
	"fmt"
	"image/color"
	"net/rpc"
	"time"

	"github.com/cloudkucooland/AlarmClock/internal/ledserver"
)

func (g *Game) ledAllOn() {
	g.sendLEDCommand(&ledserver.Command{
		Command: ledserver.AllOn,
		Color:   color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	})
}

func (g *Game) ledRainbow() {
	g.sendLEDCommand(&ledserver.Command{
		Command: ledserver.Rainbow,
	})
}

func (g *Game) ledAllOff() {
	g.sendLEDCommand(&ledserver.Command{
		Command: ledserver.AllOn,
		Color:   color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
	})
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
	if g.ledclient != nil {
		_ = g.ledclient.Close()
		g.ledclient = nil
	}
}

func (g *Game) sendLEDCommand(cmd *ledserver.Command) {
	go func() {
		g.ledMu.Lock()
		defer g.ledMu.Unlock()

		if g.ledclient == nil && !g.ledConnect() {
			return
		}

		var res ledserver.Result
		done := make(chan *rpc.Call, 1)
		g.ledclient.Go("LED.Set", cmd, &res, done)

		select {
		case call := <-done:
			if call.Error != nil {
				g.debug(call.Error.Error())
				g.ledDisconnect()
			}
		case <-time.After(500 * time.Millisecond):
			g.debug("ledserver RPC timeout")
			g.ledDisconnect()
		}
	}()
}
