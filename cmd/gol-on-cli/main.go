package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gol-on-cli/internal/app"
	"gol-on-cli/internal/cli"
	"gol-on-cli/internal/renderer"
)

const version = "v0.1.0"

type noopLoader struct{}

func (n noopLoader) Load(url string) error { return nil }

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
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

	sim := app.NewSimulation(60, 20, *seed)
	source := patternSource(*patternURL)

	if !isTerminal(stdout) {
		status := renderer.BuildStatusBar(renderer.StatusBarData{Generation: sim.Generation(), Paused: false, PatternSource: source})
		fmt.Fprintln(stdout, status)
		return 0
	}

	return runFullscreen(sim, *fps, source, stdout)
}

func runFullscreen(sim *app.Simulation, fps int, source string, stdout io.Writer) int {
	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	fmt.Fprint(stdout, "\x1b[?25l")
	defer fmt.Fprint(stdout, "\x1b[?25h\n")

	for {
		frame := renderer.BuildFrame(sim.Board(), renderer.StatusBarData{
			Generation:    sim.Generation(),
			Paused:        false,
			PatternSource: source,
		})
		fmt.Fprint(stdout, "\x1b[H\x1b[2J")
		fmt.Fprint(stdout, frame)

		select {
		case <-ticker.C:
			sim.Tick()
		case <-sigCh:
			return 0
		}
	}
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

func patternSource(patternURL string) string {
	if patternURL == "" {
		return "random"
	}
	return patternURL
}
