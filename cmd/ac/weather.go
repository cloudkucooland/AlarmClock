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

func (g *Game) drawWeather(screen *ebiten.Image) {
	if g.weather == nil || len(g.weather.Weather) == 0 {
		return
	}

	weatherstring := fmt.Sprintf("Current %.1f Feels Like: %.1f High: %.1f Low: %.1f ", g.weather.Main.Temp, g.weather.Main.FeelsLike, g.weather.Main.TempMax, g.weather.Main.TempMin)

	textwidth, _ := text.Measure(weatherstring, weatherfont, 1)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(screensize.X/2)-float64(textwidth/2), float64(screensize.Y-40))
	op.LineSpacing = weatherfont.Size
	text.Draw(screen, weatherstring, weatherfont, op)
}

/* &{GeoPos:{Longitude:-96.7524 Latitude:33.1377} Sys:{Type:2 ID:2001388 Message:0 Country:US Sunrise:1733923242 Sunset:1733959236} Base:stations Weather:[{ID:802 Main:Clouds Description:scattered clouds Icon:03n}] Main:{Temp:44.67 TempMin:41.72 TempMax:47.44 FeelsLike:41.32 Pressure:1027 SeaLevel:1027 GrndLevel:1003 Humidity:48} Visibility:10000 Wind:{Speed:5.99 Deg:167} Clouds:{All:43} Rain:{OneH:0 ThreeH:0} Snow:{OneH:0 ThreeH:0} Dt:1733978331 ID:0 Name:Frisco Cod:200 Timezone:-21600 Unit:imperial Lang:EN Key:derp Settings:0x14000116058} */

func (g *Game) runWeather(ctx context.Context) {
	apikey := os.Getenv("OWM_API_KEY")
	if apikey == "" {
		err := fmt.Errorf("OWM_API_KEY not set; not running weather poller")
		fmt.Println(err.Error())
	}
	w, err := owm.NewCurrent("F", "EN", apikey)
	if err != nil {
		fmt.Println(err.Error())
	}
	g.weather = w

	w.CurrentByZipcode("75035", "US")

	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.CurrentByZipcode("75035", "US"); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
