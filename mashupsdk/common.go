package mashupsdk

const (
	Init    int64 = 0
	Rest    int64 = 1
	Clicked int64 = 2
	Moved   int64 = 3
)

type MashupDisplayState int

const (
	Configured MashupDisplayState = 1
	Position   MashupDisplayState = 2
	Frame      MashupDisplayState = 4
	AppInitted MashupDisplayState = 8
)

type MashupDisplayContext struct {
	MainWinDisplay *MashupDisplayHint
	yOffset        int
	settled        MashupDisplayState
}

func (m *MashupDisplayContext) GetYoffset() int {
	return m.yOffset
}

func (m *MashupDisplayContext) SetYoffset(a int) int {
	m.yOffset = a
	return m.yOffset
}

func (m *MashupDisplayContext) GetSettled() MashupDisplayState {
	return m.settled
}

func (m *MashupDisplayContext) ApplySettled(s MashupDisplayState, override bool) MashupDisplayState {
	if override {
		m.settled = s
	} else {
		m.settled |= s
	}
	return m.settled
}

func (m *MashupDisplayContext) OnResize(displayHint *MashupDisplayHint) bool {
	resize := false
	if m.MainWinDisplay == nil {
		resize = true
		m.MainWinDisplay = &MashupDisplayHint{}
	}

	if displayHint == nil {
		return false
	}

	if displayHint.Xpos != 0 && (*m.MainWinDisplay).Xpos != displayHint.Xpos {
		m.MainWinDisplay.Xpos = displayHint.Xpos
		resize = true
	}
	if displayHint.Ypos != 0 && (*m.MainWinDisplay).Ypos != displayHint.Ypos {
		m.MainWinDisplay.Ypos = displayHint.Ypos + int64(m.yOffset)
		resize = true
	}
	if displayHint.Width != 0 && (*m.MainWinDisplay).Width != displayHint.Width {
		m.MainWinDisplay.Width = displayHint.Width
		resize = true
	}
	if displayHint.Height != 0 && (*m.MainWinDisplay).Height != displayHint.Height+int64(m.yOffset) {
		m.MainWinDisplay.Height = displayHint.Height
		resize = true
	}

	if m.settled < 15 {
		return false
	} else if m.settled == 15 {
		resize = true
		m.settled = 31
	}
	return resize
}

// MashupApiHandler -- mashups implement this to handle all events sent from
// other mashups.
type MashupApiHandler interface {
	OnResize(displayHint *MashupDisplayHint)
	UpsertMashupElements(detailedElementBundle *MashupDetailedElementBundle) (*MashupElementStateBundle, error)
	UpsertMashupElementsState(elementStateBundle *MashupElementStateBundle) (*MashupElementStateBundle, error)
}

type MashupContextInitHandler interface {
	RegisterContext(*MashupContext)
}
