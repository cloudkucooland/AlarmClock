package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

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
		g.alarmStateIcon = &alarmstateicon{}
		g.alarmStateIcon.sprite = getSprite("Swan Mommy", "KERA", chirp)

		alarm, ok := g.config.Alarms[g.config.EnabledAlarmID]
		if !ok {
			g.debug("unable to get configured alarm")
			return
		}

		if alarm.station != nil && alarm.station.sprite != nil {
			g.alarmStateIcon.sprite.name = alarm.station.name
			g.alarmStateIcon.sprite.image = alarm.station.image
			g.alarmStateIcon.sprite.spritelabel = alarm.station.spritelabel
		}

		// random location
		g.alarmStateIcon.location()
	}

	g.alarmStateIcon.drawWithLabel(screen)
}
