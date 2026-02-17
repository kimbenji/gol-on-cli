package engine

import "fmt"

type Board struct {
	width  int
	height int
	cells  [][]bool
}

func NewBoardValidated(width, height int) (Board, error) {
	if width <= 0 || height <= 0 {
		return Board{}, fmt.Errorf("invalid board size: width and height must be greater than zero")
	}
	return NewBoard(width, height), nil
}

func NewBoard(width, height int) Board {
	cells := make([][]bool, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]bool, width)
	}
	return Board{width: width, height: height, cells: cells}
}

func (b *Board) SetAlive(x, y int, alive bool) {
	if !b.inBounds(x, y) {
		return
	}
	b.cells[y][x] = alive
}

func (b Board) IsAlive(x, y int) bool {
	if !b.inBounds(x, y) {
		return false
	}
	return b.cells[y][x]
}

func (b Board) inBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < b.width && y < b.height
}

func (b Board) Width() int {
	return b.width
}

func (b Board) Height() int {
	return b.height
}

func (b Board) NextGeneration() Board {
	next := NewBoard(b.width, b.height)
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			neighbors := b.aliveNeighbors(x, y)
			if !b.IsAlive(x, y) && neighbors == 3 {
				next.SetAlive(x, y, true)
			}
			if b.IsAlive(x, y) && (neighbors == 2 || neighbors == 3) {
				next.SetAlive(x, y, true)
			}
		}
	}
	return next
}

func (b Board) aliveNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx := (x + dx + b.width) % b.width
			ny := (y + dy + b.height) % b.height
			if b.IsAlive(nx, ny) {
				count++
			}
		}
	}
	return count
}
