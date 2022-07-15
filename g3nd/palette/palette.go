package palette

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

var WHITE *math32.Color = &math32.Color{R: 1.0, G: 1.0, B: 1.0}
var DARK_RED *math32.Color = math32.NewColor("DarkRed")
var DARK_GREEN *math32.Color = math32.NewColor("DarkGreen")
var DARK_BLUE *math32.Color = math32.NewColor("DarkBlue")
var GREY *math32.Color = &math32.Color{R: 0.5, G: 0.5, B: 0.5}

var backgroundColorCache *math32.Color

func RefreshBackgroundColor(glState *gls.GLS, color *math32.Color, alpha float32) {
	if backgroundColorCache == nil || backgroundColorCache.R != color.R || backgroundColorCache.G != color.G || backgroundColorCache.B != color.B {
		glState.ClearColor(color.R, color.G, color.B, alpha)
		backgroundColorCache = color
	}
}
