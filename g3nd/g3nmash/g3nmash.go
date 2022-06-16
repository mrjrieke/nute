package g3nmash

import (
	"errors"
	"log"
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
		log.Printf("We have attitude: %v\n", g3n.attitudes)
	}
	return &g3n
}

func CloneG3nDetailedElement(newId int64, g3nElement *G3nDetailedElement) *G3nDetailedElement {
	g3n := NewG3nDetailedElement(g3nElement.detailedElement)
	g3n.detailedElement.Id = newId
	return g3n
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

func (g *G3nDetailedElement) IsBackgroundColor() bool {
	return g.detailedElement.Genre == "Space"
}

func (g *G3nDetailedElement) HasAttitudeAdjustment() bool {
	return g.detailedElement.Genre == "Attitude"
}

func (g *G3nDetailedElement) AdjustAttitude(parentG3Elements []*G3nDetailedElement) error {
	if g.HasAttitudeAdjustment() {
		switch len(g.attitudes) {
		case 1:
			return g.ApplyRotation(parentG3Elements, g.attitudes[0], 0, 0)
		case 2:
			return g.ApplyRotation(parentG3Elements, g.attitudes[0], g.attitudes[1], 0)
		case 3:
			return g.ApplyRotation(parentG3Elements, g.attitudes[0], g.attitudes[1], g.attitudes[2])
		}
	} else {
		return g.ApplyRotation(parentG3Elements, 0.0, 0.0, 0.0)
	}
	return errors.New("no adjustment")
}

func (g *G3nDetailedElement) IsLibraryElement() bool {
	return g.detailedElement.State.State == int64(mashupsdk.Mutable)
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

func (g *G3nDetailedElement) SetChildElements(childIds []int64) {
	g.detailedElement.Childids = childIds
}

func (g *G3nDetailedElement) GetParentElements() []int64 {
	if g.detailedElement.Parentids != nil {
		return g.detailedElement.Parentids
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

func (g *G3nDetailedElement) ApplyRotation(parentG3Elements []*G3nDetailedElement, x float32, y float32, z float32) error {
	log.Printf("Apply rotation: %d\n", len(parentG3Elements))
	for _, parentG3Element := range parentG3Elements {
		if rootMesh, rootOk := parentG3Element.meshComposite[parentG3Element.detailedElement.Name]; rootOk { // Hello friend.
			log.Printf("Apply rotation: %f %f %f\n", x, y, z)
			rootMesh.SetRotationX(x)
			rootMesh.SetRotationY(y)
			rootMesh.SetRotationZ(z)
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
