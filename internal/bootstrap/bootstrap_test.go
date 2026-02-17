package bootstrap

import (
	"testing"

	"gol-on-cli/internal/app"
	"gol-on-cli/internal/input"
)

func TestShouldReportBootstrapReady(t *testing.T) {
	if !IsReady() {
		t.Fatalf("expected bootstrap to be ready")
	}
}

func TestShouldStartSimulationWithRandomPatternOnDefaultRun(t *testing.T) {
	sim := app.NewSimulation(6, 4, 42)

	if sim.Generation() != 0 {
		t.Fatalf("expected default run to start at generation 0, got %d", sim.Generation())
	}
	if sim.Board().Width() != 6 || sim.Board().Height() != 4 {
		t.Fatalf("expected default run board to be initialized, got %dx%d", sim.Board().Width(), sim.Board().Height())
	}
}

func TestShouldRunCoreUserScenarioWithoutFailure(t *testing.T) {
	sim := app.NewSimulation(5, 5, 11)
	state := input.NewState()

	state.HandleKey("space")
	if !state.Paused {
		t.Fatalf("expected pause to be enabled after first space")
	}

	sim.Restart()
	if sim.Generation() != 0 {
		t.Fatalf("expected generation reset after restart, got %d", sim.Generation())
	}

	err := sim.LoadPatternFromWikiContent("x = 3, y = 3\nbo$2bo$3o!")
	if err != nil {
		t.Fatalf("expected external pattern load to succeed, got %v", err)
	}
}

func TestShouldSatisfyAllDoDChecklistItems(t *testing.T) {
	if !AllDoDItemsSatisfied() {
		t.Fatalf("expected DoD checklist to be fully satisfied")
	}
}
