package custosworld

import (
	"embed"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	"github.com/mrjrieke/nute/mashupsdk"
	"github.com/mrjrieke/nute/mashupsdk/client"
	"github.com/mrjrieke/nute/mashupsdk/guiboot"
	"github.com/mrjrieke/nute/mashupsdk/server"
)

type mashupSdkApiHandler struct {
}

type worldClientInitHandler struct {
}

type IG3nRenderer interface {
	Layout(worldApp *CustosWorldApp, g3nRenderableElements []*g3nmash.G3nDetailedElement)
	InitRenderLoop(worldApp *CustosWorldApp) bool
	RenderElement(worldApp *CustosWorldApp, g3n *g3nmash.G3nDetailedElement) bool
}

type fyneMashupApiHandler struct {
}

type CustosContext struct {
	mashupContext *mashupsdk.MashupContext // Needed for callbacks to other mashups
}

type FyneWidgetBundle struct {
	mashupsdk.GuiWidgetBundle
}

func (fwb *FyneWidgetBundle) OnStatusChanged() {
	selectedDetailedElement := fwb.MashupDetailedElement
	if CUWorldApp.HeadsupFyneContext.mashupContext == nil {
		return
	}

	elementStateBundle := mashupsdk.MashupElementStateBundle{
		AuthToken:     client.GetServerAuthToken(),
		ElementStates: []*mashupsdk.MashupElementState{selectedDetailedElement.State},
	}
	CUWorldApp.HeadsupFyneContext.mashupContext.Client.ResetG3NDetailedElementStates(CUWorldApp.HeadsupFyneContext.mashupContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})

	log.Printf("Display fields set to: %d", selectedDetailedElement.State.State)
	CUWorldApp.HeadsupFyneContext.mashupContext.Client.UpsertMashupElementsState(CUWorldApp.HeadsupFyneContext.mashupContext, &elementStateBundle)

}

type CustosWorldApp struct {
	Headless                     bool // Mode for troubleshooting.
	mashupSdkApiHandler          *mashupSdkApiHandler
	wClientInitHandler           *worldClientInitHandler
	HeadsupFyneContext           *CustosContext
	mainWin                      fyne.Window
	mashupDisplayContext         *mashupsdk.MashupDisplayContext
	DetailedElements             []*mashupsdk.MashupDetailedElement
	mashupDetailedElementLibrary map[int64]*mashupsdk.MashupDetailedElement
	elementLoaderIndex           map[string]int64 // mashup indexes by Name
	fyneWidgetElements           map[string]*FyneWidgetBundle
	ClickedElements              []*mashupsdk.MashupDetailedElement // g3n indexes by string...
}

//go:embed gophericon.png
var gopherIcon embed.FS

var CUWorldApp *CustosWorldApp

func (w *CustosWorldApp) InitServer(callerCreds string, insecure bool) {
	if callerCreds != "" {
		server.InitServer(callerCreds, insecure, w.mashupSdkApiHandler, w.wClientInitHandler)
	} else {
		// TODO: These might not make sense in Custos.
		// go func() {
		// 	w.displaySetupChan <- &mashupsdk.MashupDisplayHint{Xpos: 0, Ypos: 0, Width: 400, Height: 800}
		// }()
	}
}

