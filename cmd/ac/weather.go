package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	owm "github.com/briandowns/openweathermap"
)

func weatherDialog(g *Game) {
	g.state = inWeather
}

func (g *Game) drawWeather(screen *ebiten.Image) {
	g.drawModal(screen)

	x := 40
	y := 40

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = controlfont.Size * 1.5
	text.Draw(screen, g.weather, controlfont, op)
}

func (g *Game) runWeather(ctx context.Context) error {
	apikey := os.Getenv("OWM_API_KEY")
	if apikey == "" {
		err := fmt.Errorf("OWM_API_KEY not set; not running weather poller")
		g.weather = err.Error()
		return err
	}
	w, err := owm.NewCurrent("F", "EN", apikey)
	if err != nil {
		g.weather = err.Error()
		return err
	}

	w.CurrentByName("Frisco,TX")
	g.formatWeatherString(w)

	ticker := time.NewTicker(time.Hour)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			w.CurrentByName("Frisco,TX")
			g.formatWeatherString(w)
		}
	}
}

func (g *Game) formatWeatherString(c *owm.CurrentWeatherData) {
	g.weather = fmt.Sprintf("%+v", c)
}
