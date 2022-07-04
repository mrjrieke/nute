package g3nrender

import (
	"log"
	"sort"

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

func (mr *MashupRenderer) Sort(worldApp *g3nworld.WorldApp, g3nRenderableElements G3nCollection) G3nCollection {
	if len(g3nRenderableElements) == 0 {
		return g3nRenderableElements
	}
	if renderer, ok := mr.renderers[[]*g3nmash.G3nDetailedElement(g3nRenderableElements)[0].GetDetailedElement().GetRenderer()]; ok {
		return renderer.Sort(worldApp, g3nRenderableElements)
	} else {
		sort.Sort(g3nRenderableElements)
		return g3nRenderableElements
	}
}

func (mr *MashupRenderer) NextCoordinate(g3n *g3nmash.G3nDetailedElement) (*g3nmash.G3nDetailedElement, *math32.Vector3) {
	if renderer, ok := mr.renderers[g3n.GetDetailedElement().GetRenderer()]; ok {
		return renderer.NextCoordinate(g3n)
	} else {
		// Don't touch...
		return g3n, nil
	}
}

func (mr *MashupRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElementCollection []*g3nmash.G3nDetailedElement) {
	elementsByRenderer := map[string][]*g3nmash.G3nDetailedElement{}

	for _, g3nRenderableElement := range g3nRenderableElementCollection {
		concreteG3nRenderableElement := g3nRenderableElement
		if g3nRenderableElement.IsAbstract() {
			if tc, tErr := worldApp.GetG3nDetailedElementById(g3nRenderableElement.GetChildElements()[0]); tErr == nil {
				concreteG3nRenderableElement = tc
			} else {
				log.Printf("Skipping non-concrete abstract element: %d\n", g3nRenderableElement.GetBasisId())
				continue
			}
		}
		rendererName := concreteG3nRenderableElement.GetDetailedElement().GetRenderer()
		if _, ok := mr.renderers[rendererName]; ok {
			if _, renderOk := elementsByRenderer[rendererName]; !renderOk {
				elementsByRenderer[rendererName] = []*g3nmash.G3nDetailedElement{}
			}
			elementsByRenderer[rendererName] = append(elementsByRenderer[rendererName], concreteG3nRenderableElement)
		}
	}

	for _, g3nRenderableElements := range elementsByRenderer {
		g3nRenderableElements = mr.Sort(worldApp, G3nCollection(g3nRenderableElements))
		mr.GenericRenderer.LayoutBase(worldApp, mr, g3nRenderableElements)
	}

}
