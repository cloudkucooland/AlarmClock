package main

import (
	"image"
	"maps"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

	if len(g.config.Alarms) == 0 {
		return
	}

	mx, my, mw, mh := modalBounds()
	closeBtn := modalbuttons[0]
	closeW := int(float64(closeBtn.image.Bounds().Max.X) * closeBtn.scale)

	x := mx + modalPadding
	y := my + modalPadding
	endx := mx + mw - closeW - modalPadding*2
	innerH := mh - modalPadding*2
	rowGap := 10
	rowHeight := (innerH - rowGap*(len(g.config.Alarms)-1)) / len(g.config.Alarms)

	arrowScale := 1.25
	arrowGap := 6
	xPadding := 18

	// Sort alarm IDs so row order stays deterministic across frames
	for _, key := range slices.Sorted(maps.Keys(g.config.Alarms)) {
		a := g.config.Alarms[key]
		alarmtime := a.AlarmTime.String()
		textwidth, textheight := text.Measure(alarmtime, weatherfont, 0)

		if a.dialogButton.bounds.Min.X == 0 {
			a.dialogButton.bounds = image.Rect(x, y, endx, y+rowHeight)

			a.dialogButton.hourUp = getSprite("Up", "Hour Up", func(g *Game) {
				a.AlarmTime.Hour = (a.AlarmTime.Hour + 1) % 24
				tick(g)
			})
			a.dialogButton.hourUp.scale = arrowScale

			a.dialogButton.hourDn = getSprite("Dn", "Hour Down", func(g *Game) {
				if a.AlarmTime.Hour <= 0 {
					a.AlarmTime.Hour = 23
				} else {
					a.AlarmTime.Hour = a.AlarmTime.Hour - 1
				}
				tick(g)
			})
			a.dialogButton.hourDn.scale = arrowScale

			a.dialogButton.minUp = getSprite("Up", "Minute Up", func(g *Game) {
				a.AlarmTime.Minute = (a.AlarmTime.Minute + 15) % 60
				tick(g)
			})
			a.dialogButton.minUp.scale = arrowScale

			a.dialogButton.minDn = getSprite("Dn", "Minute Down", func(g *Game) {
				if a.AlarmTime.Minute <= 0 {
					a.AlarmTime.Minute = 45
				} else {
					a.AlarmTime.Minute = (a.AlarmTime.Minute - 15) % 60
				}
				tick(g)
			})
			a.dialogButton.minDn.scale = arrowScale

			// Position arrow pairs centered vertically inside the row card
			rawB := a.dialogButton.hourUp.image.Bounds()
			arrowW := int(float64(rawB.Max.X) * arrowScale)
			arrowH := int(float64(rawB.Max.Y) * arrowScale)
			pairH := arrowH*2 + arrowGap
			topArrowY := y + (rowHeight-pairH)/2
			botArrowY := topArrowY + arrowH + arrowGap

			hourX := x + xPadding
			a.dialogButton.hourUp.setLocation(hourX, topArrowY)
			a.dialogButton.hourDn.setLocation(hourX, botArrowY)

			minX := hourX + arrowW + xPadding + int(textwidth) + xPadding
			a.dialogButton.minUp.setLocation(minX, topArrowY)
			a.dialogButton.minDn.setLocation(minX, botArrowY)
		}

		isEnabled := key == g.config.EnabledAlarmID
		rowW := float32(a.dialogButton.bounds.Max.X - a.dialogButton.bounds.Min.X)

		if isEnabled {
			drawRoundedBox(screen, float32(a.dialogButton.bounds.Min.X), float32(a.dialogButton.bounds.Min.Y), rowW, float32(rowHeight), 14, pinkTint, pink, 2.0)
		} else {
			drawRoundedBox(screen, float32(a.dialogButton.bounds.Min.X), float32(a.dialogButton.bounds.Min.Y), rowW, float32(rowHeight), 14, cardSlate, bordergrey, 1.0)
		}

		a.dialogButton.hourUp.draw(screen)
		a.dialogButton.hourDn.draw(screen)

		rawB := a.dialogButton.hourUp.image.Bounds()
		arrowW := int(float64(rawB.Max.X) * arrowScale)
		timeX := x + xPadding + arrowW + xPadding
		timeY := float64(y) + float64(rowHeight)/2.0 - textheight/2.0

		// Subtle drop shadow + foreground time text
		shadowOp := &text.DrawOptions{}
		shadowOp.GeoM.Translate(float64(timeX), timeY+1)
		shadowOp.ColorScale.ScaleWithColor(textShadow)
		text.Draw(screen, alarmtime, weatherfont, shadowOp)

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(timeX), timeY)
		if isEnabled {
			op.ColorScale.ScaleWithColor(textPrimary)
		} else {
			op.ColorScale.ScaleWithColor(textMuted)
		}
		text.Draw(screen, alarmtime, weatherfont, op)

		a.dialogButton.minUp.draw(screen)
		a.dialogButton.minDn.draw(screen)

		// Status pill badge on the right side of the row card
		badgeText := "Off"
		badgeBg := spritelabelgrey
		if isEnabled {
			badgeText = "Enabled"
			if a.station != nil {
				badgeText = a.station.spritelabel.label
			}
			badgeBg = pink
		}

		badgeImg := button(badgeText, badgeBg, controlfont)
		bb := badgeImg.Bounds()
		badgeOp := &ebiten.DrawImageOptions{}
		badgeOp.GeoM.Translate(
			float64(a.dialogButton.bounds.Max.X-bb.Max.X-xPadding),
			float64(y)+float64(rowHeight-bb.Max.Y)/2.0,
		)
		screen.DrawImage(badgeImg, badgeOp)
		badgeImg.Deallocate()

		y += rowHeight + rowGap
	}
}

func (a Alarm) in(x int, y int) bool {
	return (x >= a.dialogButton.bounds.Min.X && x <= a.dialogButton.bounds.Max.X) && (y >= a.dialogButton.bounds.Min.Y && y <= a.dialogButton.bounds.Max.Y)
}
