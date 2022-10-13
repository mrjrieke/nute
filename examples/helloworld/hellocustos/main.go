//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"embed"
	"flag"
	"log"
	"os"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mrjrieke/nute/custos/custosworld"
	"github.com/mrjrieke/nute/g3nd/data"
	"github.com/mrjrieke/nute/mashupsdk"
)

var worldCompleteChan chan bool

//go:embed tls/mashup.crt
var mashupCert embed.FS

//go:embed tls/mashup.key
var mashupKey embed.FS

type ControllerRenderer struct {
	custosWorldApp *custosworld.CustosWorldApp
}

func (cr *ControllerRenderer) PreRender() {
	// TODO: buffered sort if desired.
}

func (cr *ControllerRenderer) GetPriority() int64 {
	return 1
}

func (cr *ControllerRenderer) BuildTabItem(childId int64, concreteElement *mashupsdk.MashupDetailedElement) {
	child := cr.custosWorldApp.MashupDetailedElementLibrary[childId]
	if child != nil && child.Alias != "" {
		log.Printf("Controller lookup on: %s name: %s\n", child.Alias, child.Name)
		if fwb, fwbOk := cr.custosWorldApp.FyneWidgetElements[child.Name]; fwbOk {
			if fwb.MashupDetailedElement != nil && fwb.GuiComponent != nil {
				fwb.MashupDetailedElement.Copy(child)
				fwb.GuiComponent.(*container.TabItem).Text = child.Name
			}
		} else {
			// No widget made yet for this alias...
			cr.custosWorldApp.DetailFyneComponent(child,
				BuildDetailMappedTabItemFyneComponent)
		}
	}

	if child != nil && len(child.GetChildids()) > 0 {
		for _, cId := range child.GetChildids() {
			cr.BuildTabItem(cId, concreteElement)
		}
	}
}

func (cr *ControllerRenderer) RenderTabItem(concreteElement *mashupsdk.MashupDetailedElement) {
	log.Printf("Controller Widget lookup: %s\n", concreteElement.Alias)

	if fyneWidgetElement, fyneOk := cr.custosWorldApp.FyneWidgetElements[concreteElement.Name]; fyneOk {
		log.Printf("ControllerRenderer lookup found: %s\n", concreteElement.Alias)
		if fyneWidgetElement.GuiComponent == nil {
			fyneWidgetElement.GuiComponent = cr.custosWorldApp.CustomTabItems[concreteElement.Name](cr.custosWorldApp, concreteElement.Name)
		}
		cr.custosWorldApp.TabItemMenu.Append(fyneWidgetElement.GuiComponent.(*container.TabItem))
	}
}

func (cr *ControllerRenderer) Refresh() {
	// TODO: buffered sort if desired.
}

type TorusRenderer struct {
	custosWorldApp   *custosworld.CustosWorldApp
	concreteElements []*mashupsdk.MashupDetailedElement
}

func (tr *TorusRenderer) PreRender() {
	tr.concreteElements = []*mashupsdk.MashupDetailedElement{}
}

func (tr *TorusRenderer) GetPriority() int64 {
	return 2
}

func (tr *TorusRenderer) BuildTabItem(childId int64, concreteElement *mashupsdk.MashupDetailedElement) {
	child := tr.custosWorldApp.MashupDetailedElementLibrary[childId]
	if child != nil && child.Alias != "" {
		log.Printf("TorusRenderer.BuildTabItem lookup on: %s name: %s\n", child.Alias, child.Name)
		if fwb, fwbOk := tr.custosWorldApp.FyneWidgetElements[child.Name]; fwbOk {
			if fwb.MashupDetailedElement != nil && fwb.GuiComponent != nil {
				fwb.MashupDetailedElement.Copy(child)
				fwb.GuiComponent.(*container.TabItem).Text = child.Name
			}
		} else {
			// No widget made yet for this alias...
			tr.custosWorldApp.DetailFyneComponent(child,
				BuildDetailMappedTabItemFyneComponent)
		}
	}

	if child != nil && len(child.GetChildids()) > 0 {
		for _, cId := range child.GetChildids() {
			tr.BuildTabItem(cId, concreteElement)
		}
	}
}

