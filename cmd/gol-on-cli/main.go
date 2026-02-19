package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gol-on-cli/internal/app"
	"gol-on-cli/internal/cli"
	"gol-on-cli/internal/engine"
	"gol-on-cli/internal/input"
	"gol-on-cli/internal/pattern"
	"gol-on-cli/internal/renderer"

	"github.com/gdamore/tcell/v2"
)

const version = "v0.1.0"
const startupPatternTimeout = 5 * time.Second
const startupPatternMaxSize int64 = 1024 * 1024
const frameMarginCols = 6
const frameMarginRows = 4

type noopLoader struct{}

func (n noopLoader) Load(url string) error { return nil }

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("gol-on-cli", flag.ContinueOnError)
	flags.SetOutput(stderr)

	help := flags.Bool("help", false, "show usage")
	showVersion := flags.Bool("version", false, "show version")
	fps := flags.Int("fps", 5, "updates per second")
	seed := flags.Int64("seed", 0, "random seed")
	patternURL := flags.String("pattern-url", "", "startup pattern URL")
	flags.String("alive-color", "", "alive cell color")
	flags.String("dead-color", "", "dead cell color")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if *help {
		fmt.Fprintln(stdout, cli.BuildHelpText())
		return 0
	}
	if *showVersion {
		fmt.Fprintln(stdout, cli.BuildVersionText(version))
		return 0
	}

	if _, err := cli.Start(cli.StartOptions{PatternURL: *patternURL, FPS: *fps}, noopLoader{}); err != nil {
		fmt.Fprintf(stderr, "failed to start: %v\n", err)
		return 1
	}

	source := patternSource(*patternURL)
	if !isTerminal(stdout) {
		sim := app.NewSimulation(20, 10, *seed)
		status := renderer.BuildStatusBar(renderer.StatusBarData{Generation: sim.Generation(), Paused: false, PatternSource: source})
		fmt.Fprintln(stdout, status)
		return 0
	}

	fileIn, inOK := stdin.(*os.File)
	_, outOK := stdout.(*os.File)
	if !inOK || !outOK {
		fmt.Fprintln(stderr, "failed to start: interactive mode requires file stdin/stdout")
		return 1
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(stderr, "failed to start: %v\n", err)
		return 1
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(stderr, "failed to start: %v\n", err)
		return 1
	}
	defer screen.Fini()

	w, h := boardSizeForScreen(screen)
	sim := app.NewSimulation(w, h, *seed)
	if *patternURL != "" {
		if err := tryLoadPatternForSimulation(sim, *patternURL); err != nil {
			fmt.Fprintf(stderr, "failed to load startup pattern: %v\n", err)
		}
	}
	_ = fileIn
	return runFullscreen(screen, sim, *fps, source, *patternURL)
}

