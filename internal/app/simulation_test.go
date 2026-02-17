package app

import (
	"testing"

	"gol-on-cli/internal/engine"
	"gol-on-cli/internal/pattern"
)

func TestShouldInitializeRandomBoardWithRequestedSize(t *testing.T) {
	sim := NewSimulation(12, 8, 42)

	if sim.Board().Width() != 12 || sim.Board().Height() != 8 {
		t.Fatalf("expected board size 12x8, got %dx%d", sim.Board().Width(), sim.Board().Height())
	}
}

func TestShouldCreateSameInitialBoardForSameSeed(t *testing.T) {
	left := NewSimulation(10, 6, 1234)
	right := NewSimulation(10, 6, 1234)

	if !boardsEqual(left.Board(), right.Board()) {
		t.Fatalf("expected same initial board for same seed")
	}
}

func TestShouldIncreaseGenerationByOneAfterSingleTick(t *testing.T) {
	sim := NewSimulation(5, 5, 1)

	sim.Tick()

	if sim.Generation() != 1 {
		t.Fatalf("expected generation to increase by 1 after tick, got %d", sim.Generation())
	}
}

func TestShouldNotIncreaseGenerationWhenPaused(t *testing.T) {
	sim := NewSimulation(5, 5, 1)
	sim.Pause()

	sim.Tick()

	if sim.Generation() != 0 {
		t.Fatalf("expected generation to remain unchanged while paused, got %d", sim.Generation())
	}
}

func TestShouldResumeGenerationProgressAfterResume(t *testing.T) {
	sim := NewSimulation(5, 5, 1)
	sim.Pause()
	sim.Tick()
	sim.Resume()

	sim.Tick()

	if sim.Generation() != 1 {
		t.Fatalf("expected generation to increase after resume, got %d", sim.Generation())
	}
}

func TestShouldResetGenerationAndReinitializeBoardOnRestart(t *testing.T) {
	first := engine.NewBoard(3, 3)
	first.SetAlive(0, 0, true)
	second := engine.NewBoard(3, 3)
	second.SetAlive(2, 2, true)

	calls := 0
	factory := func(width, height int) engine.Board {
		calls++
		if calls == 1 {
			return first
		}
		return second
	}

	sim := NewSimulationWithFactory(3, 3, factory)
	sim.Tick()
	if sim.Generation() != 1 {
		t.Fatalf("expected generation to be 1 after one tick, got %d", sim.Generation())
	}

	sim.Restart()

	if sim.Generation() != 0 {
		t.Fatalf("expected generation reset to 0 on restart, got %d", sim.Generation())
	}
	if calls != 2 {
		t.Fatalf("expected restart to reinitialize random board, factory calls=%d", calls)
	}
	if !sim.Board().IsAlive(2, 2) {
		t.Fatalf("expected board to be replaced with restarted random initialization")
	}
}

func boardsEqual(left, right engine.Board) bool {
	if left.Width() != right.Width() || left.Height() != right.Height() {
		return false
	}

	for y := 0; y < left.Height(); y++ {
		for x := 0; x < left.Width(); x++ {
			if left.IsAlive(x, y) != right.IsAlive(x, y) {
				return false
			}
		}
	}
	return true
}

func TestShouldKeepCurrentBoardAndReturnRecoverableErrorWhenPatternParsingFails(t *testing.T) {
	board := engine.NewBoard(3, 3)
	board.SetAlive(1, 1, true)

	sim := NewSimulationWithFactory(3, 3, func(width, height int) engine.Board {
		return board
	})
	before := sim.Board()

	err := sim.LoadPatternFromWikiContent("#Life 1.06\ninvalid")

	if err == nil {
		t.Fatalf("expected recoverable error when parsing fails")
	}
	if _, ok := err.(pattern.RecoverableError); !ok {
		t.Fatalf("expected recoverable error type, got %T", err)
	}
	if !boardsEqual(before, sim.Board()) {
		t.Fatalf("expected board to remain unchanged on parsing failure")
	}
}
