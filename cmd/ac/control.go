package main

import (
	"fmt"
	"math"
	// "image"
	// "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
	// "github.com/cloudkucooland/AlarmClock/resources"
)

const (
	controlScale  = 1.5                         // the default scale for the icons
	controlIconY  = 32                          // how many pixles the icon is wide (should not be a const, determine from img)
	controlYspace = controlIconY * controlScale // how far down to position the text below the icon
)

type control struct {
	sprite *sprite
	label  string
	x      int
	y      int
	do     func(*Game)
	ani    *controlanimation
}

var controls = []control{
	{
		sprite: getSprite("Mad"),
		label:  "Alarms",
		x:      700,
		y:      20,
		do:     alarmConfigDialog,
		ani:    &controlanimation{},
	},
	{
		sprite: getSprite("Happy"),
		label:  "Radio",
		x:      700,
		y:      120,
		do:     radioDialog,
		ani:    &controlanimation{},
	},
	{
		sprite: getSprite("Pinwheel"),
		label:  "Weather",
		x:      700,
		y:      220,
		do:     weatherDialog,
		ani:    &controlanimation{},
	},
}

func (g *Game) drawControls(screen *ebiten.Image) {
	for x := range controls {
		if !controls[x].onscreen() {
			continue
		}

		if controls[x].ani.in {
			controls[x].aniStep(screen)
			continue
		}
		controls[x].stillIcon(screen)
	}
}

func (c *control) stillIcon(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(controlScale, controlScale)
	op.GeoM.Translate(float64(c.x), float64(c.y))
	screen.DrawImage(c.sprite.image, op)

	top := &text.DrawOptions{}
	top.GeoM.Translate(float64(c.x), float64(c.y+controlYspace))
	top.LineSpacing = controlfont.Size * 1
	text.Draw(screen, c.label, controlfont, top)
}

func (c *control) in(x, y int) bool {
	return (x >= c.x && x <= c.x+controlIconY*controlScale) && (y >= c.y && y <= c.y+controlIconY*controlScale)
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

type controlanimation struct {
	in   bool
	step int
}

func (c *control) aniStep(screen *ebiten.Image) {
	c.ani.step = c.ani.step + 1

	scale := controlScale + scaleWibble(float64(c.ani.step))
	theta := thetaWibble(float64(c.ani.step))
	recenterx, recentery := locWibble(float64(c.x), float64(c.y), float64(c.ani.step))

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Rotate(theta)
	op.GeoM.Translate(recenterx, recentery)
	screen.DrawImage(c.sprite.image, op)

	top := &text.DrawOptions{}
	top.GeoM.Translate(float64(c.x), float64(c.y+controlYspace))
	top.LineSpacing = controlfont.Size
	text.Draw(screen, c.label, controlfont, top)

	if c.ani.step > (hz / 4) { // quarter of a second
		c.ani.step = 0
		c.ani.in = false
	}
}

func scaleWibble(i float64) float64 {
	return math.Sin(i/4) / 3
}

func thetaWibble(i float64) float64 {
	return math.Sin(i/6) / 6
}

func locWibble(x, y, step float64) (float64, float64) {
	z := thetaWibble(step) * 25
	return x + z, y + z
}

func (c *control) startanimation() {
	if c.ani.step != 0 {
		return
	}
	c.ani.step = 1
	c.ani.in = true
	c.playchirp()
}

func (c *control) playchirp() {
	// TODO
}
