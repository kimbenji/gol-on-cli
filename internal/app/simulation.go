package app

import (
	"math/rand"

	"gol-on-cli/internal/engine"
	"gol-on-cli/internal/pattern"
)

type BoardFactory func(width, height int) engine.Board

type Simulation struct {
	board        engine.Board
	generation   int
	stableGenerations int
	paused       bool
	width        int
	height       int
	boardFactory BoardFactory
}

func NewSimulation(width, height int, seed int64) *Simulation {
	rng := rand.New(rand.NewSource(seed))
	factory := func(w, h int) engine.Board {
		return randomBoard(rng, w, h)
	}
	return NewSimulationWithFactory(width, height, factory)
}

func NewSimulationWithFactory(width, height int, factory BoardFactory) *Simulation {
	return &Simulation{
		board:        factory(width, height),
		generation:   0,
		stableGenerations: 0,
		paused:       false,
		width:        width,
		height:       height,
		boardFactory: factory,
	}
}

func (s *Simulation) Tick() {
	if s.paused {
		return
	}
	next := s.board.NextGeneration()
	if boardsMatch(s.board, next) {
		s.stableGenerations++
	} else {
		s.stableGenerations = 0
	}
	if s.stableGenerations >= 100 {
		s.Restart()
		s.stableGenerations = 0
		return
	}
	s.board = next
	s.generation++
}

func (s *Simulation) Pause() {
	s.paused = true
}

func (s *Simulation) Resume() {
	s.paused = false
}

func (s *Simulation) Restart() {
	s.board = s.boardFactory(s.width, s.height)
	s.generation = 0
	s.stableGenerations = 0
}

func (s *Simulation) LoadPatternFromWikiContent(content string) error {
	parsedBoard, err := pattern.LoadBoardFromWikiContent(content, s.width, s.height)
	if err != nil {
		return err
	}

	s.board = parsedBoard
	s.generation = 0
	s.stableGenerations = 0
	return nil
}

func (s *Simulation) Generation() int {
	return s.generation
}

func (s *Simulation) Board() engine.Board {
	return s.board
}

func (s *Simulation) Resize(width, height int) {
	resized := engine.NewBoard(width, height)
	copyWidth := min(s.board.Width(), width)
	copyHeight := min(s.board.Height(), height)
	for y := 0; y < copyHeight; y++ {
		for x := 0; x < copyWidth; x++ {
			if s.board.IsAlive(x, y) {
				resized.SetAlive(x, y, true)
			}
		}
	}
	s.board = resized
	s.width = width
	s.height = height
	s.stableGenerations = 0
}

func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}
func randomBoard(rng *rand.Rand, width, height int) engine.Board {
	board := engine.NewBoard(width, height)
	if width <= 0 || height <= 0 {
		return board
	}

	windowW := width
	if windowW > 10 {
		windowW = 10
	}
	windowH := height
	if windowH > 10 {
		windowH = 10
	}

	startX := (width - windowW) / 2
	startY := (height - windowH) / 2

	maxCells := windowW * windowH
	target := 16
	if maxCells < target {
		target = maxCells
	}

	picked := make(map[int]struct{}, target)
	for len(picked) < target {
		index := rng.Intn(maxCells)
		if _, exists := picked[index]; exists {
			continue
		}
		picked[index] = struct{}{}
		x := index % windowW
		y := index / windowW
		board.SetAlive(startX+x, startY+y, true)
	}
	return board
}

func boardsMatch(left, right engine.Board) bool {
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
