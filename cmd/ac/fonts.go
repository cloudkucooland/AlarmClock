package main

import (
	"bytes"

	// "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/cloudkucooland/AlarmClock/resources"
)

var (
	spaceMonoSource *text.GoTextFaceSource
	clockfont       *text.GoTextFace
	controlfont     *text.GoTextFace
	bigbuttonfont   *text.GoTextFace
	weatherfont     *text.GoTextFace
)

func loadfonts() error {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(resources.SpaceMonoBold_ttf))
	if err != nil {
		return err
	}
	spaceMonoSource = s

	clockfont = &text.GoTextFace{
		Source: spaceMonoSource,
		Size:   192,
	}
	controlfont = &text.GoTextFace{
		Source: spaceMonoSource,
		Size:   16,
	}
	bigbuttonfont = &text.GoTextFace{
		Source: spaceMonoSource,
		Size:   96,
	}
	weatherfont = &text.GoTextFace{
		Source: spaceMonoSource,
		Size:   24,
	}
	return nil
}
