package g3nrender

import (
	"log"
	"sort"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	"github.com/mrjrieke/nute/g3nd/g3nworld"
)

type MashupRenderer struct {
	GenericRenderer
	renderers map[string]IG3nRenderer
}

func (mr *MashupRenderer) AddRenderer(genre string, renderer IG3nRenderer) {
	if mr.renderers == nil {
		mr.renderers = map[string]IG3nRenderer{}
	}
	mr.renderers[genre] = renderer
}

func (mr *MashupRenderer) GetRenderer(rendererName string) IG3nRenderer {
	return mr.renderers[rendererName]
}

func (mr *MashupRenderer) NewSolidAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) core.INode {
	if g3n == nil {
		return nil
	}
	if renderer, ok := mr.renderers[g3n.GetDetailedElement().GetRenderer()]; ok {
		return renderer.NewSolidAtPosition(g3n, vpos)
	} else {
		return nil
	}
}

func (mr *MashupRenderer) NewInternalMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) core.INode {
	if g3n == nil {
		return nil
	}
	if renderer, ok := mr.renderers[g3n.GetDetailedElement().GetRenderer()]; ok {
		return renderer.NewInternalMeshAtPosition(g3n, vpos)
	} else {
		return nil
	}
}

func (mr *MashupRenderer) NewRelatedMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3, vprevpos *math32.Vector3) core.INode {
	if g3n == nil {
		return nil
	}
	if renderer, ok := mr.renderers[g3n.GetDetailedElement().GetRenderer()]; ok {
		return renderer.NewRelatedMeshAtPosition(g3n, vpos, vprevpos)
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

func (mr *MashupRenderer) Collaborate(worldApp *g3nworld.WorldApp, collaboratingRenderer IG3nRenderer) {
	collaboratingRenderer.(IG3nRenderer).GetRenderer("").Collaborate(worldApp, collaboratingRenderer)
}

func (mr *MashupRenderer) NextCoordinate(g3n *g3nmash.G3nDetailedElement, totalElements int) (*g3nmash.G3nDetailedElement, *math32.Vector3) {
	if renderer, ok := mr.renderers[g3n.GetDetailedElement().GetRenderer()]; ok {
		return renderer.NextCoordinate(g3n, totalElements)
	} else {
		// Don't touch...
		return g3n, nil
	}
}
func (mr *MashupRenderer) CollectElementByRenderer(worldApp *g3nworld.WorldApp, elementsByRenderer map[string][]*g3nmash.G3nDetailedElement, concreteG3nRenderableElement *g3nmash.G3nDetailedElement) {
	rendererName := concreteG3nRenderableElement.GetDetailedElement().GetRenderer()

	if _, ok := mr.renderers[rendererName]; ok {
		if _, renderOk := elementsByRenderer[rendererName]; !renderOk {
			elementsByRenderer[rendererName] = []*g3nmash.G3nDetailedElement{}
		}
		elementsByRenderer[rendererName] = append(elementsByRenderer[rendererName], concreteG3nRenderableElement)
	}
}

func (mr *MashupRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElementCollection []*g3nmash.G3nDetailedElement) {
	elementsByRenderer := map[string][]*g3nmash.G3nDetailedElement{}

	for _, g3nRenderableElement := range g3nRenderableElementCollection {
		if g3nRenderableElement.IsAbstract() {
			for _, concreteComponentRenderable := range g3nRenderableElement.GetChildElementIds() {
				if tc, tErr := worldApp.GetG3nDetailedElementById(concreteComponentRenderable); tErr == nil {
					mr.CollectElementByRenderer(worldApp, elementsByRenderer, tc)
				} else {
					log.Printf("Skipping non-concrete abstract element: %d\n", g3nRenderableElement.GetBasisId())
					continue
				}
			}
			continue
		} else {
			mr.CollectElementByRenderer(worldApp, elementsByRenderer, g3nRenderableElement)
		}
	}

	for rendererName, g3nRenderableElements := range elementsByRenderer {
		g3nRenderableElements = mr.Sort(worldApp, G3nCollection(g3nRenderableElements))
		if renderer := mr.GetRenderer(rendererName); renderer != nil {
			if renderer.GetRendererType() == LAYOUT {
				renderer.Layout(worldApp, g3nRenderableElements)
			} else {
				mr.GenericRenderer.LayoutBase(worldApp, mr, g3nRenderableElements)
			}
		} else {
			mr.GenericRenderer.LayoutBase(worldApp, mr, g3nRenderableElements)
		}
	}
}

func (mr *MashupRenderer) InitRenderLoop(worldApp *g3nworld.WorldApp) bool {
	result := true
	for _, renderer := range mr.renderers {
		result = result && renderer.InitRenderLoop(worldApp)
	}

	return result
}

func (mr *MashupRenderer) RenderElement(worldApp *g3nworld.WorldApp, g3n *g3nmash.G3nDetailedElement) bool {
	if renderer, ok := mr.renderers[g3n.GetDetailedElement().GetRenderer()]; ok {
		return renderer.RenderElement(worldApp, g3n)
	} else {
		// Don't touch...
		return false
	}
}
