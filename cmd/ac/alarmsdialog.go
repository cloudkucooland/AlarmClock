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

		if a.dialogButton.Min.X == 0 {
			a.dialogButton = image.Rect(x, y, endx, y+int(rowheight))
		}
		vector.StrokeRect(screen, float32(a.dialogButton.Min.X), float32(a.dialogButton.Min.Y), float32(a.dialogButton.Max.X-a.dialogButton.Min.X), rowheight, float32(2), bordergrey, false)
		alarmtime := fmt.Sprintf("%0.2d:%0.2d", a.AlarmTime.Hour, a.AlarmTime.Minute)

		textwidth, textheight := text.Measure(alarmtime, weatherfont, 0)

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x+64), float64(y)+float64(rowheight/2.0)-float64(textheight/2))
		op.ColorScale.ScaleWithColor(color.Black)

		text.Draw(screen, alarmtime, weatherfont, op)

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
	return (x >= a.dialogButton.Min.X && x <= a.dialogButton.Max.X) && (y >= a.dialogButton.Min.Y && y <= a.dialogButton.Max.Y)
}
