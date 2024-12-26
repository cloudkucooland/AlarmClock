package ledserver

import (
	"fmt"
)

const Pipefile = "/var/run/ledserver/command.sock"

type LED struct{}

type CommandCode int

const (
	AllOn CommandCode = iota
	BackOn
	FrontOn
	Rainbow
)

type Command struct {
	Command    CommandCode
	Brightness uint8
	// RGB?
}

type Result bool

func (l *LED) Set(cmd *Command, res *Result) error {
	fmt.Printf("%+v", cmd)

	switch cmd.Command {
	case AllOn:
		// turn them all on
	case BackOn:
		// just the back
	case FrontOn:
		// just the front
	case Rainbow:
		// do a rainbow
	}
	return nil
}
