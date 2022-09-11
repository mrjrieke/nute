package custosworld

import (
	"log"
	"os"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/mrjrieke/nute/mashupsdk"
	"github.com/mrjrieke/nute/mashupsdk/client"
	"github.com/mrjrieke/nute/mashupsdk/guiboot"
	"github.com/mrjrieke/nute/mashupsdk/server"
)

type mashupSdkApiHandler struct {
}

type worldClientInitHandler struct {
}

type ICustosRenderer interface {
	OnSelected(tabItem *container.TabItem)
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
		AuthToken:     server.GetServerAuthToken(),
		ElementStates: []*mashupsdk.MashupElementState{selectedDetailedElement.State},
	}
	CUWorldApp.HeadsupFyneContext.mashupContext.Client.ResetG3NDetailedElementStates(CUWorldApp.HeadsupFyneContext.mashupContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})

	log.Printf("Display fields set to: %d", selectedDetailedElement.State.State)
	CUWorldApp.HeadsupFyneContext.mashupContext.Client.UpsertMashupElementsState(CUWorldApp.HeadsupFyneContext.mashupContext, &elementStateBundle)
	log.Printf("Finished status change.\n")
}

type ITabItemRenderer interface {
	GetPriority() int64
	BuildTabItem(id int64, concreteElement *mashupsdk.MashupDetailedElement)
	PreRender() // Called at end of all tab item renders.
	RenderTabItem(concreteElement *mashupsdk.MashupDetailedElement)
	Refresh() // Called at end of all tab item renders.
}

type CustosWorldApp struct {
	Headless                     bool // Mode for troubleshooting.
	mashupSdkApiHandler          *mashupSdkApiHandler
	Title                        string
	Icon                         *fyne.StaticResource
	MainWindowSize               fyne.Size
	wClientInitHandler           *worldClientInitHandler
	HeadsupFyneContext           *CustosContext
	MainWin                      fyne.Window
	mashupDisplayContext         *mashupsdk.MashupDisplayContext
	DetailedElements             []*mashupsdk.MashupDetailedElement
	MashupDetailedElementLibrary map[int64]*mashupsdk.MashupDetailedElement
	ElementLoaderIndex           map[string]int64 // mashup indexes by Name
	FyneWidgetElements           map[string]*FyneWidgetBundle
	TabItemMenu                  *container.AppTabs
	CustomTabItems               map[string]func(custosWorlApp *CustosWorldApp, id string) *container.TabItem
	CustomTabItemRenderer        map[string]ITabItemRenderer
	CustosRenderer               ICustosRenderer
}

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

func NewCustosWorldApp(headless bool,
	detailedElements []*mashupsdk.MashupDetailedElement,
	renderer ICustosRenderer) *CustosWorldApp {
	CUWorldApp = &CustosWorldApp{
		Headless:                     headless,
		mashupSdkApiHandler:          &mashupSdkApiHandler{},
		HeadsupFyneContext:           &CustosContext{},
		DetailedElements:             detailedElements,
		MainWin:                      nil,
		mashupDisplayContext:         &mashupsdk.MashupDisplayContext{MainWinDisplay: &mashupsdk.MashupDisplayHint{}},
		MashupDetailedElementLibrary: map[int64]*mashupsdk.MashupDetailedElement{}, // mashupDetailedElementLibrary,
		ElementLoaderIndex:           map[string]int64{},                           // elementLoaderIndex
		FyneWidgetElements:           map[string]*FyneWidgetBundle{},
		CustomTabItems:               map[string]func(custosWorlApp *CustosWorldApp, id string) *container.TabItem{},
		CustomTabItemRenderer:        map[string]ITabItemRenderer{},
		CustosRenderer:               renderer,
	}

	return CUWorldApp
}

type InitEvent struct {
}

func (w *CustosWorldApp) ResetChangeStates() []*mashupsdk.MashupElementState {
	changedElements := []*mashupsdk.MashupElementState{}
	for _, g3nDetailedElement := range w.MashupDetailedElementLibrary {
		if mashupsdk.DisplayElementState(g3nDetailedElement.GetMashupElementState().State) != mashupsdk.Init {
			g3nDetailedElement.ApplyState(mashupsdk.Clicked, false)
			changedElements = append(changedElements, g3nDetailedElement.GetMashupElementState())
		}
	}

	return changedElements
}

