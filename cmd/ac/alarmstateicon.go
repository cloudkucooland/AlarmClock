package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// alarmstateicon is just a sprite with a few extra methods
type alarmstateicon struct {
	*sprite
}

func (a *alarmstateicon) location() {
	if a == nil || a.sprite == nil {
		return
	}

	// #nosec G404 - we don't need strong randomness
	a.X = rand.Int()%screensize.X - 32
	// #nosec G404 - we don't need strong randomness
	a.Y = rand.Int()%screensize.Y - 32
}

func (g *Game) drawAlarmStateIcon(screen *ebiten.Image) {
	if g.config.EnabledAlarmID == disabledAlarmID {
		return // no alarm enabled, we are done
	}

	if g.alarmStateIcon == nil {
		g.setAlarmStateIcon()
		if g.alarmStateIcon == nil {
			return
		}
	}

	g.alarmStateIcon.drawWithLabel(screen)
}

func (g *Game) setAlarmStateIcon() {
	if g.config.EnabledAlarmID == disabledAlarmID {
		g.clearAlarmStateIcon()
		return // no alarm enabled, we are done
	}

	alarm, ok := g.config.Alarms[g.config.EnabledAlarmID]
	if !ok {
		g.clearAlarmStateIcon()
		return
	}

	g.alarmStateIcon = &alarmstateicon{}
	// default
	g.alarmStateIcon.sprite = getSprite("Swan Mommy", alarm.AlarmTime.String(), chirp)

	if alarm.station != nil && alarm.station.sprite != nil {
		g.alarmStateIcon.sprite.image = alarm.station.image
	}

	// random location
	g.alarmStateIcon.location()
}

func (g *Game) clearAlarmStateIcon() {
	g.alarmStateIcon = nil
}

func (a *alarmstateicon) draw(screen *ebiten.Image) {
	op := &colorm.DrawImageOptions{}
	op.Blend = ebiten.BlendCopy

	var cm colorm.ColorM
	cm.Scale(1.0, 1.0, 1.0, 0.25)

	op.GeoM.Scale(a.scale, a.scale)
	op.GeoM.Translate(float64(a.X), float64(a.Y))

	colorm.DrawImage(screen, a.image, cm, op)
}

func (a *alarmstateicon) drawWithLabel(screen *ebiten.Image) {
	a.draw(screen)
	a.drawlabel(screen) // from sprite
}
