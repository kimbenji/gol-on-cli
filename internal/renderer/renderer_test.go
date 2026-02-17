package renderer

import (
	"strings"
	"testing"
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

func assertContains(t *testing.T, got, expected string) {
	t.Helper()
	if !strings.Contains(got, expected) {
		t.Fatalf("expected %q to contain %q", got, expected)
	}
}
