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
		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		encoder := json.NewEncoder(f)
		if err := encoder.Encode(&config); err != nil {
			panic(err)
		}
	} else {
		// Load the existing file.
		f, err := os.Open(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		decoder := json.NewDecoder(f)
		decoder.Decode(&config)
	}
	g.config = &config
}

func (g *Game) storeconfig() error {
	f, err := os.OpenFile(configFile, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		g.debug("unable to open file " + err.Error())
		return err
	}
	defer f.Close()
	// f.Seek(0, io.SeekStart)

	encoder := json.NewEncoder(f)
	if err := encoder.Encode(g.config); err != nil {
		g.debug("unable to write file: " + err.Error())
		return err
	}
	return nil
}