func (tr *TorusRenderer) renderTabItemHelper(concreteElement *mashupsdk.MashupDetailedElement) {
	log.Printf("TorusRender Widget lookup: %s\n", concreteElement.Alias)
	tr.custosWorldApp.TabItemMenu.Hide()
	if concreteElement.IsStateSet(mashupsdk.Clicked) {
		log.Printf("TorusRender Widget looking up: %s\n", concreteElement.Alias)
		if fyneWidgetElement, fyneOk := tr.custosWorldApp.FyneWidgetElements[concreteElement.Name]; fyneOk {
			log.Printf("TorusRender Widget lookup found: %s\n", concreteElement.Alias)
			if fyneWidgetElement.GuiComponent == nil {
				if customTabFun, customTabFunOk := tr.custosWorldApp.CustomTabItems[concreteElement.Name]; customTabFunOk {
					fyneWidgetElement.GuiComponent = customTabFun(tr.custosWorldApp, concreteElement.Name)
				}
			}
			if fyneWidgetElement.GuiComponent != nil {
				tr.custosWorldApp.TabItemMenu.Append(fyneWidgetElement.GuiComponent.(*container.TabItem))
			}
		}
	} else {
		// Remove it if torus.
		// CUWorldApp.fyneWidgetElements["Inside"].GuiComponent.(*container.TabItem),
		// Remove the formerly clicked elements..
		log.Printf("TorusRender Widget lookingup for remove: %s\n", concreteElement.Alias)
		if fyneWidgetElement, fyneOk := tr.custosWorldApp.FyneWidgetElements[concreteElement.Name]; fyneOk {
			log.Printf("TorusRender Widget lookup found for remove: %s %v\n", concreteElement.Alias, fyneWidgetElement)
			if fyneWidgetElement.GuiComponent != nil {
				tr.custosWorldApp.TabItemMenu.Remove(fyneWidgetElement.GuiComponent.(*container.TabItem))
			}
		}
	}
	tr.custosWorldApp.TabItemMenu.Show()
	log.Printf("End TorusRender Widget lookup: %s\n", concreteElement.Alias)
}

func (tr *TorusRenderer) RenderTabItem(concreteElement *mashupsdk.MashupDetailedElement) {
	tr.concreteElements = append(tr.concreteElements, concreteElement)
}

func (tr *TorusRenderer) Refresh() {
	sort.Slice(tr.concreteElements, func(i, j int) bool {
		return strings.Compare(tr.concreteElements[i].Name, tr.concreteElements[j].Name) == -1
	})
	for _, concreteElement := range tr.concreteElements {
		tr.renderTabItemHelper(concreteElement)
	}
}

func (tr *TorusRenderer) OnSelected(tabItem *container.TabItem) {
	// Too bad fyne doesn't have the ability for user to assign an id to TabItem...
	// Lookup by name instead and try to keep track of any name changes instead...
	log.Printf("Selected: %s\n", tabItem.Text)
	if mashupItemIndex, miOk := tr.custosWorldApp.ElementLoaderIndex[tabItem.Text]; miOk {
		if mashupDetailedElement, mOk := tr.custosWorldApp.MashupDetailedElementLibrary[mashupItemIndex]; mOk {
			if mashupDetailedElement.Name != "" {
				if mashupDetailedElement.Genre != "Collection" {
					mashupDetailedElement.State.State |= int64(mashupsdk.Clicked)
				}
				if fyneWidget, fOk := tr.custosWorldApp.FyneWidgetElements[mashupDetailedElement.Name]; fOk {
					fyneWidget.MashupDetailedElement = mashupDetailedElement
					fyneWidget.OnStatusChanged()
				} else {
					log.Printf("Unexpected widget request: %s\n", mashupDetailedElement.Name)
				}
				return
			}
		}
	}
	//CUWorldApp.fyneWidgetElements[tabItem.Text].OnStatusChanged()
}

func OutsideClone(custosWorldApp *custosworld.CustosWorldApp, childId int64, concreteElement *mashupsdk.MashupDetailedElement) {
	custosWorldApp.FyneWidgetElements["Outside"].MashupDetailedElement.Copy(concreteElement)
}