func NewCustosWorldApp(headless bool, detailedElements []*mashupsdk.MashupDetailedElement, renderer IG3nRenderer) *CustosWorldApp {
	CUWorldApp = &CustosWorldApp{
		Headless:                     headless,
		mashupSdkApiHandler:          &mashupSdkApiHandler{},
		HeadsupFyneContext:           &CustosContext{},
		DetailedElements:             detailedElements,
		mainWin:                      nil,
		mashupDisplayContext:         &mashupsdk.MashupDisplayContext{MainWinDisplay: &mashupsdk.MashupDisplayHint{}},
		mashupDetailedElementLibrary: map[int64]*mashupsdk.MashupDetailedElement{}, // mashupDetailedElementLibrary,
		elementLoaderIndex:           map[string]int64{},                           // elementLoaderIndex
		fyneWidgetElements: map[string]*FyneWidgetBundle{
			"Inside": {
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent:          nil,
					MashupDetailedElement: &mashupsdk.MashupDetailedElement{}, // mashupDetailedElementLibrary["{0}-AxialCircle"],
				},
			},
			"Outside": {
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent:          nil,
					MashupDetailedElement: &mashupsdk.MashupDetailedElement{}, //mashupDetailedElementLibrary["Outside"],
				},
			},
			"It": {
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent:          nil,
					MashupDetailedElement: &mashupsdk.MashupDetailedElement{}, //mashupDetailedElementLibrary["{0}-Torus"],
				},
			},
			"Up-Side-Down": {
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent:          nil,
					MashupDetailedElement: &mashupsdk.MashupDetailedElement{}, //mashupDetailedElementLibrary["{0}-SharedAttitude"],
				},
			},
			"All": {
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent:          nil,
					MashupDetailedElement: &mashupsdk.MashupDetailedElement{}, //mashupDetailedElementLibrary["{0}-SharedAttitude"],
				},
			},
		},
	}

	return CUWorldApp
}

type InitEvent struct {
}

func (w *CustosWorldApp) ResetChangeStates() []*mashupsdk.MashupElementState {
	changedElements := []*mashupsdk.MashupElementState{}
	for _, g3nDetailedElement := range w.mashupDetailedElementLibrary {
		if mashupsdk.DisplayElementState(g3nDetailedElement.GetMashupElementState().State) != mashupsdk.Init {
			g3nDetailedElement.ApplyState(mashupsdk.Clicked, false)
			changedElements = append(changedElements, g3nDetailedElement.GetMashupElementState())
		}
	}

	return changedElements
}

