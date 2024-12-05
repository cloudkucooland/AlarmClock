package main

import (
	"fmt"
	// "image"
	// "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
	// "github.com/cloudkucooland/AlarmClock/resources"
)

const (
	controlScale  = 1.5
	controlIcony  = 32
	controlYspace = controlIcony * controlScale
)

type control struct {
	sprite *sprite
	label  string
	x      int
	y      int
	do     func(*Game)
}

var controls = []control{
	{
		sprite: &sprites[0],
		label:  "Alarms",
		x:      700,
		y:      20,
		do:     defaultAction,
	},
	{
		sprite: &sprites[1],
		label:  "Radio",
		x:      700,
		y:      120,
		do:     radioDialog,
	},
	{
		sprite: &sprites[2],
		label:  "Weather",
		x:      700,
		y:      220,
		do:     defaultAction,
	},
}

func (g *Game) drawControls(screen *ebiten.Image) {
	// white := color.RGBA{0xff, 0xff, 0xff, 0xdd}

	if !g.inScreenSaver() {
		for x := range controls {
			if !controls[x].onscreen() {
				continue
			}

			// vector.DrawFilledRect(screen, float32(controls[x].x), float32(controls[x].y), float32(50), float32(50), white, false)

			// w, h := text.Measure(controls[x].label, controlfont, controlfontfont.Size*1.2)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(controlScale, controlScale)
			op.GeoM.Translate(float64(controls[x].x), float64(controls[x].y))
			screen.DrawImage(controls[x].sprite.image, op)

			top := &text.DrawOptions{}
			top.GeoM.Translate(float64(controls[x].x), float64(controls[x].y+controlYspace))
			top.LineSpacing = controlfont.Size * 1
			text.Draw(screen, controls[x].label, controlfont, top)
		}
	}
}

func (c control) in(x, y int) bool {
	if (x >= c.x && x <= c.x+controlIcony*controlScale) && (y >= c.y && y <= c.y+controlIcony*controlScale) {
		return true
	}
	return false
}

func buildControls() error {
	return nil
}

func defaultAction(g *Game) {
	fmt.Println("control defaultAction")
}

func (c *control) onscreen() bool {
	if c.x == 0 && c.y == 0 {
		return false
	}
	return true
}
