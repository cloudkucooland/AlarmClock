package sounds 

import (
	_ "embed"
)

var (
	//go:embed alarm.mp3
	BackupAlarm []byte

	//go:embed khew.mp3
	Khew []byte
)

var Sounds = map[string][]byte{
	"BackupAlarm":     BackupAlarm,
	"Khew":     Khew,
}
