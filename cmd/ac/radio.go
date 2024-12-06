package main

import (
	// "fmt"
	// "image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type radiobutton struct {
	sprite *sprite
	label  string
	url    string
	x int
	y int
}

var radiobuttons = []radiobutton{
	{
		sprite: getSprite("Indignent"),
		label:  "BBC 6Music",
	},
	{
		sprite: getSprite("Love"),
		label:  "WRR",
	},
	{
		sprite: getSprite("Pinwheel"),
		label:  "BBC 4",
	},
	{
		sprite: getSprite("Mad"),
		label:  "/dev/io",
	},
	{
		sprite: getSprite("Swan Mommy"),
		label:  "Test",
	},
	{
		sprite: getSprite("Spring"),
		label:  "Another test",
	},
	{
		sprite: getSprite("Indignent"),
		label:  "Quit",
	},
}

func radioDialog(g *Game) {
	g.state = inRadio
}

func (g *Game) drawRadio(screen *ebiten.Image) {
	grey := color.RGBA{0xaa, 0xaa, 0xaa, 0x99}
	border := color.RGBA{0x66, 0x66, 0x66, 0x00}
	vector.DrawFilledRect(screen, float32(20), float32(20), float32(760), float32(440), grey, false)
	vector.StrokeRect(screen, float32(20), float32(20), float32(760), float32(440), float32(4), border, false)
	vector.StrokeRect(screen, float32(30), float32(30), float32(740), float32(420), float32(2), border, false)

	x := 40
	y := 40

	for idx := range radiobuttons {

		radiobuttons[idx].x = x
		radiobuttons[idx].y = y

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(controlScale, controlScale)
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(radiobuttons[idx].sprite.image, op)

		top := &text.DrawOptions{}
		top.GeoM.Translate(float64(x), float64(y+controlYspace))
		top.LineSpacing = controlfont.Size
		text.Draw(screen, radiobuttons[idx].label, controlfont, top)

		x = x + 112
		if (idx % 5) == 4 {
			x = 40
			y = y + 112
		}
	}

	x = 720
	y = 380
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(controlScale, controlScale)
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(radiobuttons[0].sprite.image, op)

	top := &text.DrawOptions{}
	top.GeoM.Translate(float64(x), float64(y+controlYspace))
	top.LineSpacing = controlfont.Size
	text.Draw(screen, "Close", controlfont, top)

}

func (r radiobutton) in(x, y int) bool {
	return (x >= r.x && x <= r.x+controlIconY*controlScale) && (y >= r.y && y <= r.y+controlIconY*controlScale)
}

func (r *radiobutton) onscreen() bool {
	return !(r.x == 0 && r.y == 0)
}
