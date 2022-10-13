package g3nworld

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/mrjrieke/nute/g3nd/g3nmash"
	g3ndpalette "github.com/mrjrieke/nute/g3nd/palette"
	"github.com/mrjrieke/nute/g3nd/worldg3n/g3ndisplay"
	"github.com/mrjrieke/nute/mashupsdk"
	"github.com/mrjrieke/nute/mashupsdk/client"
	"github.com/mrjrieke/nute/mashupsdk/guiboot"
	"github.com/mrjrieke/nute/mashupsdk/server"
)

type mashupSdkApiHandler struct {
}

type worldClientInitHandler struct {
}

// TODO: Can I get rid of this definition?
// -- duplicated in genericrenderer.go.
type IG3nRenderer interface {
	Layout(worldApp *WorldApp, g3nRenderableElements []*g3nmash.G3nDetailedElement)
	InitRenderLoop(worldApp *WorldApp) bool
	RenderElement(worldApp *WorldApp, g3n *g3nmash.G3nDetailedElement) bool
}

type WorldApp struct {
	custos              bool // Run in guardian mode
	headless            bool // Mode for troubleshooting.
	MSdkApiHandler      *mashupSdkApiHandler
	wClientInitHandler  *worldClientInitHandler
	displaySetupChan    chan *mashupsdk.MashupDisplayHint
	displayPositionChan chan *mashupsdk.MashupDisplayHint
	mainWin             *app.Application
	currentTargetFPS    uint             // Current target frame rate
	frameRater          *util.FrameRater // Render loop frame rater
	scene               *core.Node
	cam                 *camera.Camera
	oc                  *camera.OrbitControl
	IG3nRenderer        IG3nRenderer
	IG3nDisplayRenderer g3ndisplay.IG3nDisplayRenderer

	MashupContext *mashupsdk.MashupContext // Needed for callbacks to other mashups

	// Library for mashup objects
	elementLibraryDictionary map[int64]*g3nmash.G3nDetailedElement

	maxElementId       int64
	RootElements       []*g3nmash.G3nDetailedElement
	ConcreteElements   map[int64]*g3nmash.G3nDetailedElement // g3n indexes by string...
	elementLoaderIndex map[string]int64                      // g3n indexes by loader id...
	ClickedElements    []*g3nmash.G3nDetailedElement         // g3n indexes by string...
	backgroundG3n      *g3nmash.G3nDetailedElement
	Sticky             bool

	Focused bool // Whether current window has focus.

	isInit bool
}

var worldApp WorldApp

func NewWorldApp(headless bool, custos bool, renderer IG3nRenderer, displayRenderer g3ndisplay.IG3nDisplayRenderer) *WorldApp {
	worldApp = WorldApp{
		custos:                   custos,
		headless:                 headless,
		MSdkApiHandler:           &mashupSdkApiHandler{},
		elementLibraryDictionary: map[int64]*g3nmash.G3nDetailedElement{},
		ConcreteElements:         map[int64]*g3nmash.G3nDetailedElement{},
		elementLoaderIndex:       map[string]int64{},
		ClickedElements:          []*g3nmash.G3nDetailedElement{},
		displaySetupChan:         make(chan *mashupsdk.MashupDisplayHint, 1),
		displayPositionChan:      make(chan *mashupsdk.MashupDisplayHint, 1),
		IG3nRenderer:             renderer,
		IG3nDisplayRenderer:      displayRenderer,
	}
	return &worldApp
}

type InitEvent struct {
}

func (w *WorldApp) SetFrameRate(targetFPS uint) {
	if w.currentTargetFPS == targetFPS {
		return
	} else {
		fmt.Printf("Setting framerate to: %d\n", targetFPS)
		w.currentTargetFPS = targetFPS
		if targetFPS > 0 {
			w.frameRater = util.NewFrameRater(targetFPS)
		} else {
			w.frameRater = util.NewFrameRater(1)
		}
	}
}

func (w *WorldApp) GetAuthToken() string {
	if w.custos {
		return client.GetServerAuthToken()
	} else {
		return server.GetServerAuthToken()
	}
}

