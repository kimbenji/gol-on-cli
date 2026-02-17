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
	s.board = s.board.NextGeneration()
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
}

func (s *Simulation) LoadPatternFromWikiContent(content string) error {
	parsedBoard, err := pattern.LoadBoardFromWikiContent(content, s.width, s.height)
	if err != nil {
		return err
	}

	s.board = parsedBoard
	s.generation = 0
	return nil
}

func (s *Simulation) Generation() int {
	return s.generation
}

func (s *Simulation) Board() engine.Board {
	return s.board
}

func randomBoard(rng *rand.Rand, width, height int) engine.Board {
	board := engine.NewBoard(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rng.Intn(2) == 1 {
				board.SetAlive(x, y, true)
			}
		}
	}
	return board
}
