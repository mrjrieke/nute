package hfworld

import (
	"log"

	"fyne.io/fyne/v2"
	"github.com/g3n/engine/app"
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

var hfWorld HFWorldApp

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

	hfWorld = HFWorldApp{
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
	return &hfWorld
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

	initHandler := func(a *app.Application) {
		log.Printf("InitHandler.")
	}
	runtimeHandler := func() {
		w.mainWin.ShowAndRun()
	}

	guiboot.InitMainWindow(guiboot.G3n, initHandler, runtimeHandler)
}

func (w *worldClientInitHandler) RegisterContext(context *mashupsdk.MashupContext) {
	hfWorld.HeadsupFyneContext.mashupContext = context
}

// Sets all elements to a "Rest state."
func (w *mashupSdkApiHandler) ResetG3NDetailedElementStates() {
	log.Printf("G3n Received ResetG3NDetailedElementStates\n")
	for _, wes := range hfWorld.mashupDetailedElementLibrary {
		wes.SetElementState(mashupsdk.Init)
	}
	log.Printf("G3n finished ResetG3NDetailedElementStates handle.\n")
}

func (mSdk *mashupSdkApiHandler) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	// if hfWorld.mainWin != nil && (*hfWorld.mainWin).IWindow != nil {
	// 	log.Printf("G3n Received onResize xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	// 	hfWorld.displayPositionChan <- displayHint
	// } else {
	// 	if displayHint.Width != 0 && displayHint.Height != 0 {
	// 		log.Printf("G3n initializing with: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	// 		hfWorld.displaySetupChan <- displayHint
	// 		hfWorld.displayPositionChan <- displayHint
	// 	} else {
	// 		log.Printf("G3n Could not apply xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	// 	}
	// 	log.Printf("G3n finished onResize handle.")
	// }
}

func (mSdk *mashupSdkApiHandler) GetMashupElements() (*mashupsdk.MashupDetailedElementBundle, error) {
	log.Printf("HFWorld Received GetMashupElements\n")
	var concreteElementBundle *mashupsdk.MashupDetailedElementBundle
	DetailedElements := []*mashupsdk.MashupDetailedElement{}

	for _, detailedElement := range hfWorld.mashupDetailedElementLibrary {
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

	child := hfWorld.mashupDetailedElementLibrary[g3nId]
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
		if g3nDetailedElement, ok := hfWorld.mashupDetailedElementLibrary[es.GetId()]; ok {
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
		for _, clickedElement := range hfWorld.ClickedElements {
			if _, ok := ClickedElements[clickedElement.GetId()]; !ok {
				clickedElement.ApplyState(mashupsdk.Clicked, false)
			}
		}

		hfWorld.ClickedElements = hfWorld.ClickedElements[:0]

		// Impossible to determine ordering of clicks from upsert at this time.
		for _, g3nDetailedElement := range ClickedElements {
			hfWorld.ClickedElements = append(hfWorld.ClickedElements, g3nDetailedElement)
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
	if hfWorld.mainWin != nil {
		// TODO: Can we get rid of this?
		//		hfWorld.mainWin.Dispatch(gui.OnFocus, nil)
	}
	log.Printf("HFWorld End UpsertMashupElementsState called\n")
	return &mashupsdk.MashupElementStateBundle{}, nil
}
