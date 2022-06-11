package g3nmash

import (
	"errors"
	"strconv"
	"strings"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"tini.com/nute/mashupsdk"
)

type G3nDetailedElement struct {
	detailedElement *mashupsdk.MashupDetailedElement
	meshComposite   map[string]*graphic.Mesh // One or more meshes associated with element.
	attitudes       []float32
}

func NewG3nDetailedElement(detailedElement *mashupsdk.MashupDetailedElement) *G3nDetailedElement {
	g3n := G3nDetailedElement{detailedElement: detailedElement, meshComposite: map[string]*graphic.Mesh{}}
	if detailedElement.GetGenre() == "Attitude" {
		attitudes := detailedElement.GetSubgenre()
		attutudeSlice := strings.Split(attitudes, ",")
		for _, attitude := range attutudeSlice {
			if a, aErr := strconv.ParseFloat(attitude, 32); aErr == nil {
				g3n.attitudes = append(g3n.attitudes, float32(a))
			}
		}
	}
	return &g3n
}

func (g *G3nDetailedElement) SetNamedMesh(meshName string, mesh *graphic.Mesh) {
	g.meshComposite[meshName] = mesh
}

func (g *G3nDetailedElement) GetDisplayId() int64 {
	return g.detailedElement.Id
}

func (g *G3nDetailedElement) IsBackground() bool {
	return g.detailedElement.Genre == "Space" && g.detailedElement.Subgenre == "Exo"
}

func (g *G3nDetailedElement) HasAttitudeAdjustment() bool {
	return g.detailedElement.Genre == "Attitude"
}

func (g *G3nDetailedElement) AdjustAttitude(parentG3Elements []*G3nDetailedElement) error {
	if g.HasAttitudeAdjustment() {
		switch len(g.attitudes) {
		case 1:
			g.SetRotation(parentG3Elements, g.attitudes[0], 0, 0)
		case 2:
			g.SetRotation(parentG3Elements, g.attitudes[0], g.attitudes[1], 0)
		case 3:
			g.SetRotation(parentG3Elements, g.attitudes[0], g.attitudes[1], g.attitudes[2])
		}
	}
	return errors.New("no adjustment")
}

// TODO: Find a better name for this.
func (g *G3nDetailedElement) IsComposite() bool {
	return len(g.detailedElement.Parentids) == 0
}

func (g *G3nDetailedElement) IsItemActive() bool {
	return g.GetDisplayState() != mashupsdk.Rest
}

func (g *G3nDetailedElement) IsItemClicked(node core.INode) bool {
	if node == nil {
		return false
	} else {
		return node.GetNode().LoaderID() == g.detailedElement.Name
	}
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

func (g *G3nDetailedElement) SetDisplayState(x mashupsdk.DisplayElementState) bool {
	if g.detailedElement.State.State != int64(x) {
		g.detailedElement.State.State = int64(x)
		return true
	}
	return false
}

func (g *G3nDetailedElement) SetRotationX(x float32) error {
	if rootMesh, rootOk := g.meshComposite[g.detailedElement.Name]; rootOk {
		rootMesh.SetRotationX(x)
		return nil
	}
	return errors.New("missing components")
}

func (g *G3nDetailedElement) SetRotation(parentG3Elements []*G3nDetailedElement, x float32, y float32, z float32) error {
	if len(g.detailedElement.GetChildids()) > 0 {
		for _, parentG3Element := range parentG3Elements {
			if rootMesh, rootOk := parentG3Element.meshComposite[g.detailedElement.Name]; rootOk { // Hello friend.
				rootMesh.SetRotationX(x)
				rootMesh.SetRotationY(y)
				rootMesh.SetRotationZ(z)
			}
		}

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