func BuildDetailMappedTabItemFyneComponent(custosWorldApp *custosworld.CustosWorldApp, id string) *container.TabItem {
	de := custosWorldApp.FyneWidgetElements[id].MashupDetailedElement

	tabLabel := widget.NewLabel(de.Description)
	tabLabel.Wrapping = fyne.TextWrapWord
	tabItem := container.NewTabItem(id, container.NewBorder(nil, nil, layout.NewSpacer(), nil, container.NewVBox(tabLabel, container.NewAdaptiveGrid(2,
		widget.NewButton("Show", func() {
			// Workaround... mashupdetailedelement points at wrong element sometimes, but shouldn't!
			if len(custosWorldApp.ElementLoaderIndex) > 0 {
				fyneWidgetBundle := custosWorldApp.FyneWidgetElements[de.Name]

				mashupIndex := custosWorldApp.ElementLoaderIndex[fyneWidgetBundle.GuiComponent.(*container.TabItem).Text]
				fyneWidgetBundle.MashupDetailedElement = custosWorldApp.MashupDetailedElementLibrary[mashupIndex]

				fyneWidgetBundle.MashupDetailedElement.ApplyState(mashupsdk.Hidden, false)
				if fyneWidgetBundle.MashupDetailedElement.Genre == "Collection" {
					fyneWidgetBundle.MashupDetailedElement.ApplyState(mashupsdk.Recursive, true)
				}
				fyneWidgetBundle.OnStatusChanged()
			}
		}), widget.NewButton("Hide", func() {
			if len(custosWorldApp.ElementLoaderIndex) > 0 {
				// Workaround... mashupdetailedelement points at wrong element sometimes, but shouldn't!
				fyneWidgetBundle := custosWorldApp.FyneWidgetElements[de.Name]
				mashupIndex := custosWorldApp.ElementLoaderIndex[fyneWidgetBundle.GuiComponent.(*container.TabItem).Text]
				fyneWidgetBundle.MashupDetailedElement = custosWorldApp.MashupDetailedElementLibrary[mashupIndex]

				fyneWidgetBundle.MashupDetailedElement.ApplyState(mashupsdk.Hidden, true)
				if fyneWidgetBundle.MashupDetailedElement.Genre == "Collection" {
					fyneWidgetBundle.MashupDetailedElement.ApplyState(mashupsdk.Recursive, true)
				}
				fyneWidgetBundle.OnStatusChanged()
			}
		})))),
	)
	return tabItem
}

//go:embed gophericon.png
var gopherIcon embed.FS

func main() {
	//runtime.LockOSThread()
	callerCreds := flag.String("CREDS", "", "Credentials of caller")
	insecure := flag.Bool("insecure", false, "Skip server validation")
	headless := flag.Bool("headless", false, "Run headless")
	titlebar := flag.Bool("titlebar", false, "Run with title bar")
	flag.Parse()
	worldLog, err := os.OpenFile("custos.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(worldLog)

	mashupsdk.InitCertKeyPair(mashupCert, mashupKey)

	detailedElements := data.GetExampleLibrary()

	torusRenderer := &TorusRenderer{}
	custosWorld := custosworld.NewCustosWorldApp(*headless, *titlebar, detailedElements, torusRenderer)
	torusRenderer.custosWorldApp = custosWorld

	// Initialize a tab item renderer
	// This will be called during upsert elements phase.
	// indexed by subgenre
	custosWorld.CustomTabItemRenderer["TabItemRenderer"] = torusRenderer
	custosWorld.CustomTabItemRenderer["ControllerTabItemRenderer"] = &ControllerRenderer{custosWorldApp: custosWorld}

	custosWorld.Title = "Hello Custos"
	custosWorld.MainWindowSize = fyne.NewSize(800, 100)
	gopherIconBytes, _ := gopherIcon.ReadFile("gophericon.png")
	custosWorld.Icon = fyne.NewStaticResource("Gopher", gopherIconBytes)

	if !custosWorld.Headless {
		custosWorld.InitServer(*callerCreds, *insecure, 0)
	}

	// Initialize the main window.
	custosWorld.InitMainWindow()

	<-worldCompleteChan
}
