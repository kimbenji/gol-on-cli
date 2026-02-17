package cli

import (
	"fmt"
	"strings"
)

type Loader interface {
	Load(url string) error
}

type StartOptions struct {
	PatternURL string
}

type StartResult struct {
	PatternLoadAttempted bool
}

func Start(options StartOptions, loader Loader) (StartResult, error) {
	if options.PatternURL == "" {
		return StartResult{}, nil
	}
	if err := loader.Load(options.PatternURL); err != nil {
		return StartResult{PatternLoadAttempted: true}, err
	}
	return StartResult{PatternLoadAttempted: true}, nil
}

func BuildHelpText() string {
	return strings.Join([]string{
		"Usage: gol-on-cli [options]",
		"",
		"Options:",
		"  --help          Show usage and options",
		"  --version       Show version",
		"  --fps <n>       Set updates per second",
		"  --seed <n>      Set random seed",
		"  --pattern-url   Load ConwayLife Wiki pattern on startup",
		"",
		"Shortcuts:",
		"  q, h/?, space, r, l",
		"",
		"URL Example:",
		"  https://conwaylife.com/wiki/Glider",
	}, "\n")
}

func BuildVersionText(version string) string {
	return fmt.Sprintf("%s", version)
}
