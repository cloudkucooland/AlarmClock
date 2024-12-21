package sounds 

import (
	_ "embed"
)

var (
	//go:embed alarm.mp3
	BackupAlarm []byte
)

var Sounds = map[string][]byte{
	"BackupAlarm":     BackupAlarm,
}
