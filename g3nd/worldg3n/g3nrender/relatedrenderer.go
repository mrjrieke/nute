package g3nrender

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	"github.com/mrjrieke/nute/g3nd/g3nworld"
	g3ndpalette "github.com/mrjrieke/nute/g3nd/palette"
)

type RelatedRenderer struct {
	GenericRenderer
	iOffset int
}

type RelatedMesh struct {
	graphic.Mesh
	PrevPos *math32.Vector3 // Position of related mesh
}

func (tr *RelatedRenderer) NewSolidAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh {
	return nil
}

func (tr *RelatedRenderer) NewInternalMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh {
	return nil
}

func (tr *RelatedRenderer) NewRelatedMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3, vprevpos *math32.Vector3) graphic.Mesh {
	spherGeom := geometry.NewSphere(1.0, 5, 5)
	sphereMat := material.NewStandard(g3ndpalette.GREY)
	sphereMesh := graphic.NewMesh(spherGeom, sphereMat)
	sphereMesh.SetPositionVec(vpos)
	sphereMesh.SetLoaderID(g3n.GetDisplayName())
	relatedMesh := RelatedMesh{Mesh: *sphereMesh, PrevPos: vprevpos}

	return relatedMesh
}

func (tr *RelatedRenderer) NextCoordinate(g3n *g3nmash.G3nDetailedElement) (*g3nmash.G3nDetailedElement, *math32.Vector3) {
	if tr.iOffset == 0 {
		tr.iOffset = 1
		return g3n, math32.NewVector3(float32(-2.0), float32(-2.0), float32(-2.0))
	} else {
		return g3n, math32.NewVector3(float32(2.0), float32(2.0), float32(2.0))
	}
}

func (tr *RelatedRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	tr.GenericRenderer.LayoutBase(worldApp, tr, g3nRenderableElements)
}
