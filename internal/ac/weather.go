package ac

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	owm "github.com/briandowns/openweathermap"
)

func (g *Game) drawWeather(screen *ebiten.Image) {
	if g.weathercache == nil {
		return
	}

	b := g.weathercache.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(screensize.X/2)-float64(b.Max.X/2), float64(screensize.Y-b.Max.Y-16))
	screen.DrawImage(g.weathercache, op)
}

func (g *Game) runWeather(ctx context.Context) {
	if g.config.OWM_API_key == "" {
		g.debug("OWM API key not set; not running weather poller")
		return
	}
	w, err := owm.NewCurrent("F", "EN", g.config.OWM_API_key)
	if err != nil {
		g.debug(err.Error())
		return
	}
	g.weather = w

	if err := w.CurrentByZipcode(g.config.WeatherZipcode, g.config.WeatherCountry); err != nil {
		g.debug(err.Error())
	} else {
		updateweathercache(g)
	}

	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.CurrentByZipcode(g.config.WeatherZipcode, g.config.WeatherCountry); err != nil {
				g.debug(err.Error())
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
	line1 := fmt.Sprintf("Current %.0f°  •  Feels %.0f°", g.weather.Main.Temp, g.weather.Main.FeelsLike)
	line2 := fmt.Sprintf("High %.0f°  •  Low %.0f°", g.weather.Main.TempMax, g.weather.Main.TempMin)

	w1, h1 := text.Measure(line1, weatherfont, 0)
	w2, h2 := text.Measure(line2, weatherfont, 0)

	padX := float32(22)
	padY := float32(10)
	lineGap := float64(6)

	contentW := w1
	if w2 > contentW {
		contentW = w2
	}
	boxW := float32(contentW) + padX*2
	boxH := float32(h1+h2+lineGap) + padY*2

	if g.weathercache != nil {
		g.weathercache.Deallocate()
		g.weathercache = nil
	}

	g.weathercache = ebiten.NewImage(int(boxW), int(boxH))
	drawRoundedBox(g.weathercache, 0, 0, boxW, boxH, 16, modalSlate, rimLight, 1.5)

	drawCenteredWeatherLine(g.weathercache, line1, float64(boxW)/2.0, float64(padY), textPrimary)
	drawCenteredWeatherLine(g.weathercache, line2, float64(boxW)/2.0, float64(padY)+h1+lineGap, textMuted)
}

func drawCenteredWeatherLine(dst *ebiten.Image, str string, centerX, y float64, col color.RGBA) {
	w, _ := text.Measure(str, weatherfont, 0)
	x := centerX - w/2.0

	shadowOp := &text.DrawOptions{}
	shadowOp.GeoM.Translate(x, y+1)
	shadowOp.ColorScale.ScaleWithColor(textShadow)
	text.Draw(dst, str, weatherfont, shadowOp)

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(col)
	text.Draw(dst, str, weatherfont, op)
}
