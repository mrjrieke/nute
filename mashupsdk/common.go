package mashupsdk

const (
	Init    int64 = 0
	Rest    int64 = 1
	Clicked int64 = 2
	Moved   int64 = 3
)

// MashupApiHandler -- mashups implement this to handle all events sent from
// other mashups.
type MashupApiHandler interface {
	OnResize(displayHint *MashupDisplayHint)
	UpsertMashupSociety(societyBundle *MashupSocietyBundle) (*MashupSocietyStateBundle, error)
	UpsertMashupSocietyState(societyStateBundle *MashupSocietyStateBundle) (*MashupSocietyStateBundle, error)
}
