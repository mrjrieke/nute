package g3nrender

import (
	"log"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	"github.com/mrjrieke/nute/g3nd/g3nworld"
	g3ndpalette "github.com/mrjrieke/nute/g3nd/palette"
	"github.com/mrjrieke/nute/mashupsdk"
)

type BackgroundRenderer struct {
	GenericRenderer
	CollaboratingRenderer IG3nRenderer
	ActiveColor           *math32.Color
}

func (br *BackgroundRenderer) NewSolidAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) core.INode {
	return nil
}

func (br *BackgroundRenderer) NewInternalMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) core.INode {
	return nil
}

func (br *BackgroundRenderer) NewRelatedMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3, vprevpos *math32.Vector3) core.INode {
	return nil
}

func (br *BackgroundRenderer) NextCoordinate(g3n *g3nmash.G3nDetailedElement, totalElements int) (*g3nmash.G3nDetailedElement, *math32.Vector3) {
	return nil, nil
}

func (br *BackgroundRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	return
}

func (br *BackgroundRenderer) GetRenderer(rendererName string) IG3nRenderer {
	if br.CollaboratingRenderer != nil {
		return br.CollaboratingRenderer
	}
	return nil
}

func (br *BackgroundRenderer) InitRenderLoop(worldApp *g3nworld.WorldApp) bool {
	// TODO: noop
	return true
}

func (br *BackgroundRenderer) RenderElement(worldApp *g3nworld.WorldApp, g3nDetailedElement *g3nmash.G3nDetailedElement) bool {
	var g3nColor *math32.Color

	if br.ActiveColor == nil {
		br.ActiveColor = g3ndpalette.DARK_RED
	}

	if g3nDetailedElement.IsStateSet(mashupsdk.Clicked) {
		g3nColor = g3ndpalette.DARK_RED
		if br.ActiveColor == g3ndpalette.DARK_RED {
			br.ActiveColor = g3ndpalette.DARK_GREEN
		} else {
			br.ActiveColor = g3ndpalette.DARK_RED
		}
		log.Printf("Active color: %v\n", br.ActiveColor)
	} else {
		g3nColor = g3ndpalette.GREY
	}

	return g3nDetailedElement.SetColor(g3nColor, 1.0)
}

func (br *BackgroundRenderer) Collaborate(worldApp *g3nworld.WorldApp, collaboratingRenderer IG3nRenderer) {
	br.CollaboratingRenderer.Collaborate(worldApp, br)
}
