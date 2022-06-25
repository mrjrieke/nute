package g3nworld

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"github.com/go-gl/glfw/v3.3/glfw"
	"tini.com/nute/g3nd/g3nmash"
	g3ndpalette "tini.com/nute/g3nd/palette"
	"tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/guiboot"
	"tini.com/nute/mashupsdk/server"
)

type mashupSdkApiHandler struct {
}

type worldClientInitHandler struct {
}

type WorldApp struct {
	headless            bool // Mode for troubleshooting.
	MSdkApiHandler      *mashupSdkApiHandler
	wClientInitHandler  *worldClientInitHandler
	displaySetupChan    chan *mashupsdk.MashupDisplayHint
	displayPositionChan chan *mashupsdk.MashupDisplayHint
	mainWin             *app.Application
	frameRater          *util.FrameRater // Render loop frame rater
	scene               *core.Node
	cam                 *camera.Camera
	oc                  *camera.OrbitControl

	mashupContext *mashupsdk.MashupContext // Needed for callbacks to other mashups

	// Library for mashup objects
	elementLibraryDictionary map[int64]*g3nmash.G3nDetailedElement

	maxElementId      int64
	elementDictionary map[int64]*g3nmash.G3nDetailedElement // g3n indexes by string...
	isInit            bool
}

var worldApp WorldApp

func NewWorldApp(headless bool) *WorldApp {
	worldApp = WorldApp{
		headless:                 headless,
		MSdkApiHandler:           &mashupSdkApiHandler{},
		elementLibraryDictionary: map[int64]*g3nmash.G3nDetailedElement{},
		elementDictionary:        map[int64]*g3nmash.G3nDetailedElement{},
		displaySetupChan:         make(chan *mashupsdk.MashupDisplayHint, 1),
		displayPositionChan:      make(chan *mashupsdk.MashupDisplayHint, 1),
	}
	return &worldApp
}

type InitEvent struct {
}