func (w *CustosWorldApp) InitMainWindow() {
	log.Printf("Initializing MainWin.")

	initHandler := func(a fyne.App) {
		log.Printf("InitHandler.")
		CUWorldApp.MainWin = a.NewWindow(w.Title)
		CUWorldApp.MainWin.SetIcon(w.Icon)
		CUWorldApp.MainWin.Resize(CUWorldApp.MainWindowSize)
		CUWorldApp.MainWin.SetFixedSize(false)

		CUWorldApp.TabItemMenu = container.NewAppTabs()

		if CUWorldApp.CustosRenderer != nil {
			CUWorldApp.TabItemMenu.OnSelected = CUWorldApp.CustosRenderer.OnSelected
		}

		CUWorldApp.TabItemMenu.SetTabLocation(container.TabLocationTop)

		CUWorldApp.MainWin.SetContent(CUWorldApp.TabItemMenu)
		CUWorldApp.MainWin.SetCloseIntercept(func() {
			if CUWorldApp.HeadsupFyneContext.mashupContext != nil {
				CUWorldApp.HeadsupFyneContext.mashupContext.Client.Shutdown(CUWorldApp.HeadsupFyneContext.mashupContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})
			}
			os.Exit(0)
		})
	}
	runtimeHandler := func() {
		if w.MainWin != nil {
			log.Printf("CustosWorld main win initialized\n")
			if CUWorldApp.mashupDisplayContext != nil &&
				(CUWorldApp.Headless ||
					(CUWorldApp.mashupDisplayContext.GetSettled()&mashupsdk.AppInitted) == mashupsdk.AppInitted) {
				log.Printf("CustosWorld app settled... starting up.\n")
				w.MainWin.ShowAndRun()
			} else {
				if !CUWorldApp.Headless {
					w.MainWin.ShowAndRun()
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
	for _, wes := range CUWorldApp.MashupDetailedElementLibrary {
		wes.SetElementState(mashupsdk.Init)
	}
	log.Printf("CustosWorld finished ResetG3NDetailedElementStates handle.\n")
}

func (mSdk *mashupSdkApiHandler) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	log.Printf("CustosWorld OnResize - not implemented yet..\n")
	if CUWorldApp.MainWin != nil && CUWorldApp.mashupDisplayContext != nil && CUWorldApp.mashupDisplayContext.MainWinDisplay != nil {
		CUWorldApp.mashupDisplayContext.MainWinDisplay.Focused = displayHint.Focused
		// TODO: Resize without infinite looping....
		// The moment fyne is resized, it'll want to resize g3n...
		// Which then wants to resize fyne ad-infinitum
		//CUWorldApp.MainWin.PosResize(int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height))
		log.Printf("CustosWorld Received onResize xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	} else {
		log.Printf("CustosWorld Could not apply xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	}
}

func (custosWorldApp *CustosWorldApp) DetailMappedFyneComponent(id, description string, genre string, renderer string, tabItemFunc func(custosWorlApp *CustosWorldApp, id string) *container.TabItem) {
	de := &mashupsdk.MashupDetailedElement{Name: id, Description: description, Genre: genre, Renderer: renderer}
	custosWorldApp.FyneWidgetElements[id] = &FyneWidgetBundle{
		GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
			GuiComponent:          nil,
			MashupDetailedElement: de,
		},
	}
	custosWorldApp.CustomTabItems[id] = tabItemFunc
}

func (custosWorldApp *CustosWorldApp) DetailFyneComponent(de *mashupsdk.MashupDetailedElement, tabItemFunc func(custosWorlApp *CustosWorldApp, id string) *container.TabItem) {
	log.Printf("CustosWorldApp.DetailFyneComponent building on: %s name: %s\n", de.Alias, de.Name)
	custosWorldApp.FyneWidgetElements[de.Name] = &FyneWidgetBundle{
		GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
			GuiComponent:          nil,
			MashupDetailedElement: de,
		},
	}
	custosWorldApp.CustomTabItems[de.Name] = tabItemFunc
}

func (CustosWorldApp *CustosWorldApp) TorusParser(childId int64) {
	child := CUWorldApp.MashupDetailedElementLibrary[childId]
	if child != nil && child.Alias != "" {
		log.Printf("TorusParser lookup on: %s\n", child.Alias)
		if CUWorldApp.FyneWidgetElements != nil && CUWorldApp.FyneWidgetElements[child.Alias].MashupDetailedElement != nil && CUWorldApp.FyneWidgetElements[child.Alias].GuiComponent != nil {
			CUWorldApp.FyneWidgetElements[child.Alias].MashupDetailedElement.Copy(child)
			CUWorldApp.FyneWidgetElements[child.Alias].GuiComponent.(*container.TabItem).Text = child.Name
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
		CUWorldApp.MashupDetailedElementLibrary[concreteElement.Id] = concreteElement
		CUWorldApp.ElementLoaderIndex[concreteElement.Name] = concreteElement.Id
	}
	log.Printf("CustosWorld parsing tori.\n")
	for _, concreteElement := range detailedElementBundle.DetailedElements {
		if tabItemRenderer, tabItemRendererOk := CUWorldApp.CustomTabItemRenderer[concreteElement.GetCustosrenderer()]; tabItemRendererOk {
			tabItemRenderer.BuildTabItem(concreteElement.Id, concreteElement)
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

	child := CUWorldApp.MashupDetailedElementLibrary[g3nId]
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

	recursiveElements := map[int64]*mashupsdk.MashupDetailedElement{}

	// Separate clicked from declicked.
	for _, es := range elementStateBundle.ElementStates {
		if g3nDetailedElement, ok := CUWorldApp.MashupDetailedElementLibrary[es.GetId()]; ok {
			g3nDetailedElement.SetElementState(mashupsdk.DisplayElementState(es.State))
			if g3nDetailedElement.IsStateSet(mashupsdk.Recursive) {
				recursiveElements[es.GetId()] = g3nDetailedElement
			}

			log.Printf("Display fields set to: %d", g3nDetailedElement.GetMashupElementState())
			if (mashupsdk.DisplayElementState(es.State) & mashupsdk.Clicked) == mashupsdk.Clicked {
				CUWorldApp.MashupDetailedElementLibrary[g3nDetailedElement.Id].ApplyState(mashupsdk.Clicked, true)
			} else {
				CUWorldApp.MashupDetailedElementLibrary[g3nDetailedElement.Id].ApplyState(mashupsdk.Clicked, false)
			}
		}
	}

	if len(recursiveElements) > 0 {
		log.Printf("CustosWorld UpsertMashupElementsState apply recursive elements\n")

		for _, recursiveElement := range recursiveElements {
			stateBits := recursiveElement.State.State
			// Unset recursive for child elements
			stateBits &= ^int64(mashupsdk.Recursive)
			// Apply this state change to all child elements.
			mSdk.setStateHelper(recursiveElement.GetId(), mashupsdk.DisplayElementState(stateBits))
		}
	}

	// Wipe anything there out.
	CUWorldApp.TabItemMenu.SetItems([]*container.TabItem{})

	orderedRenderingMap := map[int64][]*mashupsdk.MashupDetailedElement{}

	// Get ready for new render cycle.
	for _, tabItemRenderer := range CUWorldApp.CustomTabItemRenderer {
		tabItemRenderer.PreRender()
	}

	// Impossible to determine ordering of clicks from upsert at this time.
	for _, concreteElement := range CUWorldApp.MashupDetailedElementLibrary {
		// Set all clicked elements...
		if tabItemRenderer, tabItemRendererOk := CUWorldApp.CustomTabItemRenderer[concreteElement.Custosrenderer]; tabItemRendererOk {
			orderedRenderingMap[tabItemRenderer.GetPriority()] = append(orderedRenderingMap[tabItemRenderer.GetPriority()], concreteElement)
		}
	}

	keys := make([]int64, 0, len(orderedRenderingMap))
	for k := range orderedRenderingMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		var tir ITabItemRenderer
		for _, concreteElement := range orderedRenderingMap[k] {
			if tabItemRenderer, tabItemRendererOk := CUWorldApp.CustomTabItemRenderer[concreteElement.Custosrenderer]; tabItemRendererOk {
				tabItemRenderer.RenderTabItem(concreteElement)
				tir = tabItemRenderer
			}
		}
		tir.Refresh()
	}

	log.Printf("CustosWorld dispatching focus\n")
	// if CUWorldApp.MainWin != nil {
	// 	CUWorldApp.MainWin.Hide()
	// }
	log.Printf("CustosWorld End UpsertMashupElementsState called\n")
	return &mashupsdk.MashupElementStateBundle{}, nil
}
