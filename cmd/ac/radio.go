package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"
)

type radiobutton struct {
	sprite *sprite
	label  string
	url    string
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
}

func radioDialog(g *Game) {
	g.state = inRadio
}

func (g *Game) drawRadio(screen *ebiten.Image) {
	g.drawModal(screen)

	paddedspritesize := 40 // spritesize * 1.5? modal border + 10?
	x := paddedspritesize
	y := paddedspritesize

	for idx := range radiobuttons {
		radiobuttons[idx].sprite.setLocation(x, y)
		radiobuttons[idx].sprite.setScale(spriteScale)
		radiobuttons[idx].sprite.draw(screen)
		radiobuttons[idx].sprite.drawLabel(radiobuttons[idx].label, screen)

		// draw label
		// TODO: base this on sprite size not hardcoded values
		perrow := 5       // make dynamic
		rowspacing := 112 // make dynamic
		x = x + rowspacing
		if (idx % perrow) == (perrow - 1) {
			x = paddedspritesize
			y = y + rowspacing
		}
	}
}