func (w *WorldApp) G3nOnFocus(name string, ev interface{}) {
	log.Printf("G3nWorld Focus gained\n")

	if _, iOk := ev.(InitEvent); iOk {
		var postPonedCollections []*g3nmash.G3nDetailedElement
		g3nCollections, err := w.GetG3nDetailedGenreFilteredElements("Collection")
		if err != nil {
			log.Fatalf(err.Error(), err)
		}
		if len(g3nCollections) == 0 {
			log.Fatalf("No elements to render.  If running standalone, provide -headless flag.")
		}
		for _, g3nCollection := range g3nCollections {
			if g3nCollection.GetDetailedElement().Colabrenderer != "" {
				postPonedCollections = append(postPonedCollections, g3nCollection)
			} else {
				var g3nCollectionElements []*g3nmash.G3nDetailedElement
				var g3nCollectionErr error

				if g3nCollection.GetDetailedElement().GetRenderer() != "" {
					g3nCollectionElements, g3nCollectionErr = w.GetG3nDetailedFilteredElements(g3nCollection.GetDetailedElement().GetRenderer(), true)
					if len(g3nCollectionElements) == 0 {
						// Try lookup by child elements instead...
						g3nCollectionElements, g3nCollectionErr = w.GetG3nDetailedChildElements(g3nCollection)
					}
				} else {
					g3nCollectionElements, g3nCollectionErr = w.GetG3nDetailedChildElements(g3nCollection)
				}
				if g3nCollectionErr != nil {
					log.Fatalf(err.Error(), g3nCollectionErr)
				}
				// Handoff...
				w.IG3nRenderer.Layout(w, g3nCollectionElements)
			}
		}

		for _, g3nCollection := range postPonedCollections {
			g3nCollectionElements, err := w.GetG3nDetailedChildElements(g3nCollection)
			if err != nil {
				log.Fatalf(err.Error(), err)
			}
			// Handoff...
			w.IG3nRenderer.Layout(w, g3nCollectionElements)
		}
	} else {

		// Focus gained...
		log.Printf("G3n Focus gained\n")

		w.Transform()
		log.Printf("G3n End Focus gained\n")
	}

	log.Printf("G3nWorld End Focus gained\n")
}

func (w *WorldApp) ResetChangeStates() []*mashupsdk.MashupElementState {
	changedElements := []*mashupsdk.MashupElementState{}
	for _, g3nDetailedElement := range w.ConcreteElements {
		if g3nDetailedElement.GetDisplayState() != mashupsdk.Init {
			g3nDetailedElement.ApplyState(mashupsdk.Clicked, false)
			changedElements = append(changedElements, g3nDetailedElement.GetMashupElementState())
		}
	}

	return changedElements
}

func (w *WorldApp) NewElementIdPump() int64 {
	w.maxElementId = w.maxElementId + 1
	return w.maxElementId
}

func (w *WorldApp) CloneG3nDetailedElement(g3nElement *g3nmash.G3nDetailedElement, elementStates *[]interface{}) *g3nmash.G3nDetailedElement {
	return w.indexG3nDetailedElement(g3nmash.CloneG3nDetailedElement(w.GetG3nDetailedElementById, w.GetG3nDetailedLibraryElementById, w.indexG3nDetailedElement, w.NewElementIdPump, g3nElement, elementStates))
}

func (w *WorldApp) NewG3nDetailedElement(detailedElement *mashupsdk.MashupDetailedElement, deepCopy bool) *g3nmash.G3nDetailedElement {
	return w.indexG3nDetailedElement(g3nmash.NewG3nDetailedElement(detailedElement, deepCopy))
}

func (w *WorldApp) indexG3nDetailedElement(g3nDetailedElement *g3nmash.G3nDetailedElement) *g3nmash.G3nDetailedElement {
	if g3nDetailedElement.GetBasisId() < 0 && g3nDetailedElement.GetDisplayId() == 0 {
		w.elementLibraryDictionary[g3nDetailedElement.GetBasisId()] = g3nDetailedElement
	} else {
		w.ConcreteElements[g3nDetailedElement.GetDisplayId()] = g3nDetailedElement
		w.elementLoaderIndex[g3nDetailedElement.GetDisplayName()] = g3nDetailedElement.GetDisplayId()
		if g3nDetailedElement.IsBackground() {
			w.backgroundG3n = g3nDetailedElement
		}
	}
	return g3nDetailedElement
}

func (w *WorldApp) GetG3nDetailedFilteredElements(renderer string, abstract bool) ([]*g3nmash.G3nDetailedElement, error) {
	filteredElements := []*g3nmash.G3nDetailedElement{}
	if renderer == "" {
		log.Printf("No filter provided.  No filtered elements found.\n")
		return nil, errors.New("no filter provided - no filtered elements found")
	}
	for _, element := range w.ConcreteElements {
		if element.GetDetailedElement().Renderer == renderer {
			if abstract {
				if element.IsAbstract() {
					filteredElements = append(filteredElements, element)
				}
			} else {
				filteredElements = append(filteredElements, element)
			}
		}
	}

	return filteredElements, nil
}

