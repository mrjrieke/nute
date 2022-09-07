package g3nfilter

import "github.com/mrjrieke/nute/mashupsdk"

type IG3nDisplayHintFilter interface {
	OnResize(displayHint *mashupsdk.MashupDisplayHint) *mashupsdk.MashupDisplayHint
}
