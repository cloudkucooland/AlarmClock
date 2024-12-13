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

	w.CurrentByZip(75035, "US")
	g.formatWeatherString(w)

	ticker := time.NewTicker(time.Hour)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			w.CurrentByZip(75035, "US")
			g.formatWeatherString(w)
		}
	}
}

/* &{GeoPos:{Longitude:-96.7524 Latitude:33.1377} Sys:{Type:2 ID:2001388 Message:0 Country:US Sunrise:1733923242 Sunset:1733959236} Base:stations Weather:[{ID:802 Main:Clouds Description:scattered clouds Icon:03n}] Main:{Temp:44.67 TempMin:41.72 TempMax:47.44 FeelsLike:41.32 Pressure:1027 SeaLevel:1027 GrndLevel:1003 Humidity:48} Visibility:10000 Wind:{Speed:5.99 Deg:167} Clouds:{All:43} Rain:{OneH:0 ThreeH:0} Snow:{OneH:0 ThreeH:0} Dt:1733978331 ID:0 Name:Frisco Cod:200 Timezone:-21600 Unit:imperial Lang:EN Key:24abfbff3f37884f99d8b1e6ed41c998 Settings:0x14000116058} */

func (g *Game) formatWeatherString(c *owm.CurrentWeatherData) {
	g.weather = fmt.Sprintf("Status: %s %s\nTemp %f\nFeels Like: %f\n ", c.Weather[0].Main, c.Weather[0].Description, c.Main.Temp, c.Main.FeelsLike)
}