func (w *WorldApp) GetG3nDetailedChildElements(g3n *g3nmash.G3nDetailedElement) ([]*g3nmash.G3nDetailedElement, error) {
	filteredElements := []*g3nmash.G3nDetailedElement{}
	if g3n == nil {
		log.Printf("No filter provided.  No filtered elements found.\n")
		return nil, errors.New("no filter provided - no filtered elements found")
	}
	for _, childId := range g3n.GetChildElementIds() {
		if tc, tErr := worldApp.GetG3nDetailedElementById(childId); tErr == nil {
			filteredElements = append(filteredElements, tc)
		} else {
			log.Printf("Skipping non-concrete abstract element: %d\n", childId)
			continue
		}
	}

	return filteredElements, nil
}

func (w *WorldApp) GetG3nDetailedGenreFilteredElements(genre string) ([]*g3nmash.G3nDetailedElement, error) {
	filteredElements := w.RootElements
	if len(filteredElements) > 0 {
		return filteredElements, nil
	}
	for _, element := range w.ConcreteElements {
		if element.GetDetailedElement().GetGenre() == genre {
			filteredElements = append(filteredElements, element)
		}
	}
	w.RootElements = filteredElements

	return filteredElements, nil
}

func (w *WorldApp) AddToScene(node core.INode) *core.Node {
	if node == nil {
		return nil
	}
	if w.scene.FindLoaderID(node.GetNode().LoaderID()) == nil {
		log.Printf("Item added %s: %v", node.GetNode().LoaderID(), w.scene.Add(node))
		return w.scene.Add(node)
	} else {
		return nil
	}
}

func (w *WorldApp) UpsertToScene(node core.INode) *core.Node {
	if node == nil {
		return nil
	}
	if w.scene.FindLoaderID(node.GetNode().LoaderID()) == nil {
		return w.scene.Add(node)
	}
	return nil
}

func (w *WorldApp) RemoveFromScene(node core.INode) bool {
	return w.scene.Remove(node)
}

func (w *WorldApp) GetG3nDetailedElementById(eid int64) (*g3nmash.G3nDetailedElement, error) {
	if g3nElement, g3nElementOk := w.ConcreteElements[eid]; g3nElementOk {
		return g3nElement, nil
	}
	return nil, fmt.Errorf("element does not exist: %d", eid)
}

func (w *WorldApp) GetG3nDetailedChildElementsByGenre(g3n *g3nmash.G3nDetailedElement, genre string) []*g3nmash.G3nDetailedElement {
	results := []*g3nmash.G3nDetailedElement{}
	for _, childId := range g3n.GetChildElementIds() {
		if g3nChild, err := w.GetG3nDetailedElementById(childId); err == nil {
			if g3nChild.HasGenre(genre) {
				results = append(results, g3nChild)
			}
		}
	}
	return results
}

func (w *WorldApp) GetG3nDetailedLibraryElementById(eid int64) (*g3nmash.G3nDetailedElement, error) {
	if g3nElement, g3nElementOk := w.elementLibraryDictionary[eid]; g3nElementOk {
		return g3nElement, nil
	}
	return nil, fmt.Errorf("element does not exist: %d", eid)
}

func (w *WorldApp) GetParentElements(g3nDetailedElement *g3nmash.G3nDetailedElement) []*g3nmash.G3nDetailedElement {
	parentIds := g3nDetailedElement.GetParentElementIds()
	g3nParentDetailedElements := []*g3nmash.G3nDetailedElement{}
	for _, parentId := range parentIds {
		if g3parent, gpErr := w.GetG3nDetailedElementById(parentId); gpErr == nil {
			g3nParentDetailedElements = append(g3nParentDetailedElements, g3parent)
		}
	}

	return g3nParentDetailedElements
}

func (w *WorldApp) GetSiblingElements(g3nDetailedElement *g3nmash.G3nDetailedElement) []*g3nmash.G3nDetailedElement {
	parentIds := g3nDetailedElement.GetParentElementIds()
	g3nParentDetailedElements := []*g3nmash.G3nDetailedElement{}
	for _, parentId := range parentIds {
		if g3parent, gpErr := w.GetG3nDetailedElementById(parentId); gpErr == nil {
			g3nParentDetailedElements = append(g3nParentDetailedElements, g3parent)
		}
	}
	g3nSiblingDetailedElements := []*g3nmash.G3nDetailedElement{}

	for _, g3nParentDetailedElement := range g3nParentDetailedElements {
		for _, childId := range g3nParentDetailedElement.GetChildElementIds() {
			if g3nSibling, gsErr := w.GetG3nDetailedElementById(childId); gsErr == nil {
				g3nSiblingDetailedElements = append(g3nSiblingDetailedElements, g3nSibling)
			}
		}
	}
	return g3nSiblingDetailedElements
}

