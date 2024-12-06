package main

import (
	// "fmt"
	// "image"
	// "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
)

type alarmbutton struct {
	x int
	y int
}

func alarmConfigDialog(g *Game) {
	g.state = inAlarmConfig
}

func (g *Game) drawAlarmConfig(screen *ebiten.Image) {
	g.drawModal(screen)

	/* x := 40
	y := 40

	for idx := range alarmbuttons {
		alarmbuttons[idx].x = x
		alarmbuttons[idx].y = y

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(controlScale, controlScale)
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(alarmbuttons[idx].sprite.image, op)

		top := &text.DrawOptions{}
		top.GeoM.Translate(float64(x), float64(y+controlYspace))
		top.LineSpacing = controlfont.Size
		text.Draw(screen, alarmbuttons[idx].label, controlfont, top)

		x = x + 112
		if (idx % 5) == 4 {
			x = 40
			y = y + 112
		}
	} */
}

func (r alarmbutton) in(x, y int) bool {
	return (x >= r.x && x <= r.x+controlIconY*controlScale) && (y >= r.y && y <= r.y+controlIconY*controlScale)
}

func (r *alarmbutton) onscreen() bool {
	return !(r.x == 0 && r.y == 0)
}
