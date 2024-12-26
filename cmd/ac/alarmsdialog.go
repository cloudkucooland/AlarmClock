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
	chirp(g)
	g.state = inAlarmConfig
}

func (g *Game) drawAlarmConfig(screen *ebiten.Image) {
	g.drawModal(screen)

	x := 46 // from modal
	y := 46 // from modal
	xpadding := 16
	hourControlOffsetX := 16
	endx := screensize.X - 130 // defined by modal
	rowheight := float32((screensize.Y - 79) / len(g.config.Alarms))

	// range over map doesn't always happen in the same order, causing chaos, this gives us a sorted list of alarmIDs to use
	for _, key := range slices.Sorted(maps.Keys(g.config.Alarms)) {
		a := g.config.Alarms[key]
		alarmtime := fmt.Sprintf("%0.2d:%0.2d", a.AlarmTime.Hour, a.AlarmTime.Minute)
		textwidth, textheight := text.Measure(alarmtime, weatherfont, 0)

		if a.dialogButton.bounds.Min.X == 0 { // the button is uninitialized
			inx := x
			a.dialogButton.bounds = image.Rect(inx, y, endx, y+int(rowheight))

			a.dialogButton.hourUp = getSprite("Up", "Hour Up", func(g *Game) {
				a.AlarmTime.Hour = (a.AlarmTime.Hour + 1) % 24
			})
			a.dialogButton.hourUp.scale = 1
			updownbounds := a.dialogButton.hourUp.image.Bounds()
			inx = x + hourControlOffsetX
			a.dialogButton.hourUp.setLocation(inx, y)

			a.dialogButton.hourDn = getSprite("Dn", "Hour Down", func(g *Game) {
				if a.AlarmTime.Hour <= 0 {
					a.AlarmTime.Hour = 23
				} else {
					a.AlarmTime.Hour = a.AlarmTime.Hour - 1
				}
			})
			a.dialogButton.hourDn.scale = 1
			a.dialogButton.hourDn.setLocation(inx, y+(int(rowheight/2)))

			a.dialogButton.minUp = getSprite("Up", "Minute Up", func(g *Game) {
				a.AlarmTime.Minute = (a.AlarmTime.Minute + 15) % 60
				// tick hour up if goes to 00?
			})
			a.dialogButton.minUp.scale = 1
			inx = inx + updownbounds.Max.X + xpadding + int(textwidth) + xpadding
			a.dialogButton.minUp.setLocation(inx, y)

			a.dialogButton.minDn = getSprite("Dn", "Minute Down", func(g *Game) {
				if a.AlarmTime.Minute <= 0 {
					a.AlarmTime.Minute = 45
					// a.AlarmTime.Hour = a.AlarmTime.Hour - 1
				} else {
					a.AlarmTime.Minute = (a.AlarmTime.Minute - 15) % 60
				}
			})
			a.dialogButton.minDn.scale = 1
			a.dialogButton.minDn.setLocation(inx, y+(int(rowheight/2)))
		}
		updownbounds := a.dialogButton.hourUp.image.Bounds()
		vector.StrokeRect(screen, float32(a.dialogButton.bounds.Min.X), float32(a.dialogButton.bounds.Min.Y), float32(a.dialogButton.bounds.Max.X-a.dialogButton.bounds.Min.X), rowheight, float32(2), bordergrey, false)

		a.dialogButton.hourUp.draw(screen)
		a.dialogButton.hourDn.draw(screen)

		inx := x + hourControlOffsetX + updownbounds.Max.X + xpadding

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(inx), float64(y)+float64(rowheight/2.0)-float64(textheight/2))
		op.ColorScale.ScaleWithColor(color.Black)
		text.Draw(screen, alarmtime, weatherfont, op)

		a.dialogButton.minUp.draw(screen)
		a.dialogButton.minDn.draw(screen)

		if key == g.config.EnabledAlarmID {
			inx = inx + xpadding + int(textwidth) + updownbounds.Max.X + xpadding
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(inx), float64(y)+float64(rowheight/2.0)-float64(textheight/2))
			op.ColorScale.ScaleWithColor(pink)
			text.Draw(screen, "Enabled", weatherfont, op)
		}

		y = y + int(rowheight) // + 1
	}
}

func (a Alarm) in(x int, y int) bool {
	return (x >= a.dialogButton.bounds.Min.X && x <= a.dialogButton.bounds.Max.X) && (y >= a.dialogButton.bounds.Min.Y && y <= a.dialogButton.bounds.Max.Y)
}
