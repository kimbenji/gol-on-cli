package input

import "testing"

func TestShouldTogglePlayPauseWhenSpaceIsPressed(t *testing.T) {
	state := NewState()

	state.HandleKey("space")
	if !state.Paused {
		t.Fatalf("expected paused state to be true after first space")
	}

	state.HandleKey("space")
	if state.Paused {
		t.Fatalf("expected paused state to be false after second space")
	}
}

func TestShouldToggleHelpWhenHOrQuestionMarkIsPressed(t *testing.T) {
	state := NewState()

	state.HandleKey("h")
	if !state.HelpVisible {
		t.Fatalf("expected help visible after h")
	}

	state.HandleKey("?")
	if state.HelpVisible {
		t.Fatalf("expected help hidden after ?")
	}
}

func TestShouldStartPatternLoadingFlowWhenLIsPressed(t *testing.T) {
	state := NewState()

	state.HandleKey("l")

	if !state.LoadPatternRequested {
		t.Fatalf("expected load pattern flow to start when l is pressed")
	}
}

func TestShouldSwitchToSafeExitStateWhenQIsPressed(t *testing.T) {
	state := NewState()

	state.HandleKey("q")

	if !state.ShouldQuit {
		t.Fatalf("expected safe exit state when q is pressed")
	}
}

func TestShouldResetLoadPatternFlagAfterConsume(t *testing.T) {
	state := NewState()
	state.HandleKey("l")

	if !state.ConsumeLoadPatternRequest() {
		t.Fatalf("expected first consume to observe pending load request")
	}
	if state.ConsumeLoadPatternRequest() {
		t.Fatalf("expected second consume to be false after reset")
	}
}
