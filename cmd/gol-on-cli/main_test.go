package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestShouldPrintStatusBarOnDefaultRun(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run(nil, strings.NewReader(""), &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("expected success exit code, got %d with stderr %q", exitCode, stderr.String())
	}
	if stdout.Len() == 0 {
		t.Fatalf("expected default run output, got empty stdout")
	}
	if !strings.Contains(stdout.String(), "gen:0") {
		t.Fatalf("expected status bar generation in output, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "source:random") {
		t.Fatalf("expected random source in output, got %q", stdout.String())
	}
}

func TestShouldPrintPatternSourceWhenPatternURLIsProvided(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{"--pattern-url", "https://conwaylife.com/wiki/Glider"}, strings.NewReader(""), &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("expected success exit code, got %d with stderr %q", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "source:https://conwaylife.com/wiki/Glider") {
		t.Fatalf("expected URL source in output, got %q", stdout.String())
	}
}

func TestShouldMapShortcutKeysFromRawBytes(t *testing.T) {
	cases := map[byte]string{
		' ': "space",
		'h': "h",
		'?': "?",
		'r': "r",
		'l': "l",
		'q': "q",
	}

	for in, expected := range cases {
		if got := mapKey(in); got != expected {
			t.Fatalf("expected key %q to map to %q, got %q", in, expected, got)
		}
	}
}
