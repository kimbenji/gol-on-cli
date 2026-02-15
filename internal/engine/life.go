package engine

type Board struct {
	width  int
	height int
	cells  [][]bool
}

func NewBoard(width, height int) Board {
	cells := make([][]bool, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]bool, width)
	}
	return Board{width: width, height: height, cells: cells}
}

func (b *Board) SetAlive(x, y int, alive bool) {
	b.cells[y][x] = alive
}

func (b Board) IsAlive(x, y int) bool {
	return b.cells[y][x]
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
			nx := x + dx
			ny := y + dy
			if nx < 0 || nx >= b.width || ny < 0 || ny >= b.height {
				continue
			}
			if b.IsAlive(nx, ny) {
				count++
			}
		}
	}
	return count
}