func (w *CustosWorldApp) InitMainWindow() {
	log.Printf("Initializing mainWin.")

	initHandler := func(a fyne.App) {
		log.Printf("InitHandler.")
		CUWorldApp.mainWin = a.NewWindow("Hello Custos")
		gopherIconBytes, _ := gopherIcon.ReadFile("gophericon.png")

		CUWorldApp.mainWin.SetIcon(fyne.NewStaticResource("Gopher", gopherIconBytes))
		CUWorldApp.mainWin.Resize(fyne.NewSize(800, 100))
		CUWorldApp.mainWin.SetFixedSize(false)

		CUWorldApp.fyneWidgetElements["Inside"].GuiComponent = CUWorldApp.detailMappedFyneComponent("Inside", "The magnetic field inside a toroid is always tangential to the circular closed path.  These magnetic field lines are concentric circles.", CUWorldApp.fyneWidgetElements["Inside"].MashupDetailedElement)
		CUWorldApp.fyneWidgetElements["Outside"].GuiComponent = CUWorldApp.detailMappedFyneComponent("Outside", "The magnetic field at any point outside the toroid is zero.", CUWorldApp.fyneWidgetElements["Outside"].MashupDetailedElement)
		CUWorldApp.fyneWidgetElements["It"].GuiComponent = CUWorldApp.detailMappedFyneComponent("It", "The magnetic field inside the empty space surrounded by the toroid is zero.", CUWorldApp.fyneWidgetElements["It"].MashupDetailedElement)
		CUWorldApp.fyneWidgetElements["Up-Side-Down"].GuiComponent = CUWorldApp.detailMappedFyneComponent("Up-Side-Down", "Torus is up-side-down", CUWorldApp.fyneWidgetElements["Up-Side-Down"].MashupDetailedElement)
		CUWorldApp.fyneWidgetElements["All"].GuiComponent = CUWorldApp.detailMappedFyneComponent("All", "A group of torus or a tori.", CUWorldApp.fyneWidgetElements["All"].MashupDetailedElement)

		torusMenu := container.NewAppTabs(
			CUWorldApp.fyneWidgetElements["Inside"].GuiComponent.(*container.TabItem),
			CUWorldApp.fyneWidgetElements["Outside"].GuiComponent.(*container.TabItem),
			CUWorldApp.fyneWidgetElements["It"].GuiComponent.(*container.TabItem),
			CUWorldApp.fyneWidgetElements["Up-Side-Down"].GuiComponent.(*container.TabItem),
			CUWorldApp.fyneWidgetElements["All"].GuiComponent.(*container.TabItem),
		)
		torusMenu.OnSelected = func(tabItem *container.TabItem) {
			// Too bad fyne doesn't have the ability for user to assign an id to TabItem...
			// Lookup by name instead and try to keep track of any name changes instead...
			log.Printf("Selected: %s\n", tabItem.Text)
			if mashupItemIndex, miOk := CUWorldApp.elementLoaderIndex[tabItem.Text]; miOk {
				if mashupDetailedElement, mOk := CUWorldApp.mashupDetailedElementLibrary[mashupItemIndex]; mOk {
					if mashupDetailedElement.Alias != "" {
						if mashupDetailedElement.Genre != "Collection" {
							mashupDetailedElement.State.State |= int64(mashupsdk.Clicked)
						}
						CUWorldApp.fyneWidgetElements[mashupDetailedElement.Alias].MashupDetailedElement = mashupDetailedElement
						CUWorldApp.fyneWidgetElements[mashupDetailedElement.Alias].OnStatusChanged()
						return
					}
				}
			}
			CUWorldApp.fyneWidgetElements[tabItem.Text].OnStatusChanged()
		}

		torusMenu.SetTabLocation(container.TabLocationTop)

		CUWorldApp.mainWin.SetContent(torusMenu)
		CUWorldApp.mainWin.SetCloseIntercept(func() {
			if CUWorldApp.HeadsupFyneContext.mashupContext != nil {
				CUWorldApp.HeadsupFyneContext.mashupContext.Client.Shutdown(CUWorldApp.HeadsupFyneContext.mashupContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})
			}
			os.Exit(0)
		})
	}
	runtimeHandler := func() {
		if w.mainWin != nil {
			log.Printf("CustosWorld main win initialized\n")
			if CUWorldApp.mashupDisplayContext != nil &&
				(CUWorldApp.Headless ||
					(CUWorldApp.mashupDisplayContext.GetSettled()&mashupsdk.AppInitted) == mashupsdk.AppInitted) {
				log.Printf("CustosWorld app settled... starting up.\n")
				w.mainWin.ShowAndRun()
			} else {
				if !CUWorldApp.Headless {
					w.mainWin.Hide()
				}
			}
		}
	}

	guiboot.InitMainWindow(guiboot.Fyne, initHandler, runtimeHandler)
}

func (w *worldClientInitHandler) RegisterContext(context *mashupsdk.MashupContext) {
	CUWorldApp.HeadsupFyneContext.mashupContext = context
}

// Sets all elements to a "Rest state."
func (w *mashupSdkApiHandler) ResetG3NDetailedElementStates() {
	log.Printf("CustosWorld Received ResetG3NDetailedElementStates\n")
	for _, wes := range CUWorldApp.mashupDetailedElementLibrary {
		wes.SetElementState(mashupsdk.Init)
	}
	log.Printf("CustosWorld finished ResetG3NDetailedElementStates handle.\n")
}

