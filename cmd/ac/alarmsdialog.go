package main

import (
	"fmt"
	"image"
	"image/color"
	"maps"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type alarmDialogButton struct {
	bounds image.Rectangle
	hourUp *sprite
	hourDn *sprite
	minUp  *sprite
	minDn  *sprite
}

func alarmConfigDialog(g *Game) {
	g.state = inAlarmConfig
}

func (g *Game) drawAlarmConfig(screen *ebiten.Image) {
	g.drawModal(screen)

	x := 32                    // from modal
	y := 32                    // from modal
	endx := screensize.X - 164 // defined by modal
	rowheight := float32((screensize.Y - 64) / len(g.config.Alarms))

	// range over map doesn't always happen in the same order, causing chaos, this gives us a sorted list of alarmIDs to use
	for _, key := range slices.Sorted(maps.Keys(g.config.Alarms)) {
		a := g.config.Alarms[key]
		alarmtime := fmt.Sprintf("%0.2d:%0.2d", a.AlarmTime.Hour, a.AlarmTime.Minute)
		textwidth, textheight := text.Measure(alarmtime, weatherfont, 0)

		if a.dialogButton.bounds.Min.X == 0 { // uninitialized
			a.dialogButton.bounds = image.Rect(x, y, endx, y+int(rowheight))

			a.dialogButton.hourUp = getSprite("Up", "Hour Up", func(g *Game) {
				g.debug("hour up")
				a.AlarmTime.Hour = (a.AlarmTime.Hour + 1) % 24
			})
			a.dialogButton.hourUp.scale = 1
			a.dialogButton.hourUp.setLocation(x+16, y-16+(int(rowheight/2)))

			a.dialogButton.hourDn = getSprite("Dn", "Hour Down", func(g *Game) {
				g.debug("hour dn")
				a.AlarmTime.Hour = (a.AlarmTime.Hour - 1) % 24
			})
			a.dialogButton.hourDn.scale = 1
			a.dialogButton.hourDn.setLocation(x+16, y+(int(rowheight/2)))

			a.dialogButton.minUp = getSprite("Up", "Minute Up", func(g *Game) {
				g.debug("min up")
				a.AlarmTime.Minute = (a.AlarmTime.Minute + 15) % 60
			})
			a.dialogButton.minUp.scale = 1
			a.dialogButton.minUp.setLocation(x+78+int(textwidth), y-16+(int(rowheight/2)))

			a.dialogButton.minDn = getSprite("Dn", "Minute Down", func(g *Game) {
				g.debug("min dn")
				a.AlarmTime.Minute = (a.AlarmTime.Minute - 15) % 60
			})
			a.dialogButton.minDn.scale = 1
			a.dialogButton.minDn.setLocation(x+78+int(textwidth), y+(int(rowheight/2)))

		}
		vector.StrokeRect(screen, float32(a.dialogButton.bounds.Min.X), float32(a.dialogButton.bounds.Min.Y), float32(a.dialogButton.bounds.Max.X-a.dialogButton.bounds.Min.X), rowheight, float32(2), bordergrey, false)

		a.dialogButton.hourUp.draw(screen)
		a.dialogButton.hourDn.draw(screen)
		a.dialogButton.minUp.draw(screen)
		a.dialogButton.minDn.draw(screen)

		{
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(x+64), float64(y)+float64(rowheight/2.0)-float64(textheight/2))
			op.ColorScale.ScaleWithColor(color.Black)
			text.Draw(screen, alarmtime, weatherfont, op)
		}

		if key == g.config.EnabledAlarmID {
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(x+128)+textwidth, float64(y)+float64(rowheight/2.0)-float64(textheight/2))
			op.ColorScale.ScaleWithColor(pink)
			text.Draw(screen, "Enabled", weatherfont, op)
		}

		y = y + int(rowheight) + 1
	}
}

func (a Alarm) in(x int, y int) bool {
	return (x >= a.dialogButton.bounds.Min.X && x <= a.dialogButton.bounds.Max.X) && (y >= a.dialogButton.bounds.Min.Y && y <= a.dialogButton.bounds.Max.Y)
}
