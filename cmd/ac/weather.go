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

const owmzipcode = "75035"
const owmcountry = "US"

func (g *Game) drawWeather(screen *ebiten.Image) {
	if g.weathercache == nil {
		return
	}

	b := g.weathercache.Bounds()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(screensize.X/2)-float64(b.Max.X/2), float64(screensize.Y-40))
	screen.DrawImage(g.weathercache, op)
}

func (g *Game) runWeather(ctx context.Context) {
	apikey := os.Getenv("OWM_API_KEY")
	if apikey == "" {
		err := fmt.Errorf("OWM_API_KEY not set; not running weather poller")
		fmt.Println(err.Error())
		return
	}
	w, err := owm.NewCurrent("F", "EN", apikey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	g.weather = w

	if err := w.CurrentByZipcode(owmzipcode, owmcountry); err != nil {
		fmt.Println(err.Error())
	} else {
		updateweathercache(g)
	}

	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.CurrentByZipcode(owmzipcode, owmcountry); err != nil {
				fmt.Println(err.Error())
				if g.weathercache != nil {
					g.weathercache.Deallocate()
					g.weathercache = nil
				}
				continue
			}
			updateweathercache(g)
		}
	}
}

func updateweathercache(g *Game) {
	weatherstring := fmt.Sprintf("Current: %.1f Feels Like: %.1f High: %.1f Low: %.1f ", g.weather.Main.Temp, g.weather.Main.FeelsLike, g.weather.Main.TempMax, g.weather.Main.TempMin)
	textwidth, textheight := text.Measure(weatherstring, weatherfont, 1)

	if g.weathercache != nil {
		g.weathercache.Deallocate()
		g.weathercache = nil
	}

	g.weathercache = ebiten.NewImage(int(textwidth), int(textheight))

	op := &text.DrawOptions{}
	text.Draw(g.weathercache, weatherstring, weatherfont, op)
}
