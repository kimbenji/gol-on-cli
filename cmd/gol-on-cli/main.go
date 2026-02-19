package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"gol-on-cli/internal/app"
	"gol-on-cli/internal/cli"
	"gol-on-cli/internal/engine"
	"gol-on-cli/internal/input"
	"gol-on-cli/internal/renderer"
)

const version = "v0.1.0"

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
	fps := flags.Int("fps", 10, "updates per second")
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
	fileOut, outOK := stdout.(*os.File)
	if !inOK || !outOK {
		fmt.Fprintln(stderr, "failed to start: interactive mode requires file stdin/stdout")
		return 1
	}

	w, h := initialBoardSize(fileOut.Fd())
	sim := app.NewSimulation(w, h, *seed)
	return runFullscreen(sim, *fps, source, fileIn, fileOut)
}

func runFullscreen(sim *app.Simulation, fps int, source string, stdin *os.File, stdout io.Writer) int {
	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGWINCH)
	defer signal.Stop(sigCh)

	state := input.NewState()
	palette := renderer.SelectPalette(supportsTrueColor())

	orig, err := makeRaw(stdin.Fd())
	if err != nil {
		fmt.Fprintf(stdout, "failed to start: cannot enable raw input: %v\n", err)
		return 1
	}
	defer restoreTerm(stdin.Fd(), orig)

	keyCh := startKeyReader(stdin)
	fitSimulationToTerminal(sim, stdout)

	fmt.Fprint(stdout, "\x1b[?25l")
	defer fmt.Fprint(stdout, "\x1b[?25h\n")

	var previous *engine.Board
	for {
		current := sim.Board()
		frame := renderer.BuildFrameWithHistory(current, previous, renderer.StatusBarData{
			Generation:    sim.Generation(),
			Paused:        state.Paused,
			PatternSource: source,
		}, palette)
		fmt.Fprint(stdout, "\x1b[H")
		fmt.Fprint(stdout, frame)
		previousSnapshot := current
		previous = &previousSnapshot

		select {
		case <-ticker.C:
			sim.Tick()
		case key := <-keyCh:
			handleKey(state, sim, key)
			if state.ShouldQuit {
				return 0
			}
		case sig := <-sigCh:
			if sig == syscall.SIGWINCH {
				fitSimulationToTerminal(sim, stdout)
				previous = nil
				continue
			}
			return 0
		}
	}
}

func handleKey(state *input.State, sim *app.Simulation, key string) {
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
	}
}

func startKeyReader(stdin *os.File) <-chan string {
	keys := make(chan string, 16)
	go func() {
		defer close(keys)
		buf := make([]byte, 1)
		for {
			_, err := stdin.Read(buf)
			if err != nil {
				return
			}
			if key := mapKey(buf[0]); key != "" {
				keys <- key
			}
		}
	}()
	return keys
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

func fitSimulationToTerminal(sim *app.Simulation, stdout io.Writer) {
	file, ok := stdout.(*os.File)
	if !ok {
		return
	}
	width, height, err := terminalSize(file.Fd())
	if err != nil {
		return
	}
	if height > 1 {
		height--
	}
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	sim.Resize(width, height)
}

func initialBoardSize(fd uintptr) (int, int) {
	const startupMinWidth = 20
	const startupMinHeight = 10

	width, height, err := terminalSize(fd)
	if err != nil {
		return startupMinWidth, startupMinHeight
	}
	if height > 1 {
		height--
	}
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if width > startupMinWidth {
		width = startupMinWidth
	}
	if height > startupMinHeight {
		height = startupMinHeight
	}
	return width, height
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [19]uint8
	Ispeed uint32
	Ospeed uint32
}

func terminalSize(fd uintptr) (int, int, error) {
	ws := &winsize{}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
	if errno != 0 {
		return 0, 0, errno
	}
	return int(ws.Col), int(ws.Row), nil
}

func makeRaw(fd uintptr) (*termios, error) {
	state := &termios{}
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(state)), 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}

	raw := *state
	raw.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	raw.Oflag &^= syscall.OPOST
	raw.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	raw.Cflag &^= syscall.CSIZE | syscall.PARENB
	raw.Cflag |= syscall.CS8
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	_, _, errno = syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&raw)), 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}
	return state, nil
}

func restoreTerm(fd uintptr, state *termios) {
	if state == nil {
		return
	}
	syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(state)), 0, 0, 0)
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
