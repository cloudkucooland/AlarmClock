package main

import (
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/sysfs"
)

const threshold float64 = 69.5

func thermal() (float64, error) {
	var max float64

	if _, err := host.Init(); err != nil {
		return max, err
	}

	sensors := sysfs.ThermalSensors

	for _, t := range sensors {
		e := physic.Env{}
		if err := t.Sense(&e); err != nil {
			return max, err
		}

		if e.Temperature.Celsius() > max {
			max = e.Temperature.Celsius()
		}
	}
	return max, nil
}
