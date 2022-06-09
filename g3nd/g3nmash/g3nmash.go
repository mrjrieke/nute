package g3nmash

import (
	"errors"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"tini.com/nute/mashupsdk"
)

type G3nDetailedElement struct {
	detailedElement *mashupsdk.MashupDetailedElement
	meshComposite   map[string]*graphic.Mesh // One or more meshes associated with element.
}

func NewG3nDetailedElement(detailedElement *mashupsdk.MashupDetailedElement) *G3nDetailedElement {
	return &G3nDetailedElement{detailedElement: detailedElement, meshComposite: map[string]*graphic.Mesh{}}
}

func (g *G3nDetailedElement) SetNamedMesh(meshName string, mesh *graphic.Mesh) {
	g.meshComposite[meshName] = mesh
}

func (g *G3nDetailedElement) GetDisplayId() int64 {
	return g.detailedElement.Id
}

func (g *G3nDetailedElement) GetChildElements() []int64 {
	if g.detailedElement.Childids != nil {
		return g.detailedElement.Childids
	} else {
		return []int64{}
	}
}

func (g *G3nDetailedElement) GetDisplayName() string {
	return g.detailedElement.Name
}

func (g *G3nDetailedElement) GetMashupElementState() *mashupsdk.MashupElementState {
	return g.detailedElement.State
}

func (g *G3nDetailedElement) GetDisplayState() mashupsdk.DisplayElementState {
	return mashupsdk.DisplayElementState(g.detailedElement.State.State)
}

func (g *G3nDetailedElement) SetDisplayState(x mashupsdk.DisplayElementState) {
	g.detailedElement.State.State = int64(x)
}

func (g *G3nDetailedElement) SetRotationX(x float32) error {
	if rootMesh, rootOk := g.meshComposite[g.detailedElement.Name]; rootOk {
		rootMesh.SetRotationX(x)
	}
	return errors.New("missing components")
}

func (g *G3nDetailedElement) SetColor(color *math32.Color) error {
	if rootMesh, rootOk := g.meshComposite[g.detailedElement.Name]; rootOk {
		if standardMaterial, ok := rootMesh.Graphic.GetMaterial(0).(*material.Standard); ok {
			standardMaterial.SetColor(color)
		}
	}
	return errors.New("missing components")
}