func (w *WorldApp) Cast(inode core.INode, caster *collision.Raycaster) (core.INode, []collision.Intersect) {
	// Ignore invisible nodes and their descendants
	if !inode.Visible() {
		return nil, nil
	}

	if _, ok := inode.(gui.IPanel); ok {
		// TODO: Do we care about these types at all?
	} else if igr, ok := inode.(graphic.IGraphic); ok {
		if igr.Renderable() {
			if _, meshOk := inode.(*graphic.Mesh); meshOk {
				return inode, caster.IntersectObject(inode, false)
			}
		}
		// Ignore everything else.
	}

	if inode.Children() != nil {
		for _, ichild := range inode.Children() {
			if n, intersections := w.Cast(ichild, caster); n != nil && len(intersections) > 0 {
				return n, intersections
			}
		}
	}
	return nil, nil
}

func (w *WorldApp) InitServer(callerCreds string, insecure bool, maxMessageLength int) {
	if callerCreds != "" {
		server.InitServer(callerCreds, insecure, maxMessageLength, w.MSdkApiHandler, w.wClientInitHandler)
	} else {
		go func() {
			w.displaySetupChan <- &mashupsdk.MashupDisplayHint{Xpos: 0, Ypos: 0, Width: 400, Height: 800}
		}()
	}
}

func (w *WorldApp) Transform() []*mashupsdk.MashupElementState {
	changedElements := []*mashupsdk.MashupElementState{}
	attitudeVisitedNodes := map[int64]bool{}

	// Notify renderers about to enter render loop.
	worldApp.IG3nRenderer.InitRenderLoop(w)

	for _, g3nDetailedElement := range w.ConcreteElements {
		if g3nDetailedElement.IsAbstract() {
			continue
		}
		changed := worldApp.IG3nRenderer.RenderElement(w, g3nDetailedElement)
		if !g3nDetailedElement.IsStateSet(mashupsdk.Hidden) && !g3nDetailedElement.IsBackground() {
			if g3nDetailedElement.IsStateSet(mashupsdk.Clicked) {
				if g3nDetailedElement.HasAttitudeAdjustment() {
					log.Printf("G3n Has parents\n")
					parentIds := g3nDetailedElement.GetParentElementIds()
					g3nParentDetailedElements := []*g3nmash.G3nDetailedElement{}
					for _, parentId := range parentIds {
						if g3parent, gpErr := w.GetG3nDetailedElementById(parentId); gpErr == nil {
							g3nParentDetailedElements = append(g3nParentDetailedElements, g3parent)
						}
						attitudeVisitedNodes[parentId] = true
					}
					log.Printf("G3n adjusting for parents: %d\n", len(g3nParentDetailedElements))

					g3nDetailedElement.AdjustAttitude(g3nParentDetailedElements)
				} else {
					if _, vOk := attitudeVisitedNodes[g3nDetailedElement.GetDisplayId()]; !vOk {
						g3nDetailedElement.AdjustAttitude([]*g3nmash.G3nDetailedElement{g3nDetailedElement})
						attitudeVisitedNodes[g3nDetailedElement.GetDisplayId()] = true
					}
				}
			} else {
				if _, vOk := attitudeVisitedNodes[g3nDetailedElement.GetDisplayId()]; !vOk {
					g3nDetailedElement.AdjustAttitude([]*g3nmash.G3nDetailedElement{g3nDetailedElement})
					attitudeVisitedNodes[g3nDetailedElement.GetDisplayId()] = true
				}
			}
		}

		if changed {
			changedElements = append(changedElements, g3nDetailedElement.GetMashupElementState())
		}
	}
	return changedElements
}

