package hfworld

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
	Layout(worldApp *HFWorldApp, g3nRenderableElements []*g3nmash.G3nDetailedElement)
	InitRenderLoop(worldApp *HFWorldApp) bool
	RenderElement(worldApp *HFWorldApp, g3n *g3nmash.G3nDetailedElement) bool
}

type fyneMashupApiHandler struct {
}

type HFContext struct {
	mashupContext *mashupsdk.MashupContext // Needed for callbacks to other mashups
}

type FyneWidgetBundle struct {
	mashupsdk.GuiWidgetBundle
}

func (fwb *FyneWidgetBundle) OnStatusChanged() {
	selectedDetailedElement := fwb.MashupDetailedElement
	if hfWorldApp.HeadsupFyneContext.mashupContext == nil {
		return
	}

	elementStateBundle := mashupsdk.MashupElementStateBundle{
		AuthToken:     client.GetServerAuthToken(),
		ElementStates: []*mashupsdk.MashupElementState{selectedDetailedElement.State},
	}
	hfWorldApp.HeadsupFyneContext.mashupContext.Client.ResetG3NDetailedElementStates(hfWorldApp.HeadsupFyneContext.mashupContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})

	log.Printf("Display fields set to: %d", selectedDetailedElement.State.State)
	hfWorldApp.HeadsupFyneContext.mashupContext.Client.UpsertMashupElementsState(hfWorldApp.HeadsupFyneContext.mashupContext, &elementStateBundle)

}