func (w *WorldApp) G3nOnFocus(name string, ev interface{}) {
	log.Printf("G3nWorld Focus gained\n")

	if _, iOk := ev.(InitEvent); iOk {

		torusG3ns, err := w.GetG3nDetailedFilteredElements("Torus")
		if err != nil {
			log.Fatal(err)
		}

		for _, torusG3n := range torusG3ns {
			torusConcreteG3n := torusG3n
			if torusG3n.IsAbstract() {
				if tc, tErr := worldApp.GetG3nDetailedElementById(torusG3n.GetChildElements()[0]); tErr == nil {
					torusConcreteG3n = tc
				}
			}

			torusGeom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
			mat := material.NewStandard(g3ndpalette.DARK_BLUE)
			torusMesh := graphic.NewMesh(torusGeom, mat)
			torusMesh.SetLoaderID(torusConcreteG3n.GetDisplayName())
			torusMesh.SetPositionVec(math32.NewVector3(float32(0.0), float32(0.0), float32(0.0)))
			w.scene.Add(torusMesh)
			torusConcreteG3n.SetNamedMesh(torusConcreteG3n.GetDisplayName(), torusMesh)

			for _, torusInside := range w.GetG3nDetailedElementsByGenre(torusConcreteG3n, "Space") {
				diskGeom := geometry.NewDisk(1, 32)
				diskMat := material.NewStandard(g3ndpalette.GREY)
				diskMesh := graphic.NewMesh(diskGeom, diskMat)
				diskMesh.SetPositionVec(math32.NewVector3(float32(0.0), float32(0.0), float32(0.0)))
				diskMesh.SetLoaderID(torusInside.GetDisplayName())
				w.scene.Add(diskMesh)
				torusInside.SetNamedMesh(torusInside.GetDisplayName(), diskMesh)
			}
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
	for _, g3nDetailedElement := range w.elementDictionary {
		if g3nDetailedElement.GetDisplayState() != mashupsdk.Rest {
			g3nDetailedElement.SetDisplayState(mashupsdk.Rest)
			changedElements = append(changedElements, g3nDetailedElement.GetMashupElementState())
		}
	}

	return changedElements
}

// Sets all elements to a "Rest state."
func (w *WorldApp) ResetG3nDetailedElementStates() {
	for _, wes := range w.elementDictionary {
		wes.SetDisplayState(mashupsdk.Rest)
	}
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
	if g3nDetailedElement.GetBasisId() < 0 {
		w.elementLibraryDictionary[g3nDetailedElement.GetBasisId()] = g3nDetailedElement
		if g3nDetailedElement.GetDisplayId() > 0 {
			w.elementDictionary[g3nDetailedElement.GetDisplayId()] = g3nDetailedElement
		}
	} else {
		w.elementDictionary[g3nDetailedElement.GetDisplayId()] = g3nDetailedElement
	}
	return g3nDetailedElement
}

func (w *WorldApp) GetG3nDetailedFilteredElements(elementPrefix string) ([]*g3nmash.G3nDetailedElement, error) {
	filteredElements := []*g3nmash.G3nDetailedElement{}
	for _, element := range w.elementDictionary {
		if strings.HasPrefix(element.GetDisplayName(), elementPrefix) {
			filteredElements = append(filteredElements, element)
		}
	}

	return filteredElements, nil
}

func (w *WorldApp) GetG3nDetailedElementById(eid int64) (*g3nmash.G3nDetailedElement, error) {
	if g3nElement, g3nElementOk := w.elementDictionary[eid]; g3nElementOk {
		return g3nElement, nil
	}
	return nil, fmt.Errorf("element does not exist: %d", eid)
}

func (w *WorldApp) GetG3nDetailedElementsByGenre(g3n *g3nmash.G3nDetailedElement, genre string) []*g3nmash.G3nDetailedElement {
	results := []*g3nmash.G3nDetailedElement{}
	for _, childId := range g3n.GetChildElements() {
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

func (w *WorldApp) InitServer(callerCreds string, insecure bool) {
	if callerCreds != "" {
		server.InitServer(callerCreds, insecure, w.MSdkApiHandler, w.wClientInitHandler)
	} else {
		go func() {
			w.displaySetupChan <- &mashupsdk.MashupDisplayHint{Xpos: 0, Ypos: 0, Width: 400, Height: 800}
		}()
	}
}

func (w *WorldApp) Transform() []*mashupsdk.MashupElementState {
	itemColor := g3ndpalette.DARK_BLUE
	itemClickedColor := g3ndpalette.DARK_RED

	changedElements := []*mashupsdk.MashupElementState{}
	visitedNodes := map[int64]bool{}
	for _, g3nDetailedElement := range w.elementDictionary {
		var changed bool
		g3nColor := itemColor
		if g3nDetailedElement.IsItemActive() {
			g3nColor = itemClickedColor
			if g3nDetailedElement.HasAttitudeAdjustment() {
				log.Printf("G3n Has parents\n")
				parentIds := g3nDetailedElement.GetParentElements()
				g3nParentDetailedElements := []*g3nmash.G3nDetailedElement{}
				for _, parentId := range parentIds {
					if g3parent, gpErr := w.GetG3nDetailedElementById(parentId); gpErr == nil {
						g3nParentDetailedElements = append(g3nParentDetailedElements, g3parent)
					}
					visitedNodes[parentId] = true
				}
				log.Printf("G3n adjusting for parents: %d\n", len(g3nParentDetailedElements))

				g3nDetailedElement.AdjustAttitude(g3nParentDetailedElements)
			} else {
				if _, vOk := visitedNodes[g3nDetailedElement.GetDisplayId()]; !vOk {
					g3nDetailedElement.AdjustAttitude([]*g3nmash.G3nDetailedElement{g3nDetailedElement})
					visitedNodes[g3nDetailedElement.GetDisplayId()] = true
				}
			}
		} else {
			if g3nDetailedElement.IsAbstract() {
				continue
			}

			if g3nDetailedElement.IsBackground() {
				if g3nDetailedElement.IsItemActive() {
					// No items clicked means background is clicked.
					g3nColor = itemClickedColor
				} else {
					g3nColor = g3ndpalette.GREY
				}
			} else {
				if g3nDetailedElement.IsBackgroundElement() {
					g3nColor = g3ndpalette.GREY
				}
			}
			if _, vOk := visitedNodes[g3nDetailedElement.GetDisplayId()]; !vOk {
				g3nDetailedElement.AdjustAttitude([]*g3nmash.G3nDetailedElement{g3nDetailedElement})
				visitedNodes[g3nDetailedElement.GetDisplayId()] = true
			}
		}
		changed = g3nDetailedElement.SetColor(g3nColor)

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
		w.frameRater = util.NewFrameRater(1)
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
		log.Printf("Finished Orbit Control setup.")

		// Set up callback to update viewport and camera aspect ratio when the window is resized
		onResize := func(evname string, ev interface{}) {
			// Get framebuffer size and update viewport accordingly
			width, height := a.GetSize()
			a.Gls().Viewport(0, 0, int32(width), int32(height))
			// Update the camera's aspect ratio
			w.cam.SetAspect(float32(width) / float32(height))

			xpos, ypos := (*w.mainWin).IWindow.(*window.GlfwWindow).Window.GetPos()

			if w.mashupContext != nil {
				w.mashupContext.Client.OnResize(w.mashupContext,
					&mashupsdk.MashupDisplayBundle{
						AuthToken: server.GetServerAuthToken(),
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

		(*w.mainWin).IWindow.(*window.GlfwWindow).Window.SetCloseCallback(func(glfwWindow *glfw.Window) {
			if w.mashupContext != nil {
				w.mashupContext.Client.Shutdown(w.mashupContext, &mashupsdk.MashupEmpty{AuthToken: server.GetServerAuthToken()})
			}
		})

		w.mainWin.Subscribe(gui.OnMouseUp, func(name string, ev interface{}) {
			mev := ev.(*window.MouseEvent)
			g3Width, g3Height := w.mainWin.GetSize()

			xPosNdc := 2*(mev.Xpos/float32(g3Width)) - 1
			yPosNdc := -2*(mev.Ypos/float32(g3Height)) + 1
			caster := collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
			caster.SetFromCamera(w.cam, xPosNdc, yPosNdc)

			if w.scene.Visible() {
				itemClicked, _ := w.Cast(w.scene, caster)

				itemMatched := false
				var backgroundG3n *g3nmash.G3nDetailedElement
				for _, g3nDetailedElement := range w.elementDictionary {
					if g3nDetailedElement.IsBackground() {
						backgroundG3n = g3nDetailedElement
					} else {
						if g3nDetailedElement.IsItemClicked(itemClicked) {
							g3nDetailedElement.SetDisplayState(mashupsdk.Clicked)
							itemMatched = true
						} else {
							g3nDetailedElement.SetDisplayState(mashupsdk.Rest)
						}
					}
				}
				if !itemMatched {
					backgroundG3n.SetDisplayState(mashupsdk.Clicked)
				} else {
					backgroundG3n.SetDisplayState(mashupsdk.Rest)
				}
				changedElements := w.Transform()
				if !itemMatched {
					changedElements = append(changedElements, backgroundG3n.GetMashupElementState())
				}

				elementStateBundle := mashupsdk.MashupElementStateBundle{
					AuthToken:     server.GetServerAuthToken(),
					ElementStates: changedElements,
				}

				if !w.headless {
					w.mashupContext.Client.UpsertMashupElementsState(w.mashupContext, &elementStateBundle)
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
		g3ndpalette.RefreshBackgroundColor(a.Gls(), g3ndpalette.GREY, 1.0)
		go func() {
			log.Println("Watching position events.")
			for displayHint := range w.displayPositionChan {
				log.Printf("G3n applying xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
				(*w.mainWin).IWindow.(*window.GlfwWindow).Window.SetPos(int(displayHint.Xpos), int(displayHint.Ypos+displayHint.Height))
				(*w.mainWin).IWindow.(*window.GlfwWindow).Window.SetSize(int(displayHint.Width), int(displayHint.Height))
			}
			log.Println("Exiting disply chan.")
		}()
		log.Printf("InitHandler complete.")
	}
	runtimeHandler := func(renderer *renderer.Renderer, deltaTime time.Duration) {
		for _, g3nDetailedElement := range w.elementDictionary {
			if g3nDetailedElement.GetDisplayState() != mashupsdk.Rest {
				if g3nDetailedElement.IsBackground() {
					g3ndpalette.RefreshBackgroundColor(w.mainWin.Gls(), g3ndpalette.DARK_RED, 1.0)
				} else {
					g3ndpalette.RefreshBackgroundColor(w.mainWin.Gls(), g3ndpalette.GREY, 1.0)
				}
			}
		}
		w.mainWin.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(w.scene, w.cam)

		if !w.isInit {
			w.G3nOnFocus("", InitEvent{})
			w.isInit = true
		}
		w.frameRater.Wait()
	}

	guiboot.InitMainWindow(guiboot.G3n, initHandler, runtimeHandler)
}

func (w *worldClientInitHandler) RegisterContext(context *mashupsdk.MashupContext) {
	worldApp.mashupContext = context
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
			g3nDetailedElement.SetDisplayState(mashupsdk.Rest)
		}

		for _, childId := range g3nDetailedElement.GetChildElements() {
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

			for _, childId := range incompleteG3nElement.GetChildElements() {
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

	log.Printf("G3n UpsertMashupElements updated\n")
	return result, nil
}

func (mSdk *mashupSdkApiHandler) UpsertMashupElementsState(elementStateBundle *mashupsdk.MashupElementStateBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("G3n UpsertMashupElementsState called\n")

	worldApp.ResetG3nDetailedElementStates()

	for _, es := range elementStateBundle.ElementStates {
		if worldApp.elementDictionary[es.GetId()].GetDisplayState() != mashupsdk.DisplayElementState(es.State) {
			worldApp.elementDictionary[es.GetId()].SetDisplayState(mashupsdk.DisplayElementState(es.State))
		}
	}
	log.Printf("G3n dispatching focus\n")
	if worldApp.mainWin != nil {
		worldApp.mainWin.Dispatch(gui.OnFocus, nil)
	}
	log.Printf("G3n End UpsertMashupElementsState called\n")
	return &mashupsdk.MashupElementStateBundle{}, nil
}