package main

import (
	"bytes"
	"image"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const bgpath = "/home/birdhouse/backgrounds"

func (g *Game) setBackground() {
	if g.background != nil {
		g.background.Deallocate()
		g.background = nil
	}

	// start with empty image
	g.background = ebiten.NewImage(screensize.X, screensize.Y)

	bgs, err := os.ReadDir(bgpath)
	if err != nil {
		g.debug(err.Error())
		return
	}

	// #nosec G404 G115
	key := uint8(rand.Intn(len(bgs)))

	bg := bgs[key]
	if !strings.HasSuffix(bg.Name(), ".png") {
		g.debug("file not PNG")
		return
	}
	raw, err := os.ReadFile(filepath.Join(bgpath, bg.Name()))
	if err != nil {
		g.debug(err.Error())
		return
	}
	decoded, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		g.debug(err.Error())
		return
	}
	img := ebiten.NewImageFromImage(decoded)

	// cache dim for screensaver
	if g.state == inScreenSaver {
		op := &colorm.DrawImageOptions{}
		op.Blend = ebiten.BlendCopy
		var cm colorm.ColorM
		cm.Scale(1.0, 1.0, 1.0, 0.05)
		colorm.DrawImage(g.background, img, cm, op)
		return
	}

	// cache normal
	g.background.DrawImage(img, &ebiten.DrawImageOptions{})
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	if g.background != nil {
		screen.DrawImage(g.background, &ebiten.DrawImageOptions{})
		return
	}
}
