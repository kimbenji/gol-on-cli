package app

import "testing"

func TestShouldProvideDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.FPS <= 0 {
		t.Fatalf("expected default FPS to be positive, got %d", cfg.FPS)
	}

	if cfg.SeedMode != SeedModeRandom {
		t.Fatalf("expected default seed mode to be %q, got %q", SeedModeRandom, cfg.SeedMode)
	}

	if cfg.ColorMode != ColorModeTrueColor {
		t.Fatalf("expected default color mode to be %q, got %q", ColorModeTrueColor, cfg.ColorMode)
	}
}
