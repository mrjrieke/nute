package g3nworld

import (
	"errors"
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
	"tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/guiboot"
	"tini.com/nute/mashupsdk/server"
)

type mashupSdkApiHandler struct {
}

type worldClientInitHandler struct {
}

type WorldApp struct {
	mSdkApiHandler      *mashupSdkApiHandler
	wClientInitHandler  *worldClientInitHandler
	displaySetupChan    chan *mashupsdk.MashupDisplayHint
	displayPositionChan chan *mashupsdk.MashupDisplayHint
	mainWin             *app.Application
	frameRater          *util.FrameRater // Render loop frame rater
	scene               *core.Node
	cam                 *camera.Camera
	oc                  *camera.OrbitControl

	mashupContext     *mashupsdk.MashupContext              // Needed for callbacks to other mashups
	elementIndex      map[int64]*g3nmash.G3nDetailedElement // g3n indexes by string...
	elementDictionary map[string]int64
	isInit            bool
}

var worldApp WorldApp

func NewWorldApp() *WorldApp {
	worldApp = WorldApp{
		mSdkApiHandler:      &mashupSdkApiHandler{},
		elementIndex:        map[int64]*g3nmash.G3nDetailedElement{},
		elementDictionary:   map[string]int64{},
		displaySetupChan:    make(chan *mashupsdk.MashupDisplayHint, 1),
		displayPositionChan: make(chan *mashupsdk.MashupDisplayHint, 1),
	}
	return &worldApp
}

type InitEvent struct {
}

