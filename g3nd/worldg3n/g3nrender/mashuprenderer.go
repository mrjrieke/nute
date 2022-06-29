package g3nrender

import (
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	"github.com/mrjrieke/nute/g3nd/g3nworld"
)

type MashupRenderer struct {
	GenericRenderer
	renderers map[string]G3nRenderer
}

func (mr *MashupRenderer) AddRenderer(genre string, renderer G3nRenderer) {
	if mr.renderers == nil {
		mr.renderers = map[string]G3nRenderer{}
	}
	mr.renderers[genre] = renderer
}

func (mr *MashupRenderer) NewSolidAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh {
	if g3n == nil {
		return nil
	}
	if renderer, ok := mr.renderers[g3n.GetDetailedElement().GetRenderer()]; ok {
		return renderer.NewSolidAtPosition(g3n, vpos)
	} else {
		return nil
	}
}

func (mr *MashupRenderer) NewInternalMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh {
	if g3n == nil {
		return nil
	}
	if renderer, ok := mr.renderers[g3n.GetDetailedElement().GetRenderer()]; ok {
		return renderer.NewInternalMeshAtPosition(g3n, vpos)
	} else {
		return nil
	}
}

func (mr *MashupRenderer) NextCoordinate(g3n *g3nmash.G3nDetailedElement, prevG3n *g3nmash.G3nDetailedElement, prevPos *math32.Vector3) (*g3nmash.G3nDetailedElement, *math32.Vector3) {
	if renderer, ok := mr.renderers[g3n.GetDetailedElement().GetRenderer()]; ok {
		return renderer.NextCoordinate(g3n, prevG3n, prevPos)
	} else {
		// Don't touch...
		return g3n, prevPos
	}
}

func (mr *MashupRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	mr.GenericRenderer.LayoutBase(worldApp, mr, g3nRenderableElements)
}