type HFWorldApp struct {
	mashupSdkApiHandler          *mashupSdkApiHandler
	wClientInitHandler           *worldClientInitHandler
	HeadsupFyneContext           *HFContext
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

var hfWorldApp *HFWorldApp

func (w *HFWorldApp) InitServer(callerCreds string, insecure bool) {
	if callerCreds != "" {
		server.InitServer(callerCreds, insecure, w.mashupSdkApiHandler, w.wClientInitHandler)
	} else {
		// TODO: These might not make sense in HF.
		// go func() {
		// 	w.displaySetupChan <- &mashupsdk.MashupDisplayHint{Xpos: 0, Ypos: 0, Width: 400, Height: 800}
		// }()
	}
}

func NewHFWorldApp(headless bool, detailedElements []*mashupsdk.MashupDetailedElement, renderer IG3nRenderer) *HFWorldApp {

	hfWorldApp = &HFWorldApp{
		mashupSdkApiHandler:          &mashupSdkApiHandler{},
		HeadsupFyneContext:           &HFContext{},
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
	return hfWorldApp
}

type InitEvent struct {
}

func (w *HFWorldApp) ResetChangeStates() []*mashupsdk.MashupElementState {
	changedElements := []*mashupsdk.MashupElementState{}
	for _, g3nDetailedElement := range w.mashupDetailedElementLibrary {
		if mashupsdk.DisplayElementState(g3nDetailedElement.GetMashupElementState().State) != mashupsdk.Init {
			g3nDetailedElement.ApplyState(mashupsdk.Clicked, false)
			changedElements = append(changedElements, g3nDetailedElement.GetMashupElementState())
		}
	}

	return changedElements
}

func (w *HFWorldApp) InitMainWindow() {
	log.Printf("Initializing mainWin.")

	initHandler := func(a fyne.App) {
		log.Printf("InitHandler.")
		hfWorldApp.mainWin = a.NewWindow("Hello Fyne Headsup")
		gopherIconBytes, _ := gopherIcon.ReadFile("gophericon.png")

		hfWorldApp.mainWin.SetIcon(fyne.NewStaticResource("Gopher", gopherIconBytes))
		hfWorldApp.mainWin.Resize(fyne.NewSize(800, 100))
		hfWorldApp.mainWin.SetFixedSize(false)

		hfWorldApp.fyneWidgetElements["Inside"].GuiComponent = hfWorldApp.detailMappedFyneComponent("Inside", "The magnetic field inside a toroid is always tangential to the circular closed path.  These magnetic field lines are concentric circles.", hfWorldApp.fyneWidgetElements["Inside"].MashupDetailedElement)
		hfWorldApp.fyneWidgetElements["Outside"].GuiComponent = hfWorldApp.detailMappedFyneComponent("Outside", "The magnetic field at any point outside the toroid is zero.", hfWorldApp.fyneWidgetElements["Outside"].MashupDetailedElement)
		hfWorldApp.fyneWidgetElements["It"].GuiComponent = hfWorldApp.detailMappedFyneComponent("It", "The magnetic field inside the empty space surrounded by the toroid is zero.", hfWorldApp.fyneWidgetElements["It"].MashupDetailedElement)
		hfWorldApp.fyneWidgetElements["Up-Side-Down"].GuiComponent = hfWorldApp.detailMappedFyneComponent("Up-Side-Down", "Torus is up-side-down", hfWorldApp.fyneWidgetElements["Up-Side-Down"].MashupDetailedElement)
		hfWorldApp.fyneWidgetElements["All"].GuiComponent = hfWorldApp.detailMappedFyneComponent("All", "A group of torus or a tori.", hfWorldApp.fyneWidgetElements["All"].MashupDetailedElement)

		torusMenu := container.NewAppTabs(
			hfWorldApp.fyneWidgetElements["Inside"].GuiComponent.(*container.TabItem),
			hfWorldApp.fyneWidgetElements["Outside"].GuiComponent.(*container.TabItem),
			hfWorldApp.fyneWidgetElements["It"].GuiComponent.(*container.TabItem),
			hfWorldApp.fyneWidgetElements["Up-Side-Down"].GuiComponent.(*container.TabItem),
			hfWorldApp.fyneWidgetElements["All"].GuiComponent.(*container.TabItem),
		)
		torusMenu.OnSelected = func(tabItem *container.TabItem) {
			// Too bad fyne doesn't have the ability for user to assign an id to TabItem...
			// Lookup by name instead and try to keep track of any name changes instead...
			log.Printf("Selected: %s\n", tabItem.Text)
			if mashupItemIndex, miOk := hfWorldApp.elementLoaderIndex[tabItem.Text]; miOk {
				mashupDetailedElement := hfWorldApp.mashupDetailedElementLibrary[mashupItemIndex]
				if mashupDetailedElement.Alias != "" {
					if mashupDetailedElement.Genre != "Collection" {
						mashupDetailedElement.State.State |= int64(mashupsdk.Clicked)
					}
					hfWorldApp.fyneWidgetElements[mashupDetailedElement.Alias].MashupDetailedElement = mashupDetailedElement
					hfWorldApp.fyneWidgetElements[mashupDetailedElement.Alias].OnStatusChanged()
					return
				}
			}
			hfWorldApp.fyneWidgetElements[tabItem.Text].OnStatusChanged()
		}

		torusMenu.SetTabLocation(container.TabLocationTop)

		hfWorldApp.mainWin.SetContent(torusMenu)
		hfWorldApp.mainWin.SetCloseIntercept(func() {
			if hfWorldApp.HeadsupFyneContext.mashupContext != nil {
				hfWorldApp.HeadsupFyneContext.mashupContext.Client.Shutdown(hfWorldApp.HeadsupFyneContext.mashupContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})
			}
			os.Exit(0)
		})
	}
	runtimeHandler := func() {
		go func() {
			w.mainWin.Hide()
		}()

		w.mainWin.ShowAndRun()
	}

	guiboot.InitMainWindow(guiboot.Fyne, initHandler, runtimeHandler)
}

func (w *worldClientInitHandler) RegisterContext(context *mashupsdk.MashupContext) {
	hfWorldApp.HeadsupFyneContext.mashupContext = context
}

// Sets all elements to a "Rest state."
func (w *mashupSdkApiHandler) ResetG3NDetailedElementStates() {
	log.Printf("G3n Received ResetG3NDetailedElementStates\n")
	for _, wes := range hfWorldApp.mashupDetailedElementLibrary {
		wes.SetElementState(mashupsdk.Init)
	}
	log.Printf("G3n finished ResetG3NDetailedElementStates handle.\n")
}

func (mSdk *mashupSdkApiHandler) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	// if hfWorldApp.mainWin != nil && (*hfWorldApp.mainWin).IWindow != nil {
	// 	log.Printf("G3n Received onResize xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	// 	hfWorldApp.displayPositionChan <- displayHint
	// } else {
	// 	if displayHint.Width != 0 && displayHint.Height != 0 {
	// 		log.Printf("G3n initializing with: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	// 		hfWorldApp.displaySetupChan <- displayHint
	// 		hfWorldApp.displayPositionChan <- displayHint
	// 	} else {
	// 		log.Printf("G3n Could not apply xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	// 	}
	// 	log.Printf("G3n finished onResize handle.")
	// }
}