func (w *WorldApp) G3nOnFocus(name string, ev interface{}) {
	log.Printf("G3nWorld Focus gained\n")

	if _, iOk := ev.(InitEvent); iOk {

		torusG3ns, err := w.GetG3nDetailedFilteredElements("torus")
		if err != nil {
			log.Fatal(err)
		}

		for _, torusG3n := range torusG3ns {
			torusGeom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
			mat := material.NewStandard(math32.NewColor("DarkBlue"))
			torusMesh := graphic.NewMesh(torusGeom, mat)
			torusMesh.SetLoaderID(torusG3n.GetDisplayName())
			torusMesh.SetPositionVec(math32.NewVector3(float32(0.0), float32(0.0), float32(0.0)))
			w.scene.Add(torusMesh)
			torusG3n.SetNamedMesh(torusG3n.GetDisplayName(), torusMesh)

			if torusInside, tidErr := w.GetG3nDetailedElementById(torusG3n.GetChildElements()[0]); tidErr == nil {
				diskGeom := geometry.NewDisk(1, 32)
				diskMat := material.NewStandard(&math32.Color{R: 0.5, G: 0.5, B: 0.5})
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
		torus, _ := w.GetG3nDetailedElement("torus")
		torusInnerDisk, _ := w.GetG3nDetailedElementById(torus.GetChildElements()[0])

		for _, g3nDetailedElement := range w.elementIndex {
			if g3nDetailedElement.GetDisplayState() != mashupsdk.Rest {
				switch g3nDetailedElement.GetDisplayId() {
				case 1:
					log.Printf("G3n Inside\n")
					torus.SetRotationX(0)
					torus.SetColor(math32.NewColor("DarkBlue"))

					torusInnerDisk.SetRotationX(0)
					torusInnerDisk.SetColor(math32.NewColor("DarkRed"))
				case 2:
					log.Printf("G3n Outside\n")
					torus.SetRotationX(0)
					torus.SetColor(math32.NewColor("DarkBlue"))

					torusInnerDisk.SetRotationX(0)
					torusInnerDisk.SetColor(&math32.Color{R: 0.5, G: 0.5, B: 0.5})
				case 3:
					log.Printf("G3n It\n")
					torus.SetRotationX(0)
					torus.SetColor(math32.NewColor("DarkRed"))

					torusInnerDisk.SetRotationX(0)
					torusInnerDisk.SetColor(&math32.Color{R: 0.5, G: 0.5, B: 0.5})
				case 4:
					log.Printf("G3n Up-Side-Down\n")
					torus.SetRotationX(180)
					torus.SetColor(math32.NewColor("DarkBlue"))

					torusInnerDisk.SetRotationX(180)
				}
			}
		}
		log.Printf("G3n End Focus gained\n")
	}

	log.Printf("G3nWorld End Focus gained\n")
}

func (w *WorldApp) ResetChangeStates() []*mashupsdk.MashupElementState {
	changedElements := []*mashupsdk.MashupElementState{}
	for _, g3nDetailedElement := range w.elementIndex {
		if g3nDetailedElement.GetDisplayState() != mashupsdk.Rest {
			g3nDetailedElement.SetDisplayState(mashupsdk.Rest)
			changedElements = append(changedElements, g3nDetailedElement.GetMashupElementState())
		}
	}

	return changedElements
}

// Sets all elements to a "Rest state."
func (w *WorldApp) ResetG3nDetailedElementStates() {
	for _, wes := range w.elementIndex {
		wes.SetDisplayState(mashupsdk.Rest)
	}
}

func (w *WorldApp) NewG3nDetailedElement(detailedElement *mashupsdk.MashupDetailedElement) *g3nmash.G3nDetailedElement {
	w.elementDictionary[detailedElement.GetName()] = detailedElement.Id
	g3nDetailedElement := g3nmash.NewG3nDetailedElement(detailedElement)
	w.elementIndex[detailedElement.Id] = g3nDetailedElement
	return g3nDetailedElement
}

func (w *WorldApp) GetG3nDetailedFilteredElements(elementPrefix string) ([]*g3nmash.G3nDetailedElement, error) {
	filteredElements := []*g3nmash.G3nDetailedElement{}
	for _, element := range w.elementIndex {
		if strings.HasPrefix(element.GetDisplayName(), elementPrefix) {
			filteredElements = append(filteredElements, element)
		}
	}

	return filteredElements, nil
}

func (w *WorldApp) GetG3nDetailedElement(elementName string) (*g3nmash.G3nDetailedElement, error) {
	if eid, eidOk := w.elementDictionary[elementName]; eidOk {
		if g3nElement, g3nElementOk := w.elementIndex[eid]; g3nElementOk {
			return g3nElement, nil
		}
	}
	return nil, errors.New("element does not exist: " + elementName)
}

func (w *WorldApp) GetG3nDetailedElementById(eid int64) (*g3nmash.G3nDetailedElement, error) {
	if g3nElement, g3nElementOk := w.elementIndex[eid]; g3nElementOk {
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
		server.InitServer(callerCreds, insecure, w.mSdkApiHandler, w.wClientInitHandler)
	} else {
		go func() {
			w.displaySetupChan <- &mashupsdk.MashupDisplayHint{Xpos: 0, Ypos: 0, Width: 400, Height: 800}
		}()
	}
}

func (w *WorldApp) InitMainWindow() {
	log.Printf("Initializing mainWin.")

	initHandler := func(a *app.Application) {
		log.Printf("InitHandler.")

		if w.mainWin == nil {
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
				n, intersections := w.Cast(w.scene, caster)
				if len(intersections) != 0 {
					if len(w.elementDictionary) != 0 {
						g3nElement, err := w.GetG3nDetailedElement(n.GetNode().LoaderID())
						if err != nil {
							log.Fatal(err)
						}

						log.Printf("State: %d\n", g3nElement.GetDisplayState())

						if g3nElement.GetMashupElementState() != nil {
							log.Printf("State size: %d\n", len(w.elementDictionary))
							g3nElement.SetColor(math32.NewColor("DarkRed"))
							// Zero out states of all elements to rest state.
							changedStates := w.ResetChangeStates()

							g3nElement.SetDisplayState(mashupsdk.Clicked)
							changedStates = append(changedStates, g3nElement.GetMashupElementState())

							elementStateBundle := mashupsdk.MashupElementStateBundle{
								AuthToken:     server.GetServerAuthToken(),
								ElementStates: changedStates,
							}

							w.mashupContext.Client.UpsertMashupElementsState(w.mashupContext, &elementStateBundle)
						}
					}

				} else {
					log.Printf("No intersection found\n")
					// Nothing selected...
					if torusG3n, tidErr := w.GetG3nDetailedElement("torus"); tidErr == nil {
						torusG3n.SetColor(math32.NewColor("DarkBlue"))
					}
					if len(w.elementDictionary) != 0 {
						changedElements := w.ResetChangeStates()

						g3nElement, err := w.GetG3nDetailedElement("Outside")
						if err != nil {
							log.Fatal(err)
						}

						g3nElement.SetDisplayState(mashupsdk.Clicked)
						changedElements = append(changedElements, g3nElement.GetMashupElementState())

						elementStateBundle := mashupsdk.MashupElementStateBundle{
							AuthToken:     server.GetServerAuthToken(),
							ElementStates: changedElements,
						}

						w.mashupContext.Client.UpsertMashupElementsState(w.mashupContext, &elementStateBundle)
					}
				}
			}

		})

		// Create and add lights to the scene
		w.scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
		pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
		pointLight.SetPosition(1, 0, 2)
		w.scene.Add(pointLight)

		// Create and add an axis helper to the scene
		w.scene.Add(helper.NewAxes(0.5))

		w.frameRater.Start()
		// Set background color to gray
		a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
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
		for _, g3nDetailedElement := range w.elementIndex {
			if g3nDetailedElement.GetDisplayState() != mashupsdk.Rest {
				switch g3nDetailedElement.GetDisplayId() {
				case 1:
					// Inside
					w.mainWin.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
				case 2:
					// Outside
					// Updates the background to dark red...
					w.mainWin.Gls().ClearColor(.545, 0, 0, 1.0)
				case 3:
					// It
					w.mainWin.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
				case 4:
					// Up-side-down
					w.mainWin.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
				}
				break
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

func (mSdk *mashupSdkApiHandler) UpsertMashupElements(detailedElementBundle *mashupsdk.MashupDetailedElementBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("G3n Received UpsertMashupElements\n")
	result := &mashupsdk.MashupElementStateBundle{ElementStates: []*mashupsdk.MashupElementState{}}

	for _, detailedElement := range detailedElementBundle.DetailedElements {
		g3nDetailedElement := worldApp.NewG3nDetailedElement(detailedElement)
		g3nDetailedElement.SetDisplayState(mashupsdk.Rest)

		// Add to resulting element states.
		result.ElementStates = append(result.ElementStates, detailedElement.State)
	}

	log.Printf("G3n UpsertMashupElements updated\n")
	return result, nil
}

func (mSdk *mashupSdkApiHandler) UpsertMashupElementsState(elementStateBundle *mashupsdk.MashupElementStateBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("G3n UpsertMashupElementsState called\n")

	worldApp.ResetG3nDetailedElementStates()

	for _, es := range elementStateBundle.ElementStates {
		if worldApp.elementIndex[es.GetId()].GetDisplayState() != mashupsdk.DisplayElementState(es.State) {
			worldApp.elementIndex[es.GetId()].SetDisplayState(mashupsdk.DisplayElementState(es.State))
		}
	}
	worldApp.mainWin.Dispatch(gui.OnFocus, nil)
	log.Printf("G3n End UpsertMashupElementsState called\n")
	return &mashupsdk.MashupElementStateBundle{}, nil
}
