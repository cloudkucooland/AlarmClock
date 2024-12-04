package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	// "strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/cloudkucooland/AlarmClock/resources"
)

const defaultClockLocationX = 5
const defaultClockLocationY = 20

type Game struct {
	inScreenSaver   bool
	lastAct         time.Time
	clockLocationX  int
	clockLocationY  int
	cyclesSinceTick int
	timestring      string
}

type sprite struct {
	name  string
	raw   []byte
	image *ebiten.Image
	alpha *image.Alpha
	x     int
	y     int
	// Do	*(func)()
}

func (s *sprite) In(x, y int) bool {
	return s.alpha.At(x-s.x, y-s.y).(color.Alpha).A > 0
}

var (
	spaceMonoSource *text.GoTextFaceSource
	spaceMonoFace   *text.GoTextFace
)

var sprites = []sprite{
	{
		name: "Artist",
		raw:  resources.ArtistPNG,
		x:    700,
		y:    20,
	},
	{
		name: "Baby",
		raw:  resources.BabyPNG,
		x:    700,
		y:    120,
	},
	{
		name: "Bathtime",
		raw:  resources.BathtimePNG,
		x:    700,
		y:    220,
	},
	{
		name: "Confused",
		raw:  resources.ConfusedPNG,
		x:    700,
		y:    320,
	},
	{
		name: "Happy",
		raw:  resources.HappyPNG,
		x:    700,
		y:    420,
	},
}

func (g *Game) Update() error {
	g.cyclesSinceTick = (g.cyclesSinceTick + 1) % 3600 // assumes 60Hz
	if g.cyclesSinceTick == 1 {
		now := time.Now()
		g.timestring = now.Format("03:04")
		if g.inScreenSaver {
			g.updateClockLocation()
		}
		if !g.inScreenSaver && now.After(g.lastAct.Add(5*time.Minute)) {
			// start screen saver
			g.inScreenSaver = true
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.lastAct = time.Now()
		if g.inScreenSaver {
			g.leaveScreenSaver()
		}
		x, y := ebiten.CursorPosition()
		fmt.Printf("Mouse position: x=%d, y=%d\n", x, y)

		for _, s := range sprites {
			if s.In(x, y) {
				fmt.Printf("in sprite %s\n", s.name)
				// s.Do()
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		log.Fatal("shutting down")
	}

	return nil
}

func (g *Game) updateClockLocation() {
	// get clock size, determine max range...
	g.clockLocationX = rand.Int() % 100
	g.clockLocationY = rand.Int() % 200
}

func (g *Game) leaveScreenSaver() {
	g.inScreenSaver = false
	g.clockLocationX = defaultClockLocationX
	g.clockLocationY = defaultClockLocationY
}

func (g *Game) Draw(screen *ebiten.Image) {
	// black  := color.RGBA{0x00, 0x00, 0x00, 0x00}
	// white := color.RGBA{0xff, 0xff, 0xff, 0x00}

	if !g.inScreenSaver {
		for x := range sprites {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.05, 0.05)
			op.GeoM.Translate(float64(sprites[x].x), float64(sprites[x].y))
			screen.DrawImage(sprites[x].image, op)
		}
	}

	{
		// w, h := text.Measure(g.timestring, spaceMonoFace, spaceMonoFace.Size*1.2)
		// vector.DrawFilledRect(screen, g.clockLocationX, g.clockLocationY, float32(w), float32(h), black, false)
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(g.clockLocationX), float64(g.clockLocationY))
		if g.inScreenSaver {
			op.ColorScale.ScaleAlpha(0.25)
		}
		op.LineSpacing = spaceMonoFace.Size * 1
		text.Draw(screen, g.timestring, spaceMonoFace, op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(resources.SpaceMonoBold_ttf))
	if err != nil {
		log.Fatal(err)
	}
	spaceMonoSource = s

	spaceMonoFace = &text.GoTextFace{
		Source: spaceMonoSource,
		Size:   192,
	}

	for x := range sprites {
		img, _, err := image.Decode(bytes.NewReader(sprites[x].raw))
		if err != nil {
			log.Fatal(err)
		}
		sprites[x].image = ebiten.NewImageFromImage(img)
		b := img.Bounds()
		sprites[x].alpha = image.NewAlpha(b)
		for j := b.Min.Y; j < b.Max.Y; j++ {
			for i := b.Min.X; i < b.Max.X; i++ {
				sprites[x].alpha.Set(i, j, img.At(i, j))
			}
		}
	}

	g := &Game{
		inScreenSaver:  false,
		clockLocationX: defaultClockLocationX,
		clockLocationY: defaultClockLocationY,
	}
	now := time.Now()
	g.lastAct = now
	g.inScreenSaver = false
	// attempt to get the minute-change correct...
	ms := now.Sub(now.Truncate(time.Second))
	g.cyclesSinceTick = int(ms.Milliseconds() * 60 / 1000) // assumes 60Hz

	g.timestring = now.Format("03:04")

	ebiten.SetWindowSize(800, 480)
	// ebiten.SetFullscreen(true)

	ebiten.SetWindowTitle("Alarm Clock")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
