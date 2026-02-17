package cli

import "testing"

func TestShouldPrintUsageOptionsShortcutsAndURLExampleForHelp(t *testing.T) {
	help := BuildHelpText()

	assertContains(t, help, "Usage:")
	assertContains(t, help, "--help")
	assertContains(t, help, "--version")
	assertContains(t, help, "--pattern-url")
	assertContains(t, help, "Shortcuts")
	assertContains(t, help, "q")
	assertContains(t, help, "h/?")
	assertContains(t, help, "space")
	assertContains(t, help, "https://conwaylife.com/wiki/Glider")
}

func TestShouldPrintVersionString(t *testing.T) {
	version := BuildVersionText("v1.2.3")
	if version != "v1.2.3" {
		t.Fatalf("expected exact version output, got %q", version)
	}
}

func TestShouldAttemptPatternLoadOnStartupWhenPatternURLIsProvided(t *testing.T) {
	loader := &spyLoader{}
	_, err := Start(StartOptions{PatternURL: "https://conwaylife.com/wiki/Glider", FPS: 10}, loader)
	if err != nil {
		t.Fatalf("expected startup to succeed, got error: %v", err)
	}
	if loader.calls != 1 {
		t.Fatalf("expected startup pattern load attempt once, got %d", loader.calls)
	}
	if loader.lastURL != "https://conwaylife.com/wiki/Glider" {
		t.Fatalf("expected startup pattern URL to be forwarded, got %q", loader.lastURL)
	}
}

type spyLoader struct {
	calls   int
	lastURL string
}

func (s *spyLoader) Load(url string) error {
	s.calls++
	s.lastURL = url
	return nil
}

func assertContains(t *testing.T, got, expected string) {
	t.Helper()
	if !contains(got, expected) {
		t.Fatalf("expected %q to contain %q", got, expected)
	}
}

func contains(source, needle string) bool {
	return len(needle) == 0 || (len(source) >= len(needle) && indexOf(source, needle) >= 0)
}

func indexOf(source, needle string) int {
	for i := 0; i+len(needle) <= len(source); i++ {
		if source[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

func TestShouldRejectNonPositiveFPS(t *testing.T) {
	_, err := Start(StartOptions{FPS: 0}, &spyLoader{})
	if err == nil {
		t.Fatalf("expected error for non-positive fps")
	}
}

func TestShouldRejectInvalidPatternURLBeforeLoading(t *testing.T) {
	loader := &spyLoader{}
	_, err := Start(StartOptions{PatternURL: "https://example.com/wiki/Glider", FPS: 10}, loader)
	if err == nil {
		t.Fatalf("expected invalid pattern URL to fail")
	}
	if loader.calls != 0 {
		t.Fatalf("expected loader not to be called for invalid URL")
	}
}
