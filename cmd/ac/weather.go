package main

import (
	// "fmt"
	// "image"
	// "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
)

type weatherbutton struct {
	x int
	y int
}

func weatherDialog(g *Game) {
	g.state = inWeather
}

func (g *Game) drawWeather(screen *ebiten.Image) {
	g.drawModal(screen)

	/* x := 40
	y := 40

	for idx := range weatherbuttons {
		weatherbuttons[idx].x = x
		weatherbuttons[idx].y = y

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(controlScale, controlScale)
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(weatherbuttons[idx].sprite.image, op)

		top := &text.DrawOptions{}
		top.GeoM.Translate(float64(x), float64(y+controlYspace))
		top.LineSpacing = controlfont.Size
		text.Draw(screen, weatherbuttons[idx].label, controlfont, top)

		x = x + 112
		if (idx % 5) == 4 {
			x = 40
			y = y + 112
		}
	} */
}

func (r weatherbutton) in(x, y int) bool {
	return (x >= r.x && x <= r.x+controlIconY*controlScale) && (y >= r.y && y <= r.y+controlIconY*controlScale)
}
