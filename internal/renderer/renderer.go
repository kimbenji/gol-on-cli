package renderer

import "fmt"

type PaletteMode string

const (
	ModeTrueColor PaletteMode = "truecolor"
	ModeFallback  PaletteMode = "fallback"
)

type Palette struct {
	Mode  PaletteMode
	Alive string
	Dead  string
}

type StatusBarData struct {
	Generation    int
	Paused        bool
	PatternSource string
}

func SelectPalette(supportsTrueColor bool) Palette {
	if supportsTrueColor {
		return Palette{Mode: ModeTrueColor, Alive: "#00FF87", Dead: "#1F2937"}
	}
	return Palette{Mode: ModeFallback, Alive: "46", Dead: "236"}
}

func BuildStatusBar(data StatusBarData) string {
	state := "running"
	if data.Paused {
		state = "paused"
	}
	return fmt.Sprintf(
		"gen:%d | state:%s | source:%s | keys:q h/? space r l",
		data.Generation,
		state,
		data.PatternSource,
	)
}
