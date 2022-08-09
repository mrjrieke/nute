package mashupsdk

func (m *MashupDetailedElement) Copy(source *MashupDetailedElement) {
	m.Basisid = source.Basisid
	m.Id = source.Id
	if m.State == nil {
		m.State = &MashupElementState{}
	}
	m.State.Id = source.State.Id
	m.State.State = source.State.State
	m.Name = source.Name
	m.Alias = source.Alias
	m.Description = source.Description
	m.Renderer = source.Renderer
	m.Colabrenderer = source.Colabrenderer
	m.Genre = source.Genre
	m.Subgenre = source.Subgenre
	m.Parentids = source.Parentids
	m.Childids = source.Childids
}

func (m *MashupDetailedElement) IsStateSet(stateBit DisplayElementState) bool {
	displayState := m.State.State
	return (displayState & int64(stateBit)) == int64(stateBit)
}

func (m *MashupDetailedElement) ApplyState(x DisplayElementState, isset bool) bool {
	if isset {
		m.State.State |= int64(x)
	} else {
		if m.IsStateSet(x) {
			m.State.State &= ^int64(x)
		}
	}

	return false
}
