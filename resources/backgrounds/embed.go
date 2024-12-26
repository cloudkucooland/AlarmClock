package backgrounds

import (
	_ "embed"
)

var (
	//go:embed flowersplash.png
	Default []byte

	//go:embed hummingbird_bg.png
	Hummingbird []byte

	//go:embed owleyes.png
	Owleyes []byte

	//go:embed owlmoon.png
	Owlmoon []byte
)