func (mSdk *mashupSdkApiHandler) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	log.Printf("CustosWorld OnResize - not implemented yet..\n")
	if CUWorldApp.mainWin != nil && CUWorldApp.mashupDisplayContext != nil && CUWorldApp.mashupDisplayContext.MainWinDisplay != nil {
		CUWorldApp.mashupDisplayContext.MainWinDisplay.Focused = displayHint.Focused
		// TODO: Resize without infinite looping....
		// The moment fyne is resized, it'll want to resize g3n...
		// Which then wants to resize fyne ad-infinitum
		//CUWorldApp.mainWin.PosResize(int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height))
		log.Printf("CustosWorld Received onResize xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	} else {
		log.Printf("CustosWorld Could not apply xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	}
}

func (CustosWorldApp *CustosWorldApp) detailMappedFyneComponent(id, description string, de *mashupsdk.MashupDetailedElement) *container.TabItem {
	tabLabel := widget.NewLabel(description)
	tabLabel.Wrapping = fyne.TextWrapWord
	tabItem := container.NewTabItem(id, container.NewBorder(nil, nil, layout.NewSpacer(), nil, container.NewVBox(tabLabel, container.NewAdaptiveGrid(2,
		widget.NewButton("Show", func() {
			// Workaround... mashupdetailedelement points at wrong element sometimes, but shouldn't!
			if len(CUWorldApp.elementLoaderIndex) > 0 {
				mashupIndex := CUWorldApp.elementLoaderIndex[CUWorldApp.fyneWidgetElements[de.Alias].GuiComponent.(*container.TabItem).Text]
				CUWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement = CUWorldApp.mashupDetailedElementLibrary[mashupIndex]

				CUWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Hidden, false)
				if CUWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.Genre == "Collection" {
					CUWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Recursive, true)
				}
				CUWorldApp.fyneWidgetElements[de.Alias].OnStatusChanged()
			}
		}), widget.NewButton("Hide", func() {
			if len(CUWorldApp.elementLoaderIndex) > 0 {
				// Workaround... mashupdetailedelement points at wrong element sometimes, but shouldn't!
				mashupIndex := CUWorldApp.elementLoaderIndex[CUWorldApp.fyneWidgetElements[de.Alias].GuiComponent.(*container.TabItem).Text]
				CUWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement = CUWorldApp.mashupDetailedElementLibrary[mashupIndex]

				CUWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Hidden, true)
				if CUWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.Genre == "Collection" {
					CUWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Recursive, true)
				}
				CUWorldApp.fyneWidgetElements[de.Alias].OnStatusChanged()
			}
		})))),
	)
	return tabItem
}

func (CustosWorldApp *CustosWorldApp) TorusParser(childId int64) {
	child := CUWorldApp.mashupDetailedElementLibrary[childId]
	if child != nil && child.Alias != "" {
		log.Printf("TorusParser lookup on: %s\n", child.Alias)
		if CUWorldApp.fyneWidgetElements != nil && CUWorldApp.fyneWidgetElements[child.Alias].MashupDetailedElement != nil && CUWorldApp.fyneWidgetElements[child.Alias].GuiComponent != nil {
			CUWorldApp.fyneWidgetElements[child.Alias].MashupDetailedElement.Copy(child)
			CUWorldApp.fyneWidgetElements[child.Alias].GuiComponent.(*container.TabItem).Text = child.Name
		}
	}

	if child != nil && len(child.GetChildids()) > 0 {
		for _, cId := range child.GetChildids() {
			CUWorldApp.TorusParser(cId)
		}
	}
}

func (mSdk *mashupSdkApiHandler) GetMashupElements() (*mashupsdk.MashupDetailedElementBundle, error) {
	log.Printf("CustosWorld Received GetMashupElements\n")

	return &mashupsdk.MashupDetailedElementBundle{
		AuthToken:        client.GetServerAuthToken(),
		DetailedElements: CUWorldApp.DetailedElements,
	}, nil
}

