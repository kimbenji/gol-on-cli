package renderer

import (
	"strings"
	"testing"

	"gol-on-cli/internal/engine"
)

func TestShouldRenderAliveAndDeadCellsWithDifferentColors(t *testing.T) {
	palette := SelectPalette(true)

	if palette.Alive == palette.Dead {
		t.Fatalf("expected alive and dead colors to be different")
	}
}

func TestShouldPreferTrueColorPaletteWhenTerminalSupportsIt(t *testing.T) {
	palette := SelectPalette(true)

	if palette.Mode != ModeTrueColor {
		t.Fatalf("expected truecolor mode, got %s", palette.Mode)
	}
}

func TestShouldUseFallbackPaletteWhenTrueColorIsNotSupported(t *testing.T) {
	palette := SelectPalette(false)

	if palette.Mode != ModeFallback {
		t.Fatalf("expected fallback mode, got %s", palette.Mode)
	}
}

func TestShouldShowGenerationPlayStateAndPatternSourceInStatusBar(t *testing.T) {
	status := BuildStatusBar(StatusBarData{Generation: 7, Paused: true, PatternSource: "wiki:glider"})

	assertContains(t, status, "gen:7")
	assertContains(t, status, "state:paused")
	assertContains(t, status, "source:wiki:glider")
}

func TestShouldAlwaysShowKeyboardShortcutsInStatusBar(t *testing.T) {
	status := BuildStatusBar(StatusBarData{Generation: 3, Paused: false, PatternSource: "random"})

	assertContains(t, status, "q")
	assertContains(t, status, "h/?")
	assertContains(t, status, "space")
	assertContains(t, status, "r")
	assertContains(t, status, "l")
}

func TestShouldBuildFrameWithBoardGridAndStatusBar(t *testing.T) {
	board := engine.NewBoard(3, 2)
	board.SetAlive(1, 0, true)
	board.SetAlive(2, 1, true)

	frame := BuildFrame(board, StatusBarData{Generation: 2, Paused: false, PatternSource: "random"})

	assertContains(t, frame, " █ ")
	assertContains(t, frame, "  █")
	assertContains(t, frame, "gen:2")
	assertContains(t, frame, "source:random")
}

func TestShouldColorizeFrameWhenPaletteIsProvided(t *testing.T) {
	board := engine.NewBoard(1, 1)
	board.SetAlive(0, 0, true)

	frame := BuildFrameWithPalette(board, StatusBarData{Generation: 0, PatternSource: "random"}, SelectPalette(true))

	assertContains(t, frame, "\x1b[38;2;255;215;0m")
	assertContains(t, frame, "\x1b[0m")
}

func TestShouldRenderNewbornAndRecentlyDeadCellsDynamically(t *testing.T) {
	previous := engine.NewBoard(2, 1)
	previous.SetAlive(1, 0, true)

	current := engine.NewBoard(2, 1)
	current.SetAlive(0, 0, true)

	frame := BuildFrameWithHistory(current, &previous, StatusBarData{Generation: 1, PatternSource: "random"}, SelectPalette(true))

	assertContains(t, frame, "\x1b[38;2;255;215;0m")
	assertContains(t, frame, "\x1b[38;2;255;99;71m")
	assertContains(t, frame, "\x1b[38;2;255;99;71m \x1b[0m")
}

func assertContains(t *testing.T, got, expected string) {
	t.Helper()
	if !strings.Contains(got, expected) {
		t.Fatalf("expected %q to contain %q", got, expected)
	}
}
