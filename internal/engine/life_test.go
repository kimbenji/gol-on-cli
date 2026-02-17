package engine

import "testing"

func TestShouldReturnErrorWhenBoardSizeIsNonPositive(t *testing.T) {
	if _, err := NewBoardValidated(0, 3); err == nil {
		t.Fatalf("expected error for non-positive width")
	}

	if _, err := NewBoardValidated(3, 0); err == nil {
		t.Fatalf("expected error for non-positive height")
	}
}

func TestShouldIgnoreOutOfRangeCoordinatesWithoutPanicking(t *testing.T) {
	board := NewBoard(2, 2)

	board.SetAlive(-1, 0, true)
	board.SetAlive(0, -1, true)
	board.SetAlive(2, 0, true)
	board.SetAlive(0, 2, true)

	if board.IsAlive(-1, 0) {
		t.Fatalf("expected out-of-range read to be treated as dead cell")
	}
}

func TestShouldBirthDeadCellWithThreeNeighbors(t *testing.T) {
	board := NewBoard(3, 3)
	board.SetAlive(0, 1, true)
	board.SetAlive(1, 0, true)
	board.SetAlive(2, 1, true)

	next := board.NextGeneration()
	if !next.IsAlive(1, 1) {
		t.Fatalf("expected dead cell to become alive with exactly three neighbors")
	}
}

func TestShouldSurviveLiveCellWithTwoNeighbors(t *testing.T) {
	board := NewBoard(3, 3)
	board.SetAlive(1, 1, true)
	board.SetAlive(0, 1, true)
	board.SetAlive(2, 1, true)

	next := board.NextGeneration()
	if !next.IsAlive(1, 1) {
		t.Fatalf("expected live cell to survive with exactly two neighbors")
	}
}

func TestShouldSurviveLiveCellWithThreeNeighbors(t *testing.T) {
	board := NewBoard(3, 3)
	board.SetAlive(1, 1, true)
	board.SetAlive(0, 1, true)
	board.SetAlive(2, 1, true)
	board.SetAlive(1, 0, true)

	next := board.NextGeneration()
	if !next.IsAlive(1, 1) {
		t.Fatalf("expected live cell to survive with exactly three neighbors")
	}
}

func TestShouldDieLiveCellWithOneNeighbor(t *testing.T) {
	board := NewBoard(3, 3)
	board.SetAlive(1, 1, true)
	board.SetAlive(0, 1, true)

	next := board.NextGeneration()
	if next.IsAlive(1, 1) {
		t.Fatalf("expected live cell to die with fewer than two neighbors")
	}
}

func TestShouldDieLiveCellWithFourNeighbors(t *testing.T) {
	board := NewBoard(3, 3)
	board.SetAlive(1, 1, true)
	board.SetAlive(0, 1, true)
	board.SetAlive(1, 0, true)
	board.SetAlive(2, 1, true)
	board.SetAlive(1, 2, true)

	next := board.NextGeneration()
	if next.IsAlive(1, 1) {
		t.Fatalf("expected live cell to die with more than three neighbors")
	}
}

func TestShouldUpdateAllCellsSimultaneously(t *testing.T) {
	board := NewBoard(5, 5)
	board.SetAlive(1, 2, true)
	board.SetAlive(2, 2, true)
	board.SetAlive(3, 2, true)

	next := board.NextGeneration()

	if next.IsAlive(1, 2) || next.IsAlive(3, 2) {
		t.Fatalf("expected horizontal edge cells to die after simultaneous update")
	}
	if !next.IsAlive(2, 1) || !next.IsAlive(2, 2) || !next.IsAlive(2, 3) {
		t.Fatalf("expected blinker to rotate to a vertical line after one generation")
	}
}

func TestShouldCountNeighborsWithToroidalWrapping(t *testing.T) {
	board := NewBoard(3, 3)
	board.SetAlive(0, 0, true)
	board.SetAlive(2, 0, true)
	board.SetAlive(0, 2, true)

	next := board.NextGeneration()

	if !next.IsAlive(2, 2) {
		t.Fatalf("expected corner cell to be born from wrapped neighbors across board edges")
	}
}

func TestShouldSafelyHandleNextGenerationOnZeroSizedBoard(t *testing.T) {
	board := NewBoard(0, 0)
	next := board.NextGeneration()

	if next.Width() != 0 || next.Height() != 0 {
		t.Fatalf("expected zero-sized board to stay zero-sized, got %dx%d", next.Width(), next.Height())
	}
}