func runFullscreen(screen tcell.Screen, sim *app.Simulation, fps int, source string, patternURL string) int {
	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	state := input.NewState()
	palette := renderer.SelectPalette(supportsTrueColor())

	var previous *engine.Board
	needsFullClear := true
	notice := ""
	helpVisible := false
	var transient map[cellCoord]struct{}
	dirty := true

	eventCh := make(chan tcell.Event, 16)
	go func() {
		defer close(eventCh)
		for {
			ev := screen.PollEvent()
			if ev == nil {
				return
			}
			eventCh <- ev
		}
	}()

	fitSimulationToScreen(screen, sim)
	for {
		if state.HelpVisible != helpVisible {
			helpVisible = state.HelpVisible
			needsFullClear = true
			dirty = true
		}

		if state.ConsumeLoadPatternRequest() {
			if patternURL == "" {
				notice = "no-pattern-url-configured"
			} else if err := tryLoadPatternForSimulation(sim, patternURL); err != nil {
				notice = fmt.Sprintf("pattern-load-failed: %v", err)
			} else {
				notice = "pattern-loaded"
				previous = nil
			}
			needsFullClear = true
			dirty = true
		}

		if dirty {
			current := sim.Board()
			frameNotice := notice
			if state.HelpVisible {
				if frameNotice != "" {
					frameNotice += " | "
				}
				frameNotice += "help:q h/? space r l"
			}
			status := renderer.BuildStatusBar(renderer.StatusBarData{
				Generation:    sim.Generation(),
				Paused:        state.Paused,
				PatternSource: source,
				Notice:        frameNotice,
			})
			if needsFullClear {
				screen.Clear()
				renderBoardFull(screen, current, previous, palette)
				renderStatusBar(screen, current.Height(), status)
				screen.Show()
				needsFullClear = false
				transient = nil
			} else {
				updates, nextTransient := diffCells(current, previous, transient)
				renderCellUpdates(screen, updates, current, previous, palette)
				transient = nextTransient
				renderStatusBar(screen, current.Height(), status)
				screen.Show()
			}
			previousSnapshot := current
			previous = &previousSnapshot
			dirty = false
		}

		select {
		case <-ticker.C:
			sim.Tick()
			dirty = true
		case ev := <-eventCh:
			if ev == nil {
				return 0
			}
			switch tev := ev.(type) {
			case *tcell.EventResize:
				screen.Sync()
				fitSimulationToScreen(screen, sim)
				previous = nil
				needsFullClear = true
				dirty = true
			case *tcell.EventKey:
				if handleKeyEvent(state, sim, tev) {
					return 0
				}
				dirty = true
			}
		case <-sigCh:
			return 0
		}
	}
}

type cellCoord struct {
	x int
	y int
}

func diffCells(current engine.Board, previous *engine.Board, transient map[cellCoord]struct{}) ([]cellCoord, map[cellCoord]struct{}) {
	if previous == nil {
		return nil, nil
	}
	updates := make([]cellCoord, 0)
	nextTransient := make(map[cellCoord]struct{})
	for y := 0; y < current.Height(); y++ {
		for x := 0; x < current.Width(); x++ {
			coord := cellCoord{x: x, y: y}
			isAlive := current.IsAlive(x, y)
			wasAlive := previous.IsAlive(x, y)
			changed := isAlive != wasAlive
			if changed {
				nextTransient[coord] = struct{}{}
				updates = append(updates, coord)
				continue
			}
			if _, ok := transient[coord]; ok {
				updates = append(updates, coord)
			}
		}
	}
	return updates, nextTransient
}

func renderBoardFull(screen tcell.Screen, board engine.Board, previous *engine.Board, palette renderer.Palette) {
	for y := 0; y < board.Height(); y++ {
		for x := 0; x < board.Width(); x++ {
			isAlive := board.IsAlive(x, y)
			wasAlive := previous != nil && previous.IsAlive(x, y)
			r, style := cellRenderStyle(isAlive, wasAlive, palette)
			screen.SetContent(x, y, r, nil, style)
		}
	}
}

func renderCellUpdates(screen tcell.Screen, updates []cellCoord, current engine.Board, previous *engine.Board, palette renderer.Palette) {
	if previous == nil || len(updates) == 0 {
		return
	}
	for _, coord := range updates {
		isAlive := current.IsAlive(coord.x, coord.y)
		wasAlive := previous.IsAlive(coord.x, coord.y)
		r, style := cellRenderStyle(isAlive, wasAlive, palette)
		screen.SetContent(coord.x, coord.y, r, nil, style)
	}
}

func renderStatusBar(screen tcell.Screen, row int, status string) {
	if row < 0 {
		row = 0
	}
	width, height := screen.Size()
	if row >= height {
		row = height - 1
	}
	for x := 0; x < width; x++ {
		screen.SetContent(x, row, ' ', nil, tcell.StyleDefault)
	}
	for i, r := range status {
		if i >= width {
			break
		}
		screen.SetContent(i, row, r, nil, tcell.StyleDefault)
	}
}

