package g3nrender

import (
	"fmt"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	"github.com/mrjrieke/nute/g3nd/g3nworld"
	g3ndpalette "github.com/mrjrieke/nute/g3nd/palette"
)

type TorusRenderer struct {
	GenericRenderer
	iOffset int
}

func (tr *TorusRenderer) NewSolidAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh {
	torusGeom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(g3ndpalette.DARK_BLUE)
	torusMesh := graphic.NewMesh(torusGeom, mat)
	fmt.Printf("LoaderID: %s\n", g3n.GetDisplayName())
	torusMesh.SetLoaderID(g3n.GetDisplayName())
	torusMesh.SetPositionVec(vpos)
	return torusMesh
}

func (tr *TorusRenderer) NewInternalMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3) *graphic.Mesh {
	diskGeom := geometry.NewDisk(1, 32)
	diskMat := material.NewStandard(g3ndpalette.GREY)
	diskMesh := graphic.NewMesh(diskGeom, diskMat)
	diskMesh.SetPositionVec(vpos)
	diskMesh.SetLoaderID(g3n.GetDisplayName())
	return diskMesh
}

func (tr *TorusRenderer) NewRelatedMeshAtPosition(g3n *g3nmash.G3nDetailedElement, vpos *math32.Vector3, vprevpos *math32.Vector3) *RelatedMesh {
	return nil
}

func (tr *TorusRenderer) NextCoordinate(g3n *g3nmash.G3nDetailedElement, totalElements int) (*g3nmash.G3nDetailedElement, *math32.Vector3) {
	switch tr.iOffset {
	case 0:
		tr.iOffset = tr.iOffset + 1
		return g3n, math32.NewVector3(float32(1.0), float32(1.0), float32(1.0))
	case 1:
		tr.iOffset = tr.iOffset + 1
		return g3n, math32.NewVector3(float32(-2.0), float32(-2.0), float32(-2.0))
	case 2:
		tr.iOffset = tr.iOffset + 1
		return g3n, math32.NewVector3(float32(2.0), float32(-2.0), float32(2.0))
	}
	return nil, nil
}

func (tr *TorusRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	tr.GenericRenderer.LayoutBase(worldApp, tr, g3nRenderableElements)
}

func (tr *TorusRenderer) HandleStateChange(worldApp *g3nworld.WorldApp, g3nDetailedElement *g3nmash.G3nDetailedElement) bool {
	itemColor := g3ndpalette.DARK_BLUE
	itemClickedColor := g3ndpalette.DARK_RED
	g3nColor := itemColor

	if g3nDetailedElement.IsItemActive() {
		g3nColor = itemClickedColor
	} else {
		if g3nDetailedElement.IsBackground() {
			if g3nDetailedElement.IsItemActive() {
				// No items clicked means background is clicked.
				g3nColor = itemClickedColor
			} else {
				g3nColor = g3ndpalette.GREY
			}
		} else {
			if g3nDetailedElement.IsBackgroundElement() {
				g3nColor = g3ndpalette.GREY
			}
		}
	}

	return g3nDetailedElement.SetColor(g3nColor)
}
