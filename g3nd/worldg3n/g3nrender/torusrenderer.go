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

func (tr *TorusRenderer) NewSolidAtPosition(displayName string, vpos *math32.Vector3) *graphic.Mesh {
	torusGeom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(g3ndpalette.DARK_BLUE)
	torusMesh := graphic.NewMesh(torusGeom, mat)
	fmt.Printf("LoaderID: %s\n", displayName)
	torusMesh.SetLoaderID(displayName)
	torusMesh.SetPositionVec(vpos)
	return torusMesh
}

func (tr *TorusRenderer) NewInternalMeshAtPosition(displayName string, vpos *math32.Vector3) *graphic.Mesh {
	diskGeom := geometry.NewDisk(1, 32)
	diskMat := material.NewStandard(g3ndpalette.GREY)
	diskMesh := graphic.NewMesh(diskGeom, diskMat)
	diskMesh.SetPositionVec(vpos)
	diskMesh.SetLoaderID(displayName)
	return diskMesh
}

func (tr *TorusRenderer) NextCoordinate(prevPos *math32.Vector3) *math32.Vector3 {
	if tr.iOffset == 0 {
		tr.iOffset = 1
		return math32.NewVector3(float32(-2.0), float32(-2.0), float32(-2.0))
	} else {
		return math32.NewVector3(float32(2.0), float32(2.0), float32(2.0))
	}
}

func (tr *TorusRenderer) Layout(worldApp *g3nworld.WorldApp,
	g3nRenderableElements []*g3nmash.G3nDetailedElement) {
	tr.GenericRenderer.LayoutBase(worldApp, tr, g3nRenderableElements)
}
