package input

type State struct {
	Paused               bool
	HelpVisible          bool
	LoadPatternRequested bool
	ShouldQuit           bool
}

func NewState() *State {
	return &State{}
}

func (s *State) HandleKey(key string) {
	switch key {
	case "space":
		s.Paused = !s.Paused
	case "h", "?":
		s.HelpVisible = !s.HelpVisible
	case "l":
		s.LoadPatternRequested = true
	case "q":
		s.ShouldQuit = true
	}
}
