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
