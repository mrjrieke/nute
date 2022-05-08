//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"errors"
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
	"tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/guiboot"
	"tini.com/nute/mashupsdk/server"
)

var worldCompleteChan chan bool

type worldApiHandler struct {
}

type worldClientInitHandler struct {
}

type WorldApp struct {
	wApiHandler         *worldApiHandler
	wClientInitHandler  *worldClientInitHandler
	displaySetupChan    chan *mashupsdk.MashupDisplayHint
	displayPositionChan chan *mashupsdk.MashupDisplayHint
	mainWin             *app.Application
	scene               *core.Node
	cam                 *camera.Camera

	mashupContext  *mashupsdk.MashupContext // Needed for callbacks to other mashups
	elementIndex   map[string]*mashupsdk.MashupDetailedElement
	Elements       []*mashupsdk.MashupDetailedElement
	ElementsStates mashupsdk.MashupElementStateBundle
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
			if n, intersections := w.Cast(ichild, caster); n != nil {
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
		}
		a.Subscribe(window.OnWindowSize, onResize)
		onResize("", nil)

		// Create a blue torus and add it to the scene
		geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
		mat := material.NewStandard(math32.NewColor("DarkBlue"))
		mesh := graphic.NewMesh(geom, mat)
		mesh.SetLoaderID("torus")
		w.scene.Add(mesh)

		// Create and add a button to the scene
		btn := gui.NewButton("Make Red")
		btn.SetPosition(100, 40)
		btn.SetSize(40, 40)
		btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
			mat.SetColor(math32.NewColor("DarkRed"))
		})
		w.scene.Add(btn)

		w.mainWin.Subscribe(gui.OnMouseUp, func(name string, ev interface{}) {
			mev := ev.(*window.MouseEvent)
			g3Width, g3Height := worldApp.mainWin.GetSize()

			xPosNdc := 2*(mev.Xpos/float32(g3Width)) - 1
			yPosNdc := -2*(mev.Ypos/float32(g3Height)) + 1
			caster := collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
			caster.SetFromCamera(worldApp.cam, xPosNdc, yPosNdc)

			if worldApp.scene.Visible() {
				_, intersections := worldApp.Cast(worldApp.scene, caster)
				if len(intersections) != 0 {
					// TODO: Interact!
					//me := worldApp.elementIndex[n.GetNode().LoaderID()]
					elementStateBundle := mashupsdk.MashupElementStateBundle{}

					// TODO: Make a bundle...
					worldApp.mashupContext.Client.UpsertMashupElementsState(worldApp.mashupContext, &elementStateBundle, nil)
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
		worldApp.mainWin.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(w.scene, w.cam)
	}

	guiboot.InitMainWindow(guiboot.G3n, initHandler, runtimeHandler)
}

func (w *worldClientInitHandler) RegisterContext(context *mashupsdk.MashupContext) {
	worldApp.mashupContext = context
}

func (w *worldApiHandler) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	if worldApp.mainWin != nil && (*worldApp.mainWin).IWindow != nil {
		log.Printf("G3n Received onResize xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
		worldApp.displayPositionChan <- displayHint
	} else {
		log.Printf("G3n Could not apply xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
		worldApp.displaySetupChan <- displayHint
		worldApp.displayPositionChan <- displayHint
	}
}

func (w *worldApiHandler) UpsertMashupDetailedElements(detailedElementBundle *mashupsdk.MashupDetailedElementBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("G3n Received UpsertMashupDetailedElements\n")
	worldApp.Elements = detailedElementBundle.Mashobjects
	for _, e := range worldApp.Elements {
		worldApp.elementIndex[e.GetName()] = e
	}
	worldApp.ElementsStates = mashupsdk.MashupElementStateBundle{
		Mashobjects: make([]*mashupsdk.MashupElementState, len(worldApp.Elements)),
	}

	for _, detailedElement := range worldApp.Elements {
		detailedElement.State = mashupsdk.Rest
		worldApp.ElementsStates.Mashobjects = append(worldApp.ElementsStates.Mashobjects, &mashupsdk.MashupElementState{
			Id:    detailedElement.Id,
			State: mashupsdk.Rest,
		})
	}

	log.Printf("G3n UpsertMashupSociety updated\n")
	return &worldApp.ElementsStates, nil
}

func (w *worldApiHandler) UpsertMashupElementState(elementStateBundle *mashupsdk.MashupElementStateBundle) (*mashupsdk.MashupElementStateBundle, error) {
	// Not implemented.
	log.Printf("G3n UpsertMashupElementState called\n")
	return nil, errors.New("Could not capture items.")
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
		wApiHandler:         &worldApiHandler{},
		displaySetupChan:    make(chan *mashupsdk.MashupDisplayHint, 1),
		displayPositionChan: make(chan *mashupsdk.MashupDisplayHint, 1),
	}

	if *callerCreds != "" {
		server.InitServer(*callerCreds, *insecure, worldApp.wApiHandler, worldApp.wClientInitHandler)
	} else {
		go func() {
			worldApp.displaySetupChan <- &mashupsdk.MashupDisplayHint{Xpos: 0, Ypos: 0, Width: 400, Height: 800}
		}()
	}

	// Initialize the main window.
	go worldApp.InitMainWindow()

	<-worldCompleteChan
}
