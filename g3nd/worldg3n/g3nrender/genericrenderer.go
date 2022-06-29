package g3nrender

import (
	"log"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	"github.com/mrjrieke/nute/g3nd/g3nworld"
)

type G3nRenderer interface {
	NewSolidAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh
	NewInternalMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh
	NextCoordinate(g3n *g3nmash.G3nDetailedElement, prevG3n *g3nmash.G3nDetailedElement, prev *math32.Vector3) (*g3nmash.G3nDetailedElement, *math32.Vector3)
	Layout(worldApp *g3nworld.WorldApp, g3nRenderableElements []*g3nmash.G3nDetailedElement)
}

type GenericRenderer struct {
}

func (*GenericRenderer) NewSolidAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh {
	return nil
}

func (*GenericRenderer) NewInternalMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh {
	return nil
}

func (*GenericRenderer) NextCoordinate(g3n *g3nmash.G3nDetailedElement, prevG3n *g3nmash.G3nDetailedElement, prev *math32.Vector3) (*g3nmash.G3nDetailedElement, *math32.Vector3) {
	return g3n, math32.NewVector3(float32(0.0), float32(0.0), float32(0.0))
}

func (gr *GenericRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	gr.LayoutBase(worldApp, gr, g3nRenderableElements)
}

func (gr *GenericRenderer) LayoutBase(worldApp *g3nworld.WorldApp,
	g3Renderer G3nRenderer,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	var prevPos *math32.Vector3
	var prevG3n *g3nmash.G3nDetailedElement

	for _, g3nRenderableElement := range g3nRenderableElements {
		concreteG3nRenderableElement := g3nRenderableElement
		if g3nRenderableElement.IsAbstract() {
			if tc, tErr := worldApp.GetG3nDetailedElementById(g3nRenderableElement.GetChildElements()[0]); tErr == nil {
				concreteG3nRenderableElement = tc
			} else {
				log.Printf("Skipping non-concrete abstract element: %d\n", g3nRenderableElement.GetBasisId())
				continue
			}
		}

		prevG3n, prevPos = g3Renderer.NextCoordinate(concreteG3nRenderableElement, prevG3n, prevPos)
		solidMesh := g3Renderer.NewSolidAtPosition(concreteG3nRenderableElement, prevPos)
		if solidMesh != nil {
			worldApp.AddToScene(solidMesh)
			concreteG3nRenderableElement.SetNamedMesh(concreteG3nRenderableElement.GetDisplayName(), solidMesh)
		}

		for _, innerG3n := range worldApp.GetG3nDetailedChildElementsByGenre(concreteG3nRenderableElement, "Space") {
			negativeMesh := g3Renderer.NewInternalMeshAtPosition(innerG3n, prevPos)
			if negativeMesh != nil {
				worldApp.AddToScene(negativeMesh)
				innerG3n.SetNamedMesh(innerG3n.GetDisplayName(), negativeMesh)
			}
		}
	}
}
