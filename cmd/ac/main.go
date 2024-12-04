package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	// "math"
	// "strings"


	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	// "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	inScreenSaver bool
}

var (
	mplusFaceSource *text.GoTextFaceSource
	mplusFace    *text.GoTextFace
)

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	black  := color.RGBA{0x00, 0x00, 0x00, 0xff}

	{
		const x, y = 2, 2
		w, h := text.Measure("something", mplusFace, mplusFace.Size*1.5)
		vector.DrawFilledRect(screen, x, y, float32(w), float32(h), black, false)
		op := &text.DrawOptions{}
		op.GeoM.Translate(x, y)
		op.LineSpacing = mplusFace.Size * 1.5
		text.Draw(screen, "something", mplusFace, op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	mplusFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}

	g := &Game{
		inScreenSaver:    false,
	}

	ebiten.SetWindowSize(800, 480)
	// ebiten.SetFullscreen(true)

	ebiten.SetWindowTitle("Alarm Clock")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