func (mSdk *mashupSdkApiHandler) UpsertMashupElements(detailedElementBundle *mashupsdk.MashupDetailedElementBundle) (*mashupsdk.MashupDetailedElementBundle, error) {
	log.Printf("CustosWorld Received UpsertMashupElements\n")

	for _, concreteElement := range detailedElementBundle.DetailedElements {
		//helloApp.fyneComponentCache[generatedComponent.Basisid]
		CUWorldApp.mashupDetailedElementLibrary[concreteElement.Id] = concreteElement
		CUWorldApp.elementLoaderIndex[concreteElement.Name] = concreteElement.Id

		if concreteElement.GetName() == "Outside" {
			CUWorldApp.fyneWidgetElements["Outside"].MashupDetailedElement.Copy(concreteElement)
		}
	}
	log.Printf("CustosWorld parsing tori.\n")

	for _, concreteElement := range detailedElementBundle.DetailedElements {
		if concreteElement.GetSubgenre() == "Torus" {
			CUWorldApp.TorusParser(concreteElement.Id)
		}
	}

	log.Printf("Mashup elements delivered.\n")

	CUWorldApp.mashupDisplayContext.ApplySettled(mashupsdk.AppInitted, false)

	log.Printf("CustosWorld UpsertMashupElements updated\n")
	return &mashupsdk.MashupDetailedElementBundle{
		AuthToken:        client.GetServerAuthToken(),
		DetailedElements: detailedElementBundle.DetailedElements,
	}, nil
}
func (mSdk *mashupSdkApiHandler) setStateHelper(g3nId int64, x mashupsdk.DisplayElementState) {

	child := CUWorldApp.mashupDetailedElementLibrary[g3nId]
	if child.Genre != "Attitude" {
		child.SetElementState(mashupsdk.DisplayElementState(x))
	}

	if len(child.Childids) > 0 {
		for _, cId := range child.Childids {
			mSdk.setStateHelper(cId, x)
		}
	}
}

func (mSdk *mashupSdkApiHandler) UpsertMashupElementsState(elementStateBundle *mashupsdk.MashupElementStateBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("CustosWorld UpsertMashupElementsState called\n")

	ClickedElements := map[int64]*mashupsdk.MashupDetailedElement{}
	recursiveElements := map[int64]*mashupsdk.MashupDetailedElement{}

	for _, es := range elementStateBundle.ElementStates {
		if g3nDetailedElement, ok := CUWorldApp.mashupDetailedElementLibrary[es.GetId()]; ok {
			g3nDetailedElement.SetElementState(mashupsdk.DisplayElementState(es.State))
			if g3nDetailedElement.IsStateSet(mashupsdk.Recursive) {
				recursiveElements[es.GetId()] = g3nDetailedElement
			}

			log.Printf("Display fields set to: %d", g3nDetailedElement.GetMashupElementState())
			if (mashupsdk.DisplayElementState(es.State) & mashupsdk.Clicked) == mashupsdk.Clicked {
				ClickedElements[es.GetId()] = g3nDetailedElement
			}
		}
	}

	if len(ClickedElements) > 0 {
		// Remove existing clicks.
		for _, clickedElement := range CUWorldApp.ClickedElements {
			if _, ok := ClickedElements[clickedElement.GetId()]; !ok {
				clickedElement.ApplyState(mashupsdk.Clicked, false)
			}
		}

		CUWorldApp.ClickedElements = CUWorldApp.ClickedElements[:0]

		// Impossible to determine ordering of clicks from upsert at this time.
		for _, g3nDetailedElement := range ClickedElements {
			CUWorldApp.ClickedElements = append(CUWorldApp.ClickedElements, g3nDetailedElement)
		}
	}

	if len(recursiveElements) > 0 {
		for _, recursiveElement := range recursiveElements {
			stateBits := recursiveElement.State.State
			// Unset recursive for child elements
			stateBits &= ^int64(mashupsdk.Recursive)
			// Apply this state change to all child elements.
			mSdk.setStateHelper(recursiveElement.GetId(), mashupsdk.DisplayElementState(stateBits))
		}
	}

	log.Printf("CustosWorld dispatching focus\n")
	// if CUWorldApp.mainWin != nil {
	// 	CUWorldApp.mainWin.Hide()
	// }
	log.Printf("CustosWorld End UpsertMashupElementsState called\n")
	return &mashupsdk.MashupElementStateBundle{}, nil
}
