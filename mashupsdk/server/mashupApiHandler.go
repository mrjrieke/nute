package server

import (
	sdk "tini.com/nute/mashupsdk"
)

// MashupApiHandler -- mashups implement this to handle all events sent from
// other mashups.
type MashupApiHandler interface {
	OnResize(displayHint *sdk.MashupDisplayHint)
}
