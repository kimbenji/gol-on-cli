package renderer

import (
	"fmt"
	"strings"

	"gol-on-cli/internal/engine"
)

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

func BuildFrame(board engine.Board, status StatusBarData) string {
	return BuildFrameWithPalette(board, status, Palette{})
}

func BuildFrameWithPalette(board engine.Board, status StatusBarData, palette Palette) string {
	var b strings.Builder
	aliveStart, deadStart, colorReset := colorSequences(palette)
	for y := 0; y < board.Height(); y++ {
		for x := 0; x < board.Width(); x++ {
			if board.IsAlive(x, y) {
				b.WriteString(aliveStart)
				b.WriteRune('█')
				b.WriteString(colorReset)
			} else {
				b.WriteString(deadStart)
				b.WriteRune(' ')
				b.WriteString(colorReset)
			}
		}
		b.WriteRune('\n')
	}
	b.WriteString(BuildStatusBar(status))
	b.WriteRune('\n')
	return b.String()
}

func colorSequences(palette Palette) (aliveStart, deadStart, reset string) {
	if palette.Mode == ModeTrueColor {
		return "\x1b[38;2;0;255;135m", "\x1b[38;2;31;41;55m", "\x1b[0m"
	}
	if palette.Mode == ModeFallback {
		return "\x1b[38;5;46m", "\x1b[38;5;236m", "\x1b[0m"
	}
	return "", "", ""
}