func (w *WorldApp) InitMainWindow() {
	log.Printf("Initializing mainWin.")

	initHandler := func(a *app.Application) {
		log.Printf("InitHandler.")

		if w.mainWin == nil {
			log.Printf("Main app handle initialized.")
			w.mainWin = a
		}
		log.Printf("Frame rater setup.")
		w.SetFrameRate(5)
		log.Printf("Frame rater setup complete.")

		displayHint := <-w.displaySetupChan
		log.Printf("Initializing app.")
		app.AppCustom(a, "Hello world G3n", int(displayHint.Width), int(displayHint.Height), int(displayHint.Xpos), int(displayHint.Ypos+displayHint.Height))
		log.Printf("Initializing scene.")
		w.scene = core.NewNode()

		// Set the scene to be managed by the gui manager
		gui.Manager().Set(w.scene)

		// Create perspective camera
		w.cam = camera.New(1)
		w.cam.SetPosition(0, 0, 3)
		w.scene.Add(w.cam)

		// Set up orbit control for the camera
		w.oc = camera.NewOrbitControl(w.cam)
		w.oc.Zoom(60.0)
		log.Printf("Finished Orbit Control setup.")

		// Set up callback to update viewport and camera aspect ratio when the window is resized
		onResize := func(evname string, ev interface{}) {
			// Get framebuffer size and update viewport accordingly
			width, height := a.GetSize()
			a.Gls().Viewport(0, 0, int32(width), int32(height))
			// Update the camera's aspect ratio
			w.cam.SetAspect(float32(width) / float32(height))

			xpos, ypos := (*w.mainWin).IWindow.(*window.GlfwWindow).Window.GetPos()

			if w.MashupContext != nil {
				w.MashupContext.Client.OnResize(w.MashupContext,
					&mashupsdk.MashupDisplayBundle{
						AuthToken: worldApp.GetAuthToken(),
						MashupDisplayHint: &mashupsdk.MashupDisplayHint{
							Xpos:   int64(xpos),
							Ypos:   int64(ypos),
							Width:  int64(width),
							Height: int64(height),
						},
					})
			}
		}
		a.Subscribe(window.OnWindowSize, onResize)
		onResize("", nil)

		w.mainWin.Subscribe(gui.OnFocus, w.G3nOnFocus)
		if iWindow, iWindowOk := (*w.mainWin).IWindow.(*window.GlfwWindow); iWindowOk {

			// TODO: OSX enable?
			//iWindow.Window.SetAttrib(glfw.OpenGLForwardCompatible, glfw.True)
			//iWindow.Window.SetAttrib(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
			//iWindow.Window.SetAttrib(glfw.CocoaGraphicsSwitching, 1)
			//iWindow.Window.SetAttrib(glfw.Samples, 0)
			//iWindow.Window.SetAttrib(glfw.Visible, glfw.False)
			//iWindow.Window.SetAttrib(glfw.RefreshRate, glfw.DontCare)

			iWindow.Window.SetAttrib(glfw.Floating, 1)
			iWindow.Window.SetAttrib(glfw.FocusOnShow, 1)

			iWindow.Window.SetCloseCallback(func(glfwWindow *glfw.Window) {
				if w.MashupContext != nil {
					w.MashupContext.Client.Shutdown(w.MashupContext, &mashupsdk.MashupEmpty{AuthToken: worldApp.GetAuthToken()})
				}
				os.Exit(0)
			})
		}

		w.mainWin.Subscribe(gui.OnKeyDown, func(name string, ev interface{}) {
			kev := ev.(*window.KeyEvent)
			if kev.Key == window.KeyLeftControl {
				w.Sticky = true
			}
		})
		w.mainWin.Subscribe(gui.OnKeyUp, func(name string, ev interface{}) {
			kev := ev.(*window.KeyEvent)
			if kev.Key == window.KeyLeftControl {
				w.Sticky = false
			}
		})
		w.mainWin.Subscribe(window.OnWindowFocus, func(name string, ev interface{}) {
			wev := ev.(*window.FocusEvent)
			w.Focused = wev.Focused
			// Tell fyne not to try to regain focus.
			if w.MashupContext != nil && wev.Focused {
				w.MashupContext.Client.OnResize(w.MashupContext,
					&mashupsdk.MashupDisplayBundle{
						AuthToken: worldApp.GetAuthToken(),
						MashupDisplayHint: &mashupsdk.MashupDisplayHint{
							Focused: wev.Focused,
						},
					})
			}
		})
		w.mainWin.Subscribe(gui.OnMouseDown, func(name string, ev interface{}) {
			w.SetFrameRate(30)
		})
		w.mainWin.Subscribe(gui.OnMouseUp, func(name string, ev interface{}) {
			w.SetFrameRate(0)
			mev := ev.(*window.MouseEvent)
			if mev.Mods == window.ModControl {
				w.Sticky = true
			} else {
				w.Sticky = false
			}

			g3Width, g3Height := w.mainWin.GetSize()

			xPosNdc := 2*(mev.Xpos/float32(g3Width)) - 1
			yPosNdc := -2*(mev.Ypos/float32(g3Height)) + 1
			caster := collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
			caster.SetFromCamera(w.cam, xPosNdc, yPosNdc)

			if w.scene.Visible() {
				itemClicked, _ := w.Cast(w.scene, caster)

				itemMatched := false
				if itemClicked != nil {
					if g3nDetailedIndex, ok := w.elementLoaderIndex[itemClicked.GetNode().LoaderID()]; ok {
						if g3nDetailedElement, ok := w.ConcreteElements[g3nDetailedIndex]; ok {
							g3nDetailedElement.ApplyState(mashupsdk.Clicked, true)
							fmt.Printf("matched: %s\n", g3nDetailedElement.GetDisplayName())
							itemMatched = true
							for _, clickedElement := range w.ClickedElements {
								if clickedElement.GetDisplayId() != g3nDetailedElement.GetDisplayId() {
									clickedElement.ApplyState(mashupsdk.Clicked, false)
								}
							}
							if !w.Sticky {
								w.ClickedElements = w.ClickedElements[:0]
							}
							w.ClickedElements = append(w.ClickedElements, g3nDetailedElement)
						}
					}
				}

				if !itemMatched {
					w.backgroundG3n.ApplyState(mashupsdk.Clicked, true)
					for _, clickedElement := range w.ClickedElements {
						if clickedElement.GetDisplayId() != w.backgroundG3n.GetDisplayId() {
							clickedElement.ApplyState(mashupsdk.Clicked, false)
						}
					}
					if !w.Sticky {
						w.ClickedElements = w.ClickedElements[:0]
					}
					w.ClickedElements = append(w.ClickedElements, w.backgroundG3n)
				} else {
					w.backgroundG3n.ApplyState(mashupsdk.Clicked, false)
				}
				changedElements := w.Transform()
				if !itemMatched {
					changedElements = append(changedElements, w.backgroundG3n.GetMashupElementState())
				}

				if !w.headless {
					if mev.Button == window.MouseButtonRight {
						// Mark all changed elements as unclicked for
						// mashup sake.
						for _, changedElement := range changedElements {
							if (mashupsdk.Clicked & mashupsdk.DisplayElementState(changedElement.State)) == mashupsdk.Clicked {
								changedElement.State |= int64(mashupsdk.RightClick)
							} else {
								changedElement.State &= ^int64(mashupsdk.RightClick)
							}
						}
					}

					elementStateBundle := mashupsdk.MashupElementStateBundle{
						AuthToken:     worldApp.GetAuthToken(),
						ElementStates: changedElements,
					}

					w.MashupContext.Client.UpsertMashupElementsState(w.MashupContext, &elementStateBundle)
				}
			}

		})

		// Create and add lights to the scene
		w.scene.Add(light.NewAmbient(g3ndpalette.WHITE, 0.8))
		pointLight := light.NewPoint(g3ndpalette.WHITE, 5.0)
		pointLight.SetPosition(1, 0, 2)
		w.scene.Add(pointLight)

		// Create and add an axis helper to the scene
		w.scene.Add(helper.NewAxes(0.5))

		w.frameRater.Start()
		// Set background color to gray
		if w.backgroundG3n != nil {
			w.IG3nRenderer.RenderElement(w, w.backgroundG3n)
			g3ndpalette.RefreshBackgroundColor(w.mainWin.Gls(), w.backgroundG3n.GetColor(), 1.0)
		}
		go func() {
			log.Println("Watching position events.")
			for displayHint := range w.displayPositionChan {
				log.Printf("G3n applying xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
				if !w.headless {
					if w.mainWin != nil {
						if iWindow, iWindowOk := (*w.mainWin).IWindow.(*window.GlfwWindow); iWindowOk {
							if worldApp.IG3nDisplayRenderer != nil {
								worldApp.IG3nDisplayRenderer.Render((*w.mainWin).IWindow.(*window.GlfwWindow), displayHint)
							} else {
								if !w.custos {
									iWindow.Window.SetAttrib(glfw.Decorated, 0)
									iWindow.Window.SetAttrib(glfw.FocusOnShow, 1)
									iWindow.Window.SetAttrib(glfw.Stereo, 1)
								}
								if x, y := iWindow.Window.GetPos(); x != int(displayHint.Xpos) || y != int(displayHint.Ypos+displayHint.Height) {
									iWindow.Window.SetPos(int(displayHint.Xpos), int(displayHint.Ypos+displayHint.Height))
									iWindow.Window.SetSize(int(displayHint.Width), int(displayHint.Height))
								}
								if !w.Focused && displayHint.Focused {
									log.Printf("G3n setting focus.")
									iWindow.Window.Hide()
									iWindow.Window.Show()
									displayHint.Focused = false
									w.Focused = true
								}
							}
						}
					}
				}
			}
			log.Println("Exiting disply chan.")
		}()
		log.Printf("InitHandler complete.")
	}
	runtimeHandler := func(renderer *renderer.Renderer, deltaTime time.Duration) {
		if iWindow, iWindowOk := (*w.mainWin).IWindow.(*window.GlfwWindow); iWindowOk {

			if iWindow.Window.GetAttrib(glfw.Focused) != 1 {
				w.SetFrameRate(0)
			} else {
				if w.currentTargetFPS == 0 {
					w.SetFrameRate(5)
				}
			}
		}

		w.frameRater.Start()

		if w.currentTargetFPS != 0 {
			// Nothing to render.
			if w.backgroundG3n != nil && w.backgroundG3n.GetColor() != nil {
				g3ndpalette.RefreshBackgroundColor(w.mainWin.Gls(), w.backgroundG3n.GetColor(), 1.0)
			}
			w.mainWin.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
			renderer.Render(w.scene, w.cam)

			if !w.isInit {
				w.G3nOnFocus("", InitEvent{})
				w.isInit = true
			}
		}

		w.frameRater.Wait()
	}

	guiboot.InitMainWindow(guiboot.G3n, initHandler, runtimeHandler)
}

func (w *worldClientInitHandler) RegisterContext(context *mashupsdk.MashupContext) {
	worldApp.MashupContext = context
}

// Sets all elements to a "Rest state."
func (w *mashupSdkApiHandler) ResetG3NDetailedElementStates() {
	log.Printf("G3n Received ResetG3NDetailedElementStates\n")
	for _, wes := range worldApp.ConcreteElements {
		wes.SetElementState(mashupsdk.Init)
	}
	log.Printf("G3n finished ResetG3NDetailedElementStates handle.\n")
}

func (mSdk *mashupSdkApiHandler) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	if worldApp.mainWin != nil && (*worldApp.mainWin).IWindow != nil {
		log.Printf("G3n Received onResize xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
		worldApp.displayPositionChan <- displayHint
	} else {
		if displayHint.Width != 0 && displayHint.Height != 0 {
			log.Printf("G3n initializing with: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
			worldApp.displaySetupChan <- displayHint
			worldApp.displayPositionChan <- displayHint
		} else {
			log.Printf("G3n Could not apply xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
		}
		log.Printf("G3n finished onResize handle.")
	}
}

func (w *mashupSdkApiHandler) GetMashupElements() (*mashupsdk.MashupDetailedElementBundle, error) {
	log.Printf("G3n Received GetMashupElements\n")
	log.Printf("G3n finished GetMashupElements handle.\n")
	return nil, nil
}

func (mSdk *mashupSdkApiHandler) UpsertMashupElements(detailedElementBundle *mashupsdk.MashupDetailedElementBundle) (*mashupsdk.MashupDetailedElementBundle, error) {
	log.Printf("G3n Received UpsertMashupElements\n")
	result := &mashupsdk.MashupDetailedElementBundle{DetailedElements: []*mashupsdk.MashupDetailedElement{}}
	incompleteG3nElements := []*g3nmash.G3nDetailedElement{}

	for _, detailedElement := range detailedElementBundle.DetailedElements {
		g3nDetailedElement := worldApp.NewG3nDetailedElement(detailedElement, false)
		if g3nDetailedElement.IsLibraryElement() {
			continue
		}

		if detailedElement.State.Id != int64(mashupsdk.Immutable) {
			g3nDetailedElement.ApplyState(mashupsdk.Clicked, false)
		}

		for _, childId := range g3nDetailedElement.GetChildElementIds() {
			if childId < 0 {
				incompleteG3nElements = append(incompleteG3nElements, g3nDetailedElement)
				break
			}
		}
		if worldApp.maxElementId < g3nDetailedElement.GetDisplayId() {
			worldApp.maxElementId = g3nDetailedElement.GetDisplayId()
		}

		// Add to resulting element states.
		result.DetailedElements = append(result.DetailedElements, detailedElement)
	}

	if len(incompleteG3nElements) > 0 {
		// Fill out incomplete g3n elements
		generatedElements := []interface{}{}
		for _, incompleteG3nElement := range incompleteG3nElements {
			newChildIds := []int64{}

			for _, childId := range incompleteG3nElement.GetChildElementIds() {
				if childId < 0 {
					if libElement, err := worldApp.GetG3nDetailedLibraryElementById(childId); err == nil {
						clonedChild := worldApp.CloneG3nDetailedElement(libElement, &generatedElements)
						newChildIds = append(newChildIds, clonedChild.GetDisplayId())
					} else {
						log.Printf("Missing child from library: %d\n", childId)
					}
				} else {
					// Deal with concrete element.
					if concreteElement, err := worldApp.GetG3nDetailedElementById(childId); err == nil {
						newChildIds = append(newChildIds, concreteElement.GetDisplayId())
					}
				}
			}
			if len(newChildIds) > 0 {
				incompleteG3nElement.SetChildElements(newChildIds)
			}
		}
		for _, generatedElement := range generatedElements {
			result.DetailedElements = append(result.DetailedElements, generatedElement.(*mashupsdk.MashupDetailedElement))
		}
	}
	// Reset cache so it'll be rebuilt.
	worldApp.RootElements = []*g3nmash.G3nDetailedElement{}

	log.Printf("G3n UpsertMashupElements updated\n")
	return result, nil
}

func (mSdk *mashupSdkApiHandler) applyStateHelper(g3nId int64, x mashupsdk.DisplayElementState, isset bool) {

	child := worldApp.ConcreteElements[g3nId]
	child.ApplyState(mashupsdk.DisplayElementState(x), isset)

	if len(child.GetDetailedElement().Childids) > 0 {
		for _, cId := range child.GetDetailedElement().Childids {
			mSdk.applyStateHelper(cId, x, isset)
		}
	}
}

func (mSdk *mashupSdkApiHandler) setStateHelper(g3nId int64, x mashupsdk.DisplayElementState) {

	child := worldApp.ConcreteElements[g3nId]
	if child.GetDetailedElement().Genre != "Attitude" {
		child.SetElementState(mashupsdk.DisplayElementState(x))
	}

	if len(child.GetDetailedElement().Childids) > 0 {
		for _, cId := range child.GetDetailedElement().Childids {
			mSdk.setStateHelper(cId, x)
		}
	}
}

func (mSdk *mashupSdkApiHandler) UpsertMashupElementsState(elementStateBundle *mashupsdk.MashupElementStateBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("G3n UpsertMashupElementsState called\n")

	ClickedElements := map[int64]*g3nmash.G3nDetailedElement{}
	recursiveElements := map[int64]*g3nmash.G3nDetailedElement{}

	for _, es := range elementStateBundle.ElementStates {
		if g3nDetailedElement, ok := worldApp.ConcreteElements[es.GetId()]; ok {
			g3nDetailedElement.SetElementState(mashupsdk.DisplayElementState(es.State))
			if g3nDetailedElement.IsStateSet(mashupsdk.Recursive) {
				recursiveElements[es.GetId()] = g3nDetailedElement
			}

			log.Printf("Display fields set to: %v", g3nDetailedElement.GetMashupElementState())
			if (mashupsdk.DisplayElementState(es.State) & mashupsdk.Clicked) == mashupsdk.Clicked {
				ClickedElements[es.GetId()] = g3nDetailedElement
			}
		}
	}

	if len(ClickedElements) > 0 {
		// Remove existing clicks.
		for _, clickedElement := range worldApp.ClickedElements {
			if _, ok := ClickedElements[clickedElement.GetDisplayId()]; !ok {
				clickedElement.ApplyState(mashupsdk.Clicked, false)
			}
		}

		worldApp.ClickedElements = worldApp.ClickedElements[:0]

		// Impossible to determine ordering of clicks from upsert at this time.
		for _, g3nDetailedElement := range ClickedElements {
			worldApp.ClickedElements = append(worldApp.ClickedElements, g3nDetailedElement)
		}
	}

	if len(recursiveElements) > 0 {
		for _, recursiveElement := range recursiveElements {
			stateBits := recursiveElement.GetDetailedElement().State.State
			// Unset recursive for child elements
			stateBits &= ^int64(mashupsdk.Recursive)
			// Apply this state change to all child elements.
			mSdk.setStateHelper(recursiveElement.GetDisplayId(), mashupsdk.DisplayElementState(stateBits))
		}
	}

	log.Printf("G3n dispatching focus\n")
	// TODO: Feedback from Custos is broken...  Disable for now.
	if worldApp.mainWin != nil {
		worldApp.mainWin.Dispatch(gui.OnFocus, nil)
	}
	log.Printf("G3n End UpsertMashupElementsState called\n")
	return &mashupsdk.MashupElementStateBundle{}, nil
}