func (hfWorldApp *HFWorldApp) detailMappedFyneComponent(id, description string, de *mashupsdk.MashupDetailedElement) *container.TabItem {
	tabLabel := widget.NewLabel(description)
	tabLabel.Wrapping = fyne.TextWrapWord
	tabItem := container.NewTabItem(id, container.NewBorder(nil, nil, layout.NewSpacer(), nil, container.NewVBox(tabLabel, container.NewAdaptiveGrid(2,
		widget.NewButton("Show", func() {
			// Workaround... mashupdetailedelement points at wrong element sometimes, but shouldn't!
			mashupIndex := hfWorldApp.elementLoaderIndex[hfWorldApp.fyneWidgetElements[de.Alias].GuiComponent.(*container.TabItem).Text]
			hfWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement = hfWorldApp.mashupDetailedElementLibrary[mashupIndex]

			hfWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Hidden, false)
			if hfWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.Genre == "Collection" {
				hfWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Recursive, true)
			}
			hfWorldApp.fyneWidgetElements[de.Alias].OnStatusChanged()
		}), widget.NewButton("Hide", func() {
			// Workaround... mashupdetailedelement points at wrong element sometimes, but shouldn't!
			mashupIndex := hfWorldApp.elementLoaderIndex[hfWorldApp.fyneWidgetElements[de.Alias].GuiComponent.(*container.TabItem).Text]
			hfWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement = hfWorldApp.mashupDetailedElementLibrary[mashupIndex]

			hfWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Hidden, true)
			if hfWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.Genre == "Collection" {
				hfWorldApp.fyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Recursive, true)
			}
			hfWorldApp.fyneWidgetElements[de.Alias].OnStatusChanged()
		})))),
	)
	return tabItem
}

func (mSdk *mashupSdkApiHandler) GetMashupElements() (*mashupsdk.MashupDetailedElementBundle, error) {
	log.Printf("HFWorld Received GetMashupElements\n")
	var concreteElementBundle *mashupsdk.MashupDetailedElementBundle
	DetailedElements := []*mashupsdk.MashupDetailedElement{}

	for _, detailedElement := range hfWorldApp.mashupDetailedElementLibrary {
		DetailedElements = append(DetailedElements, detailedElement)
	}
	concreteElementBundle = &mashupsdk.MashupDetailedElementBundle{
		AuthToken:        client.GetServerAuthToken(),
		DetailedElements: DetailedElements,
	}

	log.Printf("HFWorld GetMashupElements finished\n")
	return concreteElementBundle, nil
}

func (mSdk *mashupSdkApiHandler) UpsertMashupElements(detailedElementBundle *mashupsdk.MashupDetailedElementBundle) (*mashupsdk.MashupDetailedElementBundle, error) {
	log.Printf("HFWorld Received UpsertMashupElements\n")
	// TODO: Implement

	log.Printf("HFWorld UpsertMashupElements updated\n")
	return nil, nil
}
func (mSdk *mashupSdkApiHandler) setStateHelper(g3nId int64, x mashupsdk.DisplayElementState) {

	child := hfWorldApp.mashupDetailedElementLibrary[g3nId]
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
	log.Printf("HFWorld UpsertMashupElementsState called\n")

	ClickedElements := map[int64]*mashupsdk.MashupDetailedElement{}
	recursiveElements := map[int64]*mashupsdk.MashupDetailedElement{}

	for _, es := range elementStateBundle.ElementStates {
		if g3nDetailedElement, ok := hfWorldApp.mashupDetailedElementLibrary[es.GetId()]; ok {
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
		for _, clickedElement := range hfWorldApp.ClickedElements {
			if _, ok := ClickedElements[clickedElement.GetId()]; !ok {
				clickedElement.ApplyState(mashupsdk.Clicked, false)
			}
		}

		hfWorldApp.ClickedElements = hfWorldApp.ClickedElements[:0]

		// Impossible to determine ordering of clicks from upsert at this time.
		for _, g3nDetailedElement := range ClickedElements {
			hfWorldApp.ClickedElements = append(hfWorldApp.ClickedElements, g3nDetailedElement)
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

	log.Printf("HFWorld dispatching focus\n")
	if hfWorldApp.mainWin != nil {
		hfWorldApp.mainWin.Hide()
	}
	log.Printf("HFWorld End UpsertMashupElementsState called\n")
	return &mashupsdk.MashupElementStateBundle{}, nil
}
