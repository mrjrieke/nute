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

type TorusRenderer struct {
	GenericRenderer
}

func (*TorusRenderer) NewSolidAtPosition(displayName string, vpos *math32.Vector3) *graphic.Mesh {
	torusGeom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(g3ndpalette.DARK_BLUE)
	torusMesh := graphic.NewMesh(torusGeom, mat)
	torusMesh.SetLoaderID(displayName)
	torusMesh.SetPositionVec(vpos)
	return torusMesh
}

func (*TorusRenderer) NewInternalMeshAtPosition(displayName string, vpos *math32.Vector3) *graphic.Mesh {
	diskGeom := geometry.NewDisk(1, 32)
	diskMat := material.NewStandard(g3ndpalette.GREY)
	diskMesh := graphic.NewMesh(diskGeom, diskMat)
	diskMesh.SetPositionVec(vpos)
	diskMesh.SetLoaderID(displayName)
	return diskMesh
}

func (*TorusRenderer) NextCoordinate() *math32.Vector3 {
	return math32.NewVector3(float32(0.0), float32(0.0), float32(0.0))
}

func (tr *TorusRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	tr.GenericRenderer.LayoutBase(worldApp, tr, g3nRenderableElements)
}
