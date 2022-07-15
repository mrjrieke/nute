package g3nrender

import (
	"log"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	"github.com/mrjrieke/nute/g3nd/g3nworld"
	g3ndpalette "github.com/mrjrieke/nute/g3nd/palette"
)

type BackgroundRenderer struct {
	GenericRenderer
	CollaboratingRenderer G3nRenderer
	ActiveColor           *math32.Color
}

func (br *BackgroundRenderer) NewSolidAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh {
	return nil
}

func (br *BackgroundRenderer) NewInternalMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh {
	return nil
}

func (br *BackgroundRenderer) NewRelatedMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3, vprevpos *math32.Vector3) *RelatedMesh {
	return nil
}

func (br *BackgroundRenderer) NextCoordinate(g3n *g3nmash.G3nDetailedElement, totalElements int) (*g3nmash.G3nDetailedElement, *math32.Vector3) {
	return nil, nil
}

func (br *BackgroundRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	return
}

func (br *BackgroundRenderer) GetRenderer(rendererName string) G3nRenderer {
	if br.CollaboratingRenderer != nil {
		return br.CollaboratingRenderer
	}
	return nil
}

func (br *BackgroundRenderer) HandleStateChange(worldApp *g3nworld.WorldApp, g3nDetailedElement *g3nmash.G3nDetailedElement) bool {
	var g3nColor *math32.Color

	if br.ActiveColor == nil {
		br.ActiveColor = g3ndpalette.DARK_RED
	}

	if g3nDetailedElement.IsItemActive() {
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

	return g3nDetailedElement.SetColor(g3nColor)
}

func (br *BackgroundRenderer) Collaborate(worldApp *g3nworld.WorldApp, collaboratingRenderer interface{}) {
	br.CollaboratingRenderer.Collaborate(worldApp, br)
}
