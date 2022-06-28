package g3nrender

import (
	"log"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	"github.com/mrjrieke/nute/g3nd/g3nworld"
)

type G3nRenderer interface {
	NewSolidAtPosition(displayName string, vpos *math32.Vector3) *graphic.Mesh
	NewInternalMeshAtPosition(displayName string, vpos *math32.Vector3) *graphic.Mesh
	NextCoordinate() *math32.Vector3
	Layout(worldApp *g3nworld.WorldApp, g3nRenderableElements []*g3nmash.G3nDetailedElement)
}

type GenericRenderer struct {
}

func (*GenericRenderer) NewSolidAtPosition(displayName string, vpos *math32.Vector3) *graphic.Mesh {
	return nil
}

func (*GenericRenderer) NewInternalMeshAtPosition(displayName string, vpos *math32.Vector3) *graphic.Mesh {
	return nil
}

func (*GenericRenderer) NextCoordinate() *math32.Vector3 {
	return math32.NewVector3(float32(0.0), float32(0.0), float32(0.0))
}

func (gr *GenericRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	gr.LayoutBase(worldApp, gr, g3nRenderableElements)
}

func (gr *GenericRenderer) LayoutBase(worldApp *g3nworld.WorldApp,
	g3Renderer G3nRenderer,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {

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

		nextPos := g3Renderer.NextCoordinate()
		solidMesh := g3Renderer.NewSolidAtPosition(concreteG3nRenderableElement.GetDisplayName(), nextPos)
		worldApp.AddToScene(solidMesh)
		concreteG3nRenderableElement.SetNamedMesh(concreteG3nRenderableElement.GetDisplayName(), solidMesh)

		for _, innerG3n := range worldApp.GetG3nDetailedChildElementsByGenre(concreteG3nRenderableElement, "Space") {
			negativeMesh := g3Renderer.NewInternalMeshAtPosition(innerG3n.GetDisplayName(), nextPos)
			worldApp.AddToScene(negativeMesh)
			innerG3n.SetNamedMesh(innerG3n.GetDisplayName(), negativeMesh)
		}
	}
}