func cellRenderStyle(isAlive, wasAlive bool, palette renderer.Palette) (rune, tcell.Style) {
	if isAlive {
		if !wasAlive {
			return '█', tcell.StyleDefault.Foreground(paletteColor(palette, palette.Newborn))
		}
		return '█', tcell.StyleDefault.Foreground(paletteColor(palette, palette.Alive))
	}
	if wasAlive {
		return ' ', tcell.StyleDefault.Foreground(paletteColor(palette, palette.RecentlyDead))
	}
	return ' ', tcell.StyleDefault.Foreground(paletteColor(palette, palette.Dead))
}

func paletteColor(palette renderer.Palette, value string) tcell.Color {
	switch palette.Mode {
	case renderer.ModeTrueColor:
		if len(value) == 7 && strings.HasPrefix(value, "#") {
			r, errR := parseHexByte(value[1:3])
			g, errG := parseHexByte(value[3:5])
			b, errB := parseHexByte(value[5:7])
			if errR == nil && errG == nil && errB == nil {
				return tcell.NewRGBColor(int32(r), int32(g), int32(b))
			}
		}
	case renderer.ModeFallback:
		if idx, err := parseDecimal(value); err == nil {
			return tcell.PaletteColor(idx)
		}
	}
	return tcell.ColorDefault
}

func parseHexByte(value string) (int64, error) {
	return strconv.ParseInt(value, 16, 32)
}

func parseDecimal(value string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if parsed < 0 {
		return 0, fmt.Errorf("invalid color")
	}
	return parsed, nil
}

func handleKeyEvent(state *input.State, sim *app.Simulation, ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyCtrlC {
		return true
	}

	key := mapKeyEvent(ev)
	if key == "" {
		return false
	}

	state.HandleKey(key)
	switch key {
	case "space":
		if state.Paused {
			sim.Pause()
		} else {
			sim.Resume()
		}
	case "r":
		sim.Restart()
	case "q":
		return true
	}
	return false
}

func mapKeyEvent(ev *tcell.EventKey) string {
	switch ev.Key() {
	case tcell.KeyRune:
		switch ev.Rune() {
		case ' ':
			return "space"
		case 'h', 'H':
			return "h"
		case '?':
			return "?"
		case 'r', 'R':
			return "r"
		case 'l', 'L':
			return "l"
		case 'q', 'Q':
			return "q"
		}
	}
	return ""
}

func mapKey(ch byte) string {
	switch ch {
	case ' ':
		return "space"
	case 'h', 'H':
		return "h"
	case '?':
		return "?"
	case 'r', 'R':
		return "r"
	case 'l', 'L':
		return "l"
	case 'q', 'Q':
		return "q"
	default:
		return ""
	}
}

func fitSimulationToScreen(screen tcell.Screen, sim *app.Simulation) {
	width, height := boardSizeForScreen(screen)
	sim.Resize(width, height)
}

func boardSizeForScreen(screen tcell.Screen) (int, int) {
	width, height := screen.Size()
	if height > 1 {
		height--
	}
	if width > frameMarginCols {
		width -= frameMarginCols
	}
	if height > frameMarginRows {
		height -= frameMarginRows
	}
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return width, height
}

func isTerminal(w io.Writer) bool {
	t, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := t.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func supportsTrueColor() bool {
	return strings.Contains(strings.ToLower(os.Getenv("COLORTERM")), "truecolor")
}

func patternSource(patternURL string) string {
	if patternURL == "" {
		return "random"
	}
	return patternURL
}

func tryLoadPatternForSimulation(sim *app.Simulation, patternURL string) error {
	if patternURL == "" {
		return nil
	}

	loader := pattern.NewHTTPWikiLoader(startupPatternTimeout, startupPatternMaxSize)
	content, err := loader.Load(patternURL)
	if err != nil {
		return err
	}
	return sim.LoadPatternFromWikiContent(content)
}
