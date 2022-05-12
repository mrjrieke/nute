package main

import (
	"embed"
	"flag"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/client"
	"tini.com/nute/mashupsdk/guiboot"
)

type HelloContext struct {
	MashContext *mashupsdk.MashupContext
}

type fyneMashupApiHandler struct {
}

var helloContext HelloContext

type HelloApp struct {
	fyneMashupApiHandler *fyneMashupApiHandler
	HelloContext         *HelloContext
	mainWin              fyne.Window
	mainWinDisplay       *mashupsdk.MashupDisplayHint
	settled              int
	yOffset              int
	DetailedElements     []*mashupsdk.MashupDetailedElement
	elementIndex         map[int64]*mashupsdk.MashupElementState // g3n indexes by string...
	elementStateBundle   *mashupsdk.MashupElementStateBundle
}

func (ha *HelloApp) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	resize := false
	if ha.mainWinDisplay == nil {
		resize = true
		ha.mainWinDisplay = &mashupsdk.MashupDisplayHint{}
	}

	if displayHint == nil {
		return
	}

	if displayHint.Xpos != 0 && (*ha.mainWinDisplay).Xpos != displayHint.Xpos {
		ha.mainWinDisplay.Xpos = displayHint.Xpos
		resize = true
	}
	if displayHint.Ypos != 0 && (*ha.mainWinDisplay).Ypos != displayHint.Ypos {
		ha.mainWinDisplay.Ypos = displayHint.Ypos + int64(ha.yOffset)
		resize = true
	}
	if displayHint.Width != 0 && (*ha.mainWinDisplay).Width != displayHint.Width {
		ha.mainWinDisplay.Width = displayHint.Width
		resize = true
	}
	if displayHint.Height != 0 && (*ha.mainWinDisplay).Height != displayHint.Height+int64(ha.yOffset) {
		ha.mainWinDisplay.Height = displayHint.Height
		resize = true
	}

	if ha.settled < 15 {
		return
	} else if ha.settled == 15 {
		resize = true
		ha.settled = 31
	}

	if resize {
		if ha.HelloContext == nil || ha.HelloContext.MashContext == nil {
			return
		}

		if ha.HelloContext.MashContext != nil {
			ha.HelloContext.MashContext.Client.OnResize(ha.HelloContext.MashContext,
				&mashupsdk.MashupDisplayBundle{
					AuthToken:         client.GetServerAuthToken(),
					MashupDisplayHint: ha.mainWinDisplay,
				})
		}
	}
}

var helloApp HelloApp

//go:embed gophericon.png
var gopherIcon embed.FS

