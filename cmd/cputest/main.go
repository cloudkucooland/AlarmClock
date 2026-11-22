package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 480
}

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	g := &Game{}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err.Error())
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Update() error {
	return nil
}
