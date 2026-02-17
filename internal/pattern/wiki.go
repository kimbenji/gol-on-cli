package pattern

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"gol-on-cli/internal/engine"
)

type PatternFormat string

const (
	FormatRLE       PatternFormat = "RLE"
	FormatPlainText PatternFormat = "PlainText"
	FormatLife106   PatternFormat = "Life1.06"
)

type RecoverableError struct {
	Message string
}

func (e RecoverableError) Error() string {
	return e.Message
}

func ValidateWikiURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return parsed.Scheme == "https" && parsed.Host == "conwaylife.com" && strings.HasPrefix(parsed.Path, "/wiki/")
}

func SelectPreferredPattern(content string) (PatternFormat, string, error) {
	if body, ok := extractRLE(content); ok {
		return FormatRLE, body, nil
	}
	if body, ok := extractPlainText(content); ok {
		return FormatPlainText, body, nil
	}
	if body, ok := extractLife106(content); ok {
		return FormatLife106, body, nil
	}

	return "", "", RecoverableError{Message: "no supported pattern format found"}
}

func LoadBoardFromWikiContent(content string, width, height int) (engine.Board, error) {
	format, body, err := SelectPreferredPattern(content)
	if err != nil {
		return engine.Board{}, err
	}

	board, err := ParseToBoard(format, body, width, height)
	if err != nil {
		return engine.Board{}, RecoverableError{Message: err.Error()}
	}

	return board, nil
}

func ParseToBoard(format PatternFormat, body string, width, height int) (engine.Board, error) {
	switch format {
	case FormatRLE:
		return parseRLE(body, width, height)
	case FormatPlainText:
		return parsePlainText(body, width, height)
	case FormatLife106:
		return parseLife106(body, width, height)
	default:
		return engine.Board{}, fmt.Errorf("unsupported format: %s", format)
	}
}

func extractRLE(content string) (string, bool) {
	lines := strings.Split(content, "\n")
	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "x") && strings.Contains(trimmed, "=") && strings.Contains(trimmed, "y") {
			collected := []string{trimmed}
			for i := index + 1; i < len(lines); i++ {
				next := strings.TrimSpace(lines[i])
				if next == "" {
					break
				}
				collected = append(collected, next)
				if strings.Contains(next, "!") {
					break
				}
			}
			return strings.Join(collected, "\n"), true
		}
	}
	return "", false
}

func extractPlainText(content string) (string, bool) {
	lines := strings.Split(content, "\n")
	var collected []string
	collecting := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "!") {
			if collecting {
				break
			}
			continue
		}

		if isPlainTextRow(trimmed) {
			collecting = true
			collected = append(collected, trimmed)
			continue
		}

		if collecting {
			break
		}
	}

	if len(collected) == 0 {
		return "", false
	}
	return strings.Join(collected, "\n"), true
}

func isPlainTextRow(line string) bool {
	if line == "" {
		return false
	}
	for _, char := range line {
		if char != '.' && char != 'O' {
			return false
		}
	}
	return true
}

func extractLife106(content string) (string, bool) {
	lines := strings.Split(content, "\n")
	for index, line := range lines {
		if strings.TrimSpace(line) == "#Life 1.06" {
			collected := []string{"#Life 1.06"}
			for i := index + 1; i < len(lines); i++ {
				next := strings.TrimSpace(lines[i])
				if next == "" {
					break
				}
				collected = append(collected, next)
			}
			return strings.Join(collected, "\n"), true
		}
	}
	return "", false
}

func parseRLE(body string, width, height int) (engine.Board, error) {
	lines := strings.Split(body, "\n")
	if len(lines) < 2 {
		return engine.Board{}, fmt.Errorf("invalid RLE: missing body")
	}

	board := engine.NewBoard(width, height)
	x := 0
	y := 0
	runLength := 0
	for _, char := range strings.Join(lines[1:], "") {
		switch {
		case char >= '0' && char <= '9':
			runLength = runLength*10 + int(char-'0')
		case char == 'b' || char == 'o':
			count := max(1, runLength)
			for i := 0; i < count; i++ {
				if char == 'o' && x < width && y < height {
					board.SetAlive(x, y, true)
				}
				x++
			}
			runLength = 0
		case char == '$':
			count := max(1, runLength)
			y += count
			x = 0
			runLength = 0
		case char == '!':
			return board, nil
		default:
			return engine.Board{}, fmt.Errorf("invalid RLE token: %q", char)
		}
	}
	return engine.Board{}, fmt.Errorf("invalid RLE: missing terminator")
}

func parsePlainText(body string, width, height int) (engine.Board, error) {
	board := engine.NewBoard(width, height)
	lines := strings.Split(body, "\n")
	for y, line := range lines {
		if !isPlainTextRow(line) {
			return engine.Board{}, fmt.Errorf("invalid PlainText row: %q", line)
		}
		for x, char := range line {
			if char == 'O' && x < width && y < height {
				board.SetAlive(x, y, true)
			}
		}
	}
	return board, nil
}

func parseLife106(body string, width, height int) (engine.Board, error) {
	lines := strings.Split(body, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "#Life 1.06" {
		return engine.Board{}, fmt.Errorf("invalid Life 1.06 header")
	}

	board := engine.NewBoard(width, height)
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return engine.Board{}, fmt.Errorf("invalid Life 1.06 coordinate: %q", line)
		}
		x, err := strconv.Atoi(fields[0])
		if err != nil {
			return engine.Board{}, fmt.Errorf("invalid Life 1.06 x coordinate: %q", fields[0])
		}
		y, err := strconv.Atoi(fields[1])
		if err != nil {
			return engine.Board{}, fmt.Errorf("invalid Life 1.06 y coordinate: %q", fields[1])
		}
		if x >= 0 && y >= 0 && x < width && y < height {
			board.SetAlive(x, y, true)
		}
	}
	return board, nil
}

func max(left, right int) int {
	if left > right {
		return left
	}
	return right
}
