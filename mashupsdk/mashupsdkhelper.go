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