func main() {
	insecure := flag.Bool("insecure", false, "Skip server validation")
	flag.Parse()

	helloLog, err := os.OpenFile("hellofyne.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(helloLog)

	helloApp = HelloApp{
		fyneMashupApiHandler: &fyneMashupApiHandler{},
		DetailedElements: []*mashupsdk.MashupDetailedElement{
			&mashupsdk.MashupDetailedElement{
				Id:          1,
				State:       mashupsdk.Init,
				Name:        "Inside",
				Description: "",
				Genre:       "Space",
				Subgenre:    "Ento",
				Parentids:   nil,
				Childids:    nil,
			},
			&mashupsdk.MashupDetailedElement{
				Id:          2,
				State:       mashupsdk.Init,
				Name:        "Outside",
				Description: "",
				Genre:       "Space",
				Subgenre:    "Exo",
				Parentids:   nil,
				Childids:    nil,
			},
			&mashupsdk.MashupDetailedElement{
				Id:          3,
				State:       mashupsdk.Init,
				Name:        "torus",
				Description: "",
				Genre:       "Solid",
				Subgenre:    "Ento",
				Parentids:   nil,
				Childids:    []int64{4},
			},
			&mashupsdk.MashupDetailedElement{
				Id:          4,
				State:       mashupsdk.Init,
				Name:        "Up-Side-Down",
				Description: "",
				Genre:       "Attitude",
				Subgenre:    "",
				Parentids:   []int64{3},
				Childids:    nil,
			},
		},
		elementStateBundle: &mashupsdk.MashupElementStateBundle{},
		elementIndex:       map[int64]*mashupsdk.MashupElementState{},
	}

	// Sync initialization.
	initHandler := func(a fyne.App) {
		a.Lifecycle().SetOnEnteredForeground(func() {
			helloApp.elementStateBundle.ElementStates = make([]*mashupsdk.MashupElementState, len(helloApp.DetailedElements))

			// Init element states.
			for _, mashupDetailedElement := range helloApp.DetailedElements {
				es := &mashupsdk.MashupElementState{Id: mashupDetailedElement.Id, State: mashupDetailedElement.State}
				helloApp.elementStateBundle.ElementStates = append(helloApp.elementStateBundle.ElementStates, es)
				helloApp.elementIndex[es.GetId()] = es
			}

			if helloApp.HelloContext == nil {
				helloApp.HelloContext = &HelloContext{client.BootstrapInit("worldg3n", helloApp.fyneMashupApiHandler, nil, nil, insecure)}

				var upsertErr error
				// Connection with mashup fully established.  Initialize mashup elements.
				helloApp.elementStateBundle, upsertErr = helloApp.HelloContext.MashContext.Client.UpsertMashupElements(helloApp.HelloContext.MashContext,
					&mashupsdk.MashupDetailedElementBundle{
						AuthToken:        client.GetServerAuthToken(),
						DetailedElements: helloApp.DetailedElements,
					})

				if upsertErr != nil {
					log.Printf("Element state initialization failure: %s\n", upsertErr.Error())
				}

				helloApp.settled |= 8
			}
			helloApp.OnResize(helloApp.mainWinDisplay)
		})
		a.Lifecycle().SetOnResized(func(xpos int, ypos int, yoffset int, width int, height int) {
			log.Printf("Received resize: %d %d %d %d %d\n", xpos, ypos, yoffset, width, height)
			helloApp.settled |= 1
			helloApp.settled |= 2
			helloApp.settled |= 4

			if helloApp.yOffset == 0 {
				helloApp.yOffset = yoffset + 3
			}

			helloApp.OnResize(&mashupsdk.MashupDisplayHint{
				Xpos:   int64(xpos),
				Ypos:   int64(ypos),
				Width:  int64(width),
				Height: int64(height),
			})
		})
		helloApp.mainWin = a.NewWindow("Hello Fyne World")
		gopherIconBytes, _ := gopherIcon.ReadFile("gophericon.png")

		helloApp.mainWin.SetIcon(fyne.NewStaticResource("Gopher", gopherIconBytes))
		helloApp.mainWin.Resize(fyne.NewSize(800, 100))

		torusMenu := container.NewAppTabs(
			container.NewTabItem("Inside", widget.NewLabel("The magnetic field inside a toroid is always tangential to the circular closed path.  These magnetic field lines are concentric circles.")),
			container.NewTabItem("Outside", widget.NewLabel("The magnetic field at any point outside the toroid is zero.")),
			container.NewTabItem("It", widget.NewLabel("The magnetic field inside the empty space surrounded by the toroid is zero.")),
			container.NewTabItem("Up-side-down", widget.NewLabel("Torus is up-side-down")),
		)
		torusMenu.SetTabLocation(container.TabLocationTop)

		helloApp.mainWin.SetContent(torusMenu)
		helloApp.mainWin.SetCloseIntercept(func() {
			helloApp.HelloContext.MashContext.Client.Shutdown(helloApp.HelloContext.MashContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})
			os.Exit(0)
		})
	}

	// Async handler.
	runtimeHandler := func() {
		helloApp.mainWin.ShowAndRun()
	}

	guiboot.InitMainWindow(guiboot.Fyne, initHandler, runtimeHandler)

}

func (mSdk *fyneMashupApiHandler) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	if helloApp.mainWin != nil {
		log.Printf("Fyne Received onResize xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	} else {
		log.Printf("Fyne Could not apply xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	}
}

func (mSdk *fyneMashupApiHandler) UpsertMashupElements(detailedElementBundle *mashupsdk.MashupDetailedElementBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("Fyne UpsertMashupElements - not implemented\n")
	return nil, nil
}

func (mSdk *fyneMashupApiHandler) UpsertMashupElementsState(elementStateBundle *mashupsdk.MashupElementStateBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("Fyne UpsertMashupElementsState called\n")
	for _, es := range elementStateBundle.ElementStates {
		helloApp.elementIndex[es.GetId()].State = es.State
	}
	log.Printf("Fyne UpsertMashupElementsState complete\n")
	return nil, nil
}
