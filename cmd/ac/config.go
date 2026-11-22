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

type SpecialDay struct {
    Month   int    `json:"month"`   // 1–12
    Day     int    `json:"day"`     // 1–31
    Title   string `json:"title"`   // e.g., "Happy Anniversary!", "Merry Christmas!"
    Sprite  string `json:"sprite"`  // optional bird sprite override, e.g., "Party" or "Spring"
    LEDMode string `json:"ledmode"` // optional: "rainbow", "festive", etc.
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
		// #nosec G304
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
		// #nosec G304
		f, err := os.Open(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		decoder := json.NewDecoder(f)
		if err := decoder.Decode(&config); err != nil {
			panic(err)
		}
	}
	g.config = &config
}

func (g *Game) storeconfig() error {
	tmpFile := configFile + ".tmp"
	// #nosec G304
	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		g.debug("unable to open temp config: " + err.Error())
		return err
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(g.config); err != nil {
		_ = f.Close()
		g.debug("unable to write temp config: " + err.Error())
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpFile, configFile); err != nil {
		g.debug("unable to rename config: " + err.Error())
		return err
	}
	return nil
}
