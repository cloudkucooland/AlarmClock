package main

import (
	"fmt"
	// "image"
	// "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
	// "github.com/cloudkucooland/AlarmClock/resources"
)

type radiobutton struct {
	sprite *sprite
	label  string
	x      int
	y      int
	url    string
}

var radiobuttons = []radiobutton{
	{
		sprite: &sprites[0],
		label:  "BBC 6Music",
		x:      100,
		y:      20,
	},
}

func radioDialog(g *Game) {
	fmt.Println("switching to radio dialog")
}

func (g *Game) drawradio(screen *ebiten.Image) {
	// white := color.RGBA{0xff, 0xff, 0xff, 0xdd}

	if !g.inScreenSaver() {
		for x := range radiobuttons {
			if !radiobuttons[x].onscreen() {
				continue
			}

			// vector.DrawFilledRect(screen, float32(controls[x].x), float32(controls[x].y), float32(50), float32(50), white, false)

			// w, h := text.Measure(controls[x].label, controlfont, controlfontfont.Size*1.2)

			/* op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(controlScale, controlScale)
			op.GeoM.Translate(float64(controls[x].x), float64(controls[x].y))
			screen.DrawImage(controls[x].sprite.image, op)

			top := &text.DrawOptions{}
			top.GeoM.Translate(float64(controls[x].x), float64(controls[x].y+controlYspace))
			top.LineSpacing = controlfont.Size * 1
			text.Draw(screen, controls[x].label, controlfont, top)
			 */
		}
	}
}

func (r radiobutton) in(x, y int) bool {
	if (x >= r.x && x <= r.x+controlIcony*controlScale) && (y >= r.y && y <= r.y+controlIcony*controlScale) {
		return true
	}
	return false
}

func (r *radiobutton) onscreen() bool {
	if r.x == 0 && r.y == 0 {
		return false
	}
	return true
}
