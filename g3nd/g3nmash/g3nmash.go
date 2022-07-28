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
	"github.com/mrjrieke/nute/mashupsdk"
)

type G3nDetailedElement struct {
	detailedElement *mashupsdk.MashupDetailedElement
	meshComposite   map[string]core.INode // One or more meshes associated with element.
	color           *math32.Color
	attitudes       []float32
}

func NewG3nDetailedElement(detailedElement *mashupsdk.MashupDetailedElement, deepCopy bool) *G3nDetailedElement {
	detailedRef := detailedElement
	if deepCopy {
		detailedRef = &mashupsdk.MashupDetailedElement{
			Basisid:       detailedElement.Basisid,
			Id:            detailedElement.Id,
			State:         &mashupsdk.MashupElementState{Id: detailedElement.State.Id, State: detailedElement.State.State},
			Name:          detailedElement.Name,
			Alias:         detailedElement.Alias,
			Description:   detailedElement.Description,
			Renderer:      detailedElement.Renderer,
			Colabrenderer: detailedElement.Colabrenderer,
			Genre:         detailedElement.Genre,
			Subgenre:      detailedElement.Subgenre,
			Parentids:     detailedElement.Parentids,
			Childids:      detailedElement.Childids,
		}
	}

	if detailedElement.Id > 0 {
		detailedElement.State.Id = detailedElement.Id
	}

	g3n := G3nDetailedElement{detailedElement: detailedRef, meshComposite: map[string]core.INode{}}
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

func CloneG3nDetailedElement(
	getG3nDetailedElementById func(eid int64) (*G3nDetailedElement, error),
	getG3nDetailedLibraryElementById func(eid int64) (*G3nDetailedElement, error),
	indexG3nDetailedElement func(*G3nDetailedElement) *G3nDetailedElement,
	newIdPumpFunc func() int64,
	g3nElement *G3nDetailedElement,
	generatedElements *[]interface{},
) *G3nDetailedElement {
	g3n := NewG3nDetailedElement(g3nElement.detailedElement, true)
	if g3n.detailedElement.Basisid < 0 {
		// Convert from library to generated.
		g3n.detailedElement.Id = newIdPumpFunc()
		// Id state must match detailed element id.
		g3n.detailedElement.State.Id = g3n.detailedElement.Id
		g3n.detailedElement.Name = strings.Replace(g3n.detailedElement.Name, "{0}", strconv.FormatInt(g3n.detailedElement.Id, 10), 1)
		// Converted from mutable to instance...
		// Upgrade state to Init.
		g3n.SetDisplayState(mashupsdk.Init)
		*generatedElements = append(*generatedElements, g3n.GetDetailedElement())
	}

	newChildIds := []int64{}
	for _, childId := range g3n.GetChildElementIds() {
		if childId < 0 {
			if libElement, err := getG3nDetailedLibraryElementById(childId); err == nil {
				clonedChildElement := CloneG3nDetailedElement(getG3nDetailedElementById, getG3nDetailedLibraryElementById, indexG3nDetailedElement, newIdPumpFunc, libElement, generatedElements)
				clonedChildElement.SetParentElements([]int64{g3n.GetDisplayId()})
				newChildIds = append(newChildIds, clonedChildElement.GetDisplayId())
			} else {
				log.Printf("Missing child from library: %d\n", childId)
			}
		} else {
			// Deal with concrete element.
			if concreteElement, err := getG3nDetailedElementById(childId); err == nil {
				newChildIds = append(newChildIds, concreteElement.GetDisplayId())
				existingParents := concreteElement.GetParentElementIds()
				if len(existingParents) == 0 {
					existingParents = []int64{}
				}
				concreteElement.SetParentElements(append(existingParents, g3n.GetDisplayId()))
			}
		}
	}
	if len(newChildIds) > 0 {
		g3n.SetChildElements(newChildIds)
	}
	indexG3nDetailedElement(g3n)

	return g3n
}

func (g *G3nDetailedElement) SetNamedMesh(meshName string, mesh core.INode) {
	g.meshComposite[meshName] = mesh
}

func (g *G3nDetailedElement) GetNamedMesh(meshName string) core.INode {
	return g.meshComposite[meshName]
}

func (g *G3nDetailedElement) GetDetailedElement() *mashupsdk.MashupDetailedElement {
	return g.detailedElement
}

func (g *G3nDetailedElement) GetBasisId() int64 {
	return g.detailedElement.Basisid
}

func (g *G3nDetailedElement) GetDisplayId() int64 {
	return g.detailedElement.Id
}

func (g *G3nDetailedElement) IsAbstract() bool {
	return g.detailedElement.Genre == "Abstract"
}

func (g *G3nDetailedElement) IsBackground() bool {
	return g.detailedElement.Genre == "Space" && g.detailedElement.Subgenre == "Exo"
}

func (g *G3nDetailedElement) IsBackgroundElement() bool {
	return g.detailedElement.Genre == "Space"
}

func (g *G3nDetailedElement) HasGenre(genre string) bool {
	return g.detailedElement.Genre == genre
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
	return g.detailedElement.Basisid < 0 && g.detailedElement.Id == 0
}

// TODO: Find a better name for this.
func (g *G3nDetailedElement) IsComposite() bool {
	return len(g.detailedElement.Parentids) == 0
}

func (g *G3nDetailedElement) IsItemActive() bool {
	displayState := g.GetDisplayState()
	return displayState&mashupsdk.Rest == 0
}

func (g *G3nDetailedElement) IsItemClicked(itemClicked core.INode) bool {
	if itemClicked == nil {
		return false
	} else {
		return itemClicked.GetNode().LoaderID() == g.detailedElement.Name
	}
}

func (g *G3nDetailedElement) GetChildElementIds() []int64 {
	if g.detailedElement.Childids != nil {
		return g.detailedElement.Childids
	} else {
		return []int64{}
	}
}

func (g *G3nDetailedElement) SetChildElements(childIds []int64) {
	g.detailedElement.Childids = childIds
}

func (g *G3nDetailedElement) SetParentElements(parentIds []int64) {
	g.detailedElement.Parentids = parentIds
}

func (g *G3nDetailedElement) GetParentElementIds() []int64 {
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
	if (x & mashupsdk.Rest) == mashupsdk.Rest {
		if (g.detailedElement.State.State & int64(mashupsdk.Clicked)) == int64(mashupsdk.Clicked) {
			g.detailedElement.State.State &= ^int64(mashupsdk.Clicked)
		}
	} else if (x & mashupsdk.Clicked) == mashupsdk.Clicked {
		if (g.detailedElement.State.State & int64(mashupsdk.Rest)) == int64(mashupsdk.Rest) {
			g.detailedElement.State.State &= ^int64(mashupsdk.Rest)
		}
	}

	if (g.detailedElement.State.State & int64(x)) != int64(x) {
		g.detailedElement.State.State |= int64(x)
		return true
	}

	return false
}

func (g *G3nDetailedElement) SetRotationX(x float32) error {
	if rootMesh, rootOk := g.meshComposite[g.detailedElement.Name]; rootOk {
		if graphicMesh, isGraphicMesh := rootMesh.(*graphic.Mesh); isGraphicMesh {
			graphicMesh.SetRotationX(x)
		}
		return nil
	}
	return errors.New("missing components")
}

func (g *G3nDetailedElement) ApplyRotation(parentG3Elements []*G3nDetailedElement, x float32, y float32, z float32) error {
	log.Printf("Apply rotation: %d\n", len(parentG3Elements))
	for _, parentG3Element := range parentG3Elements {
		if rootMesh, rootOk := parentG3Element.meshComposite[parentG3Element.detailedElement.Name]; rootOk { // Hello friend.
			log.Printf("Apply rotation: %f %f %f\n", x, y, z)
			if graphicMesh, isGraphicMesh := rootMesh.(*graphic.Mesh); isGraphicMesh {
				graphicMesh.SetRotationX(x)
				graphicMesh.SetRotationY(y)
				graphicMesh.SetRotationZ(z)
			}
		}
	}
	return errors.New("missing components")
}

func (g *G3nDetailedElement) SetColor(color *math32.Color, opacity float32) bool {
	g.color = color
	if g.IsBackground() {
		return true
	}
	// TODO: iterate and set???
	if rootMesh, rootOk := g.meshComposite[g.detailedElement.Name]; rootOk {

		if graphicMesh, isGraphicMesh := rootMesh.(*graphic.Mesh); isGraphicMesh {
			if standardMaterial, ok := graphicMesh.GetMaterial(0).(*material.Standard); ok {
				standardMaterial.SetOpacity(opacity)
				ambient := standardMaterial.AmbientColor()
				if !color.Equals(&ambient) {
					standardMaterial.SetColor(color)
					return true
				}
			}
		}
	}
	return false
}

func (g *G3nDetailedElement) GetColor() *math32.Color {
	return g.color
}
