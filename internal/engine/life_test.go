package engine

import "testing"

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
