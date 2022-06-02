//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"flag"
	"log"
	"os"
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
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"github.com/go-gl/glfw/v3.3/glfw"
	"tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/guiboot"
	"tini.com/nute/mashupsdk/server"
)

var worldCompleteChan chan bool

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
	scene               *core.Node
	cam                 *camera.Camera

	mashupContext      *mashupsdk.MashupContext                // Needed for callbacks to other mashups
	elementIndex       map[int64]*mashupsdk.MashupElementState // g3n indexes by string...
	elementDictionary  map[string]int64
	DetailedElements   []*mashupsdk.MashupDetailedElement
	elementStateBundle mashupsdk.MashupElementStateBundle
}

var worldApp WorldApp

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

func (w *WorldApp) InitMainWindow() {
	log.Printf("Initializing mainWin.")

	initHandler := func(a *app.Application) {
		log.Printf("InitHandler.")

		if w.mainWin == nil {
			w.mainWin = a
		}
		displayHint := <-worldApp.displaySetupChan
		app.AppCustom(a, "Hello world G3n", int(displayHint.Width), int(displayHint.Height), int(displayHint.Xpos), int(displayHint.Ypos+displayHint.Height))
		w.scene = core.NewNode()

		// Set the scene to be managed by the gui manager
		gui.Manager().Set(w.scene)

		// Create perspective camera
		w.cam = camera.New(1)
		w.cam.SetPosition(0, 0, 3)
		w.scene.Add(w.cam)

		// Set up orbit control for the camera
		camera.NewOrbitControl(w.cam)

		// Set up callback to update viewport and camera aspect ratio when the window is resized
		onResize := func(evname string, ev interface{}) {
			// Get framebuffer size and update viewport accordingly
			width, height := a.GetSize()
			a.Gls().Viewport(0, 0, int32(width), int32(height))
			// Update the camera's aspect ratio
			w.cam.SetAspect(float32(width) / float32(height))

			xpos, ypos := (*worldApp.mainWin).IWindow.(*window.GlfwWindow).Window.GetPos()

			if worldApp.mashupContext != nil {
				worldApp.mashupContext.Client.OnResize(worldApp.mashupContext,
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

		// Create a blue torus and add it to the scene
		torusGeom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
		mat := material.NewStandard(math32.NewColor("DarkBlue"))
		mesh := graphic.NewMesh(torusGeom, mat)
		mesh.SetLoaderID("torus")
		w.scene.Add(mesh)

		diskGeom := geometry.NewDisk(1, 32)
		diskMat := material.NewStandard(&math32.Color{R: 0.5, G: 0.5, B: 0.5})
		diskMesh := graphic.NewMesh(diskGeom, diskMat)
		diskMesh.SetLoaderID("Inside")
		w.scene.Add(diskMesh)

		// Create and add a button to the scene
		btn := gui.NewButton("Make Red")
		btn.SetPosition(100, 40)
		btn.SetSize(40, 40)
		btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
			mat.SetColor(math32.NewColor("DarkRed"))
		})
		w.scene.Add(btn)

		w.mainWin.Subscribe(gui.OnFocus, func(name string, ev interface{}) {
			// Focus gained...
			log.Printf("G3n Focus gained\n")
			for i := 0; i < len(worldApp.elementStateBundle.ElementStates); i++ {
				if worldApp.elementStateBundle.ElementStates[i].State != mashupsdk.Rest {
					switch worldApp.elementStateBundle.ElementStates[i].Id {
					case 1:
						log.Printf("G3n Inside\n")
						mesh.SetRotationX(0)
						diskMat.SetColor(math32.NewColor("DarkRed"))
						mat.SetColor(math32.NewColor("DarkBlue"))
					case 2:
						log.Printf("G3n Outside\n")
						mesh.SetRotationX(0)
						diskMat.SetColor(&math32.Color{R: 0.5, G: 0.5, B: 0.5})
						mat.SetColor(math32.NewColor("DarkBlue"))
					case 3:
						log.Printf("G3n It\n")
						mesh.SetRotationX(0)
						diskMat.SetColor(&math32.Color{R: 0.5, G: 0.5, B: 0.5})
						mat.SetColor(math32.NewColor("DarkRed"))
					case 4:
						log.Printf("G3n Up-Side-Down\n")
						mesh.SetRotationX(180)
						mat.SetColor(math32.NewColor("DarkBlue"))
					}

				}
			}
			log.Printf("G3n End Focus gained\n")
		})

		(*worldApp.mainWin).IWindow.(*window.GlfwWindow).Window.SetCloseCallback(func(w *glfw.Window) {
			worldApp.mashupContext.Client.Shutdown(worldApp.mashupContext, &mashupsdk.MashupEmpty{AuthToken: server.GetServerAuthToken()})
		})

		w.mainWin.Subscribe(gui.OnMouseUp, func(name string, ev interface{}) {
			mev := ev.(*window.MouseEvent)
			g3Width, g3Height := worldApp.mainWin.GetSize()

			xPosNdc := 2*(mev.Xpos/float32(g3Width)) - 1
			yPosNdc := -2*(mev.Ypos/float32(g3Height)) + 1
			caster := collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
			caster.SetFromCamera(worldApp.cam, xPosNdc, yPosNdc)

			if worldApp.scene.Visible() {
				n, intersections := worldApp.Cast(worldApp.scene, caster)
				if len(intersections) != 0 {
					if len(worldApp.elementDictionary) != 0 {
						lookupId := worldApp.elementDictionary[n.GetNode().LoaderID()]
						elementState := worldApp.elementIndex[lookupId]
						log.Printf("State: %d\n", elementState.State)

						if elementState != nil {
							log.Printf("State size: %d\n", len(worldApp.elementStateBundle.ElementStates))
							mat.SetColor(math32.NewColor("DarkRed"))
							// Zero out states of all elements to rest state.
							for i := 0; i < len(worldApp.elementStateBundle.ElementStates); i++ {
								if worldApp.elementStateBundle.ElementStates[i].State != mashupsdk.Rest {
									worldApp.elementStateBundle.ElementStates[i].State = mashupsdk.Rest
								}
							}
							elementState.State = mashupsdk.Clicked
							elementStateBundle := mashupsdk.MashupElementStateBundle{
								AuthToken:     server.GetServerAuthToken(),
								ElementStates: []*mashupsdk.MashupElementState{elementState},
							}

							worldApp.mashupContext.Client.UpsertMashupElementsState(worldApp.mashupContext, &elementStateBundle)
						}
					}

				} else {
					log.Printf("No intersection found\n")
					// Nothing selected...
					mat.SetColor(math32.NewColor("DarkBlue"))
					if len(worldApp.elementDictionary) != 0 {
						changedElements := []*mashupsdk.MashupElementState{}
						for i := 0; i < len(worldApp.elementStateBundle.ElementStates); i++ {
							if worldApp.elementStateBundle.ElementStates[i].State != mashupsdk.Rest {
								worldApp.elementStateBundle.ElementStates[i].State = mashupsdk.Rest
								changedElements = append(changedElements, worldApp.elementStateBundle.ElementStates[i])
							}
						}
						// TODO: determine whether click was inside or outside toroid
						// For now, append the 'outside' clicked.
						lookupId := worldApp.elementDictionary["Outside"]
						elementState := worldApp.elementIndex[lookupId]
						elementState.State = mashupsdk.Clicked
						changedElements = append(changedElements, elementState)

						elementStateBundle := mashupsdk.MashupElementStateBundle{
							AuthToken:     server.GetServerAuthToken(),
							ElementStates: changedElements,
						}

						worldApp.mashupContext.Client.UpsertMashupElementsState(worldApp.mashupContext, &elementStateBundle)
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

		// Set background color to gray
		a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
		go func() {
			log.Println("Watching position events.")
			for displayHint := range worldApp.displayPositionChan {
				log.Printf("G3n applying xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
				(*worldApp.mainWin).IWindow.(*window.GlfwWindow).Window.SetPos(int(displayHint.Xpos), int(displayHint.Ypos+displayHint.Height))
				(*worldApp.mainWin).IWindow.(*window.GlfwWindow).Window.SetSize(int(displayHint.Width), int(displayHint.Height))
			}
			log.Println("Exiting disply chan.")
		}()
		log.Printf("InitHandler complete.")
	}
	runtimeHandler := func(renderer *renderer.Renderer, deltaTime time.Duration) {
		for i := 0; i < len(worldApp.elementStateBundle.ElementStates); i++ {
			if worldApp.elementStateBundle.ElementStates[i].State != mashupsdk.Rest {
				switch worldApp.elementStateBundle.ElementStates[i].Id {
				case 1:
					// Inside
					//					worldApp.mainWin.Gls().ClearColor(.545, 0, 0, 1.0)
					worldApp.mainWin.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
				case 2:
					// Outside
					worldApp.mainWin.Gls().ClearColor(.545, 0, 0, 1.0)
				case 3:
					// It
					worldApp.mainWin.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
				case 4:
					// Up-side-down
					worldApp.mainWin.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
				}
				break
			}
		}
		worldApp.mainWin.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(w.scene, w.cam)
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
	worldApp.DetailedElements = detailedElementBundle.DetailedElements
	worldApp.elementStateBundle = mashupsdk.MashupElementStateBundle{
		ElementStates: []*mashupsdk.MashupElementState{},
	}

	for _, detailedElement := range detailedElementBundle.DetailedElements {
		detailedElement.State.State = mashupsdk.Rest
		es := &mashupsdk.MashupElementState{
			Id:    detailedElement.Id,
			State: mashupsdk.Rest,
		}

		worldApp.elementStateBundle.ElementStates = append(worldApp.elementStateBundle.ElementStates, es)

		worldApp.elementDictionary[detailedElement.GetName()] = detailedElement.Id
		worldApp.elementIndex[detailedElement.Id] = es
	}

	log.Printf("G3n UpsertMashupElements updated\n")
	return &worldApp.elementStateBundle, nil
}

func (mSdk *mashupSdkApiHandler) UpsertMashupElementsState(elementStateBundle *mashupsdk.MashupElementStateBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("G3n UpsertMashupElementsState called\n")

	for _, wes := range worldApp.elementIndex {
		wes.State = mashupsdk.Rest
	}

	for _, es := range elementStateBundle.ElementStates {
		if worldApp.elementIndex[es.GetId()].State != es.State {
			worldApp.elementIndex[es.GetId()].State = es.State
		}
	}
	worldApp.mainWin.Dispatch(gui.OnFocus, nil)
	log.Printf("G3n End UpsertMashupElementsState called\n")
	return &mashupsdk.MashupElementStateBundle{}, nil
}

func main() {
	callerCreds := flag.String("CREDS", "", "Credentials of caller")
	insecure := flag.Bool("insecure", false, "Skip server validation")
	flag.Parse()
	worldLog, err := os.OpenFile("world.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(worldLog)

	worldApp = WorldApp{
		mSdkApiHandler:      &mashupSdkApiHandler{},
		elementIndex:        map[int64]*mashupsdk.MashupElementState{},
		elementDictionary:   map[string]int64{},
		displaySetupChan:    make(chan *mashupsdk.MashupDisplayHint, 1),
		displayPositionChan: make(chan *mashupsdk.MashupDisplayHint, 1),
	}

	if *callerCreds != "" {
		server.InitServer(*callerCreds, *insecure, worldApp.mSdkApiHandler, worldApp.wClientInitHandler)
	} else {
		go func() {
			worldApp.displaySetupChan <- &mashupsdk.MashupDisplayHint{Xpos: 0, Ypos: 0, Width: 400, Height: 800}
		}()
	}

	// Initialize the main window.
	go worldApp.InitMainWindow()

	<-worldCompleteChan
}
