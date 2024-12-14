package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

var (
	configPath = configdir.LocalConfig("Birdhouse")
	configFile = filepath.Join(configPath, "alarmclock.json")
)

type Config struct {
	WakeupStation  string
	SnoozeDuration int
	ClockFormat    string
	WeatherZipcode string
	WeatherCountry string
	OWM_API_key    string
	Alarms         map[alarmid]*Alarm
	EnabledAlarmID alarmid
}

func (g *Game) loadconfig() {
	err := configdir.MakePath(configPath)
	if err != nil {
		panic(err)
	}

	var config Config

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		config = Config{
			WakeupStation:  "Something",
			SnoozeDuration: 9,
			ClockFormat:    "3:04",
			WeatherZipcode: "75035",
			WeatherCountry: "US",
			EnabledAlarmID: disabledAlarmID,
			Alarms: map[alarmid]*Alarm{
				0: {AlarmTime: AlarmTime{7, 00}},
				1: {AlarmTime: AlarmTime{8, 00}},
				2: {AlarmTime: AlarmTime{5, 00}},
				3: {AlarmTime: AlarmTime{6, 30}},
				4: {AlarmTime: AlarmTime{4, 30}},
			},
		}
		fh, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		encoder := json.NewEncoder(fh)
		if err := encoder.Encode(&config); err != nil {
			panic(err)
		}
	} else {
		// Load the existing file.
		fh, err := os.Open(configFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		decoder := json.NewDecoder(fh)
		decoder.Decode(&config)
	}
	g.config = &config
}

func (g *Game) storeconfig() error {
	fh, err := os.Open(configFile)
	if err != nil {
		g.debug(err.Error())
		return err
	}
	defer fh.Close()

	encoder := json.NewEncoder(fh)
	if err := encoder.Encode(g.config); err != nil {
		g.debug(err.Error())
		return err
	}
	return nil
}
