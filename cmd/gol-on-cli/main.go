package main

import (
	"flag"
	"fmt"
	"os"

	"gol-on-cli/internal/cli"
)

const version = "v0.1.0"

type noopLoader struct{}

func (n noopLoader) Load(url string) error { return nil }

func main() {
	help := flag.Bool("help", false, "show usage")
	showVersion := flag.Bool("version", false, "show version")
	fps := flag.Int("fps", 10, "updates per second")
	seed := flag.Int64("seed", 0, "random seed")
	patternURL := flag.String("pattern-url", "", "startup pattern URL")
	aliveColor := flag.String("alive-color", "", "alive cell color")
	deadColor := flag.String("dead-color", "", "dead cell color")
	flag.Parse()

	_ = fps
	_ = seed
	_ = aliveColor
	_ = deadColor

	if *help {
		fmt.Println(cli.BuildHelpText())
		return
	}
	if *showVersion {
		fmt.Println(cli.BuildVersionText(version))
		return
	}

	if _, err := cli.Start(cli.StartOptions{PatternURL: *patternURL, FPS: *fps}, noopLoader{}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start: %v\n", err)
		os.Exit(1)
	}
}
