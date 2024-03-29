package g3nrender

import (
	"log"
	"sort"
	"strings"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	"github.com/mrjrieke/nute/g3nd/g3nworld"
)

type G3nCollection []*g3nmash.G3nDetailedElement

func (a G3nCollection) Len() int { return len(a) }
func (a G3nCollection) Less(i, j int) bool {
	return strings.Compare(a[i].GetDisplayName(), a[j].GetDisplayName()) < 0
}
func (a G3nCollection) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type IG3nRenderer interface {
	NewSolidAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) core.INode
	NewInternalMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) core.INode
	NewRelatedMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3, vprevpos *math32.Vector3) core.INode
	NextCoordinate(g3n *g3nmash.G3nDetailedElement, totalElements int) (*g3nmash.G3nDetailedElement, *math32.Vector3)
	Sort(worldApp *g3nworld.WorldApp, g3nRenderableElements G3nCollection) G3nCollection
	Layout(worldApp *g3nworld.WorldApp, g3nRenderableElements []*g3nmash.G3nDetailedElement)
	InitRenderLoop(worldApp *g3nworld.WorldApp) bool
	RenderElement(worldApp *g3nworld.WorldApp, g3n *g3nmash.G3nDetailedElement) bool
	GetRenderer(rendererName string) IG3nRenderer
	GetRendererType() G3nRenderType
	Collaborate(worldApp *g3nworld.WorldApp, renderer IG3nRenderer)
}

type G3nRenderType string

const (
	NONE   G3nRenderType = ""
	LAYOUT G3nRenderType = "Layout"
)

type GenericRenderer struct {
	RendererType G3nRenderType
}

func (*GenericRenderer) NewSolidAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) core.INode {
	return nil
}

func (*GenericRenderer) NewInternalMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) core.INode {
	return nil
}

func (*GenericRenderer) NewRelatedMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3, vprevpos *math32.Vector3) core.INode {
	return nil
}

func (*GenericRenderer) NextCoordinate(g3n *g3nmash.G3nDetailedElement, totalElements int) (*g3nmash.G3nDetailedElement, *math32.Vector3) {
	return g3n, math32.NewVector3(float32(0.0), float32(0.0), float32(0.0))
}

func (gr *GenericRenderer) Sort(worldApp *g3nworld.WorldApp, g3nRenderableElements G3nCollection) G3nCollection {
	sort.Sort(g3nRenderableElements)
	return g3nRenderableElements
}

func (gr *GenericRenderer) InitRenderLoop(worldApp *g3nworld.WorldApp) bool {
	// TODO: noop
	return true
}

func (gr *GenericRenderer) RenderElement(worldApp *g3nworld.WorldApp, g3n *g3nmash.G3nDetailedElement) bool {
	return false
}

func (gr *GenericRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	gr.LayoutBase(worldApp, gr, g3nRenderableElements)
}

func (gr *GenericRenderer) GetRenderer(rendererName string) IG3nRenderer {
	return nil
}

func (gr *GenericRenderer) GetRendererType() G3nRenderType {
	return gr.RendererType
}

func (gr *GenericRenderer) Collaborate(worldApp *g3nworld.WorldApp, collaboratingRenderer IG3nRenderer) {

}

func (gr *GenericRenderer) LayoutBase(worldApp *g3nworld.WorldApp,
	g3Renderer IG3nRenderer,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	var nextPos *math32.Vector3
	var prevSolidPos *math32.Vector3

	totalElements := len(g3nRenderableElements)

	if totalElements > 0 {
		if g3nRenderableElements[0].GetDetailedElement().Colabrenderer != "" {
			log.Printf("Collab examine: %v\n", g3nRenderableElements[0])
			log.Printf("Renderer name: %s\n", g3nRenderableElements[0].GetDetailedElement().GetRenderer())
			protoRenderer := g3Renderer.GetRenderer(g3nRenderableElements[0].GetDetailedElement().GetRenderer())
			log.Printf("Collaborating %v\n", protoRenderer)
			g3Renderer.Collaborate(worldApp, protoRenderer)
		}
	}

	for _, g3nRenderableElement := range g3nRenderableElements {
		concreteG3nRenderableElement := g3nRenderableElement

		prevSolidPos = nextPos
		_, nextPos = g3Renderer.NextCoordinate(concreteG3nRenderableElement, totalElements)
		solidMesh := g3Renderer.NewSolidAtPosition(concreteG3nRenderableElement, nextPos)
		if solidMesh != nil {
			log.Printf("Adding %s\n", solidMesh.GetNode().LoaderID())
			worldApp.AddToScene(solidMesh)
			concreteG3nRenderableElement.SetNamedMesh(concreteG3nRenderableElement.GetDisplayName(), solidMesh)
		}

		for _, relatedG3n := range worldApp.GetG3nDetailedChildElementsByGenre(concreteG3nRenderableElement, "Related") {
			relatedMesh := g3Renderer.NewRelatedMeshAtPosition(concreteG3nRenderableElement, nextPos, prevSolidPos)
			if relatedMesh != nil {
				worldApp.AddToScene(relatedMesh)
				concreteG3nRenderableElement.SetNamedMesh(relatedG3n.GetDisplayName(), relatedMesh)
			}
		}

		for _, innerG3n := range worldApp.GetG3nDetailedChildElementsByGenre(concreteG3nRenderableElement, "Space") {
			negativeMesh := g3Renderer.NewInternalMeshAtPosition(innerG3n, nextPos)
			if negativeMesh != nil {
				worldApp.AddToScene(negativeMesh)
				innerG3n.SetNamedMesh(innerG3n.GetDisplayName(), negativeMesh)
			}
		}
	}
}
