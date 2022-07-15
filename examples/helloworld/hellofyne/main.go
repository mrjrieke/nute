package main

import (
	"embed"
	"flag"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mrjrieke/nute/mashupsdk"
	"github.com/mrjrieke/nute/mashupsdk/client"
	"github.com/mrjrieke/nute/mashupsdk/guiboot"
)

type HelloContext struct {
	mashupContext *mashupsdk.MashupContext // Needed for callbacks to other mashups
}

type fyneMashupApiHandler struct {
}

var helloContext HelloContext

type FyneWidgetBundle struct {
	mashupsdk.GuiWidgetBundle
}

type HelloApp struct {
	fyneMashupApiHandler         *fyneMashupApiHandler
	HelloContext                 *HelloContext
	mainWin                      fyne.Window
	mashupDisplayContext         *mashupsdk.MashupDisplayContext
	mashupDetailedElementLibrary map[int64]*mashupsdk.MashupDetailedElement
	elementLoaderIndex           map[string]int64 // mashup indexes by Name
	fyneWidgetElements           map[string]*FyneWidgetBundle
}

func (fwb *FyneWidgetBundle) OnClicked() {
	fwb.MashupDetailedElement.State.State = int64(mashupsdk.Clicked)

	elementStateBundle := mashupsdk.MashupElementStateBundle{
		AuthToken:     client.GetServerAuthToken(),
		ElementStates: []*mashupsdk.MashupElementState{fwb.MashupDetailedElement.State},
	}
	helloApp.HelloContext.mashupContext.Client.UpsertMashupElementsState(helloApp.HelloContext.mashupContext, &elementStateBundle)
}

func (ha *HelloApp) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	resize := ha.mashupDisplayContext.OnResize(displayHint)

	if resize {
		if ha.HelloContext.mashupContext == nil {
			return
		}

		if ha.HelloContext.mashupContext != nil {
			ha.HelloContext.mashupContext.Client.OnResize(ha.HelloContext.mashupContext,
				&mashupsdk.MashupDisplayBundle{
					AuthToken:         client.GetServerAuthToken(),
					MashupDisplayHint: ha.mashupDisplayContext.MainWinDisplay,
				})
		}
	}
}

func (ha *HelloApp) TorusParser(childId int64) {
	child := helloApp.mashupDetailedElementLibrary[childId]
	if child.Alias != "" {
		helloApp.fyneWidgetElements[child.Alias].MashupDetailedElement = child
	}

	if len(child.GetChildids()) > 0 {
		for _, cId := range child.GetChildids() {
			ha.TorusParser(cId)
		}

	}
}

var helloApp HelloApp

//go:embed gophericon.png
var gopherIcon embed.FS

//go:embed tls/mashup.crt
var mashupCert embed.FS

//go:embed tls/mashup.key
var mashupKey embed.FS

func main() {
	insecure := flag.Bool("insecure", false, "Skip server validation")
	flag.Parse()

	helloLog, err := os.OpenFile("hellofyne.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf(err.Error(), err)
	}
	log.SetOutput(helloLog)

	mashupsdk.InitCertKeyPair(mashupCert, mashupKey)

	helloApp = HelloApp{
		fyneMashupApiHandler:         &fyneMashupApiHandler{},
		HelloContext:                 &HelloContext{},
		mainWin:                      nil,
		mashupDisplayContext:         &mashupsdk.MashupDisplayContext{MainWinDisplay: &mashupsdk.MashupDisplayHint{}},
		mashupDetailedElementLibrary: map[int64]*mashupsdk.MashupDetailedElement{}, // mashupDetailedElementLibrary,
		elementLoaderIndex:           map[string]int64{},                           // elementLoaderIndex
		fyneWidgetElements: map[string]*FyneWidgetBundle{
			"Inside": {
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent:          container.NewTabItem("Inside", widget.NewLabel("The magnetic field inside a toroid is always tangential to the circular closed path.  These magnetic field lines are concentric circles.")),
					MashupDetailedElement: nil, // mashupDetailedElementLibrary["{0}-AxialCircle"],
				},
			},
			"Outside": {
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent:          container.NewTabItem("Outside", widget.NewLabel("The magnetic field at any point outside the toroid is zero.")),
					MashupDetailedElement: nil, //mashupDetailedElementLibrary["Outside"],
				},
			},
			"It": {
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent:          container.NewTabItem("It", widget.NewLabel("The magnetic field inside the empty space surrounded by the toroid is zero.")),
					MashupDetailedElement: nil, //mashupDetailedElementLibrary["{0}-Torus"],
				},
			},
			"Up-Side-Down": {
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent:          container.NewTabItem("Up-Side-Down", widget.NewLabel("Torus is up-side-down")),
					MashupDetailedElement: nil, //mashupDetailedElementLibrary["{0}-SharedAttitude"],
				},
			},
		},
	}

	// Build G3nDetailedElement cache.
	for _, fc := range helloApp.fyneWidgetElements {
		fc.GuiComponent.(*container.TabItem).Content.(*widget.Label).Wrapping = fyne.TextWrapWord
	}

	// Sync initialization.
	initHandler := func(a fyne.App) {
		a.Lifecycle().SetOnEnteredForeground(func() {
			if helloApp.HelloContext.mashupContext == nil {
				helloApp.HelloContext.mashupContext = client.BootstrapInit("worldg3n", helloApp.fyneMashupApiHandler, nil, nil, insecure)

				var upsertErr error
				var concreteElementBundle *mashupsdk.MashupDetailedElementBundle
				DetailedElements := []*mashupsdk.MashupDetailedElement{
					{
						Basisid:       -1,
						State:         &mashupsdk.MashupElementState{Id: -1, State: int64(mashupsdk.Mutable)},
						Name:          "{0}-Torus",
						Alias:         "It",
						Description:   "",
						Renderer:      "Torus",
						Colabrenderer: "Background",
						Genre:         "Solid",
						Subgenre:      "Ento",
						Parentids:     nil,
						Childids:      []int64{-2, 4},
					},
					{
						Basisid:     -2,
						State:       &mashupsdk.MashupElementState{Id: -2, State: int64(mashupsdk.Mutable)},
						Name:        "{0}-AxialCircle",
						Alias:       "Inside",
						Description: "",
						Renderer:    "Torus",
						Genre:       "Space",
						Subgenre:    "Ento",
						Parentids:   []int64{-1},
						Childids:    []int64{4},
					},
					{
						Id:          4,
						State:       &mashupsdk.MashupElementState{Id: 4, State: int64(mashupsdk.Mutable)},
						Name:        "Up-Side-Down",
						Alias:       "Up-Side-Down",
						Description: "",
						Genre:       "Attitude",
						Subgenre:    "180,0,0",
						Parentids:   nil,
						Childids:    nil,
					},
					{
						Id:          5,
						State:       &mashupsdk.MashupElementState{Id: 5, State: int64(mashupsdk.Init)},
						Name:        "ToriOne",
						Description: "Tori",
						Genre:       "Collection",
						Subgenre:    "Torus",
						Parentids:   []int64{},
						Childids:    []int64{7, 8},
					},
					{
						Id:          6,
						State:       &mashupsdk.MashupElementState{Id: 6, State: int64(mashupsdk.Init)},
						Name:        "Outside",
						Alias:       "Outside",
						Description: "",
						Renderer:    "Background",
						Genre:       "Space",
						Subgenre:    "Exo",
						Parentids:   nil,
						Childids:    nil,
					},
					{
						Id:          7,
						State:       &mashupsdk.MashupElementState{Id: 7, State: int64(mashupsdk.Init)},
						Name:        "TorusEntity-One",
						Description: "",
						Genre:       "Abstract",
						Subgenre:    "",
						Parentids:   []int64{5},
						Childids:    []int64{-1}, // -1 -- generated and replaced by server since it is immutable.
					},
					{
						Id:          8,
						State:       &mashupsdk.MashupElementState{Id: 8, State: int64(mashupsdk.Init)},
						Name:        "TorusEntity-Two",
						Description: "",
						Genre:       "Abstract",
						Subgenre:    "",
						Parentids:   []int64{5},
						Childids:    []int64{-1}, // -1 -- generated and replaced by server since it is immutable.
					},
					{
						Id:          9,
						State:       &mashupsdk.MashupElementState{Id: 9, State: int64(mashupsdk.Init)},
						Name:        "TorusEntity-Three",
						Description: "",
						Genre:       "Abstract",
						Subgenre:    "",
						Parentids:   []int64{5},
						Childids:    []int64{-1}, // -1 -- generated and replaced by server since it is immutable.
					},
				}
				for _, detailedElement := range helloApp.mashupDetailedElementLibrary {
					DetailedElements = append(DetailedElements, detailedElement)
				}
				log.Printf("Delivering mashup elements.\n")

				// Connection with mashup fully established.  Initialize mashup elements.
				concreteElementBundle, upsertErr = helloApp.HelloContext.mashupContext.Client.UpsertMashupElements(helloApp.HelloContext.mashupContext,
					&mashupsdk.MashupDetailedElementBundle{
						AuthToken:        client.GetServerAuthToken(),
						DetailedElements: DetailedElements,
					})

				if upsertErr != nil {
					log.Printf("Element state initialization failure: %s\n", upsertErr.Error())
				}

				for _, concreteElement := range concreteElementBundle.DetailedElements {
					//helloApp.fyneComponentCache[generatedComponent.Basisid]
					helloApp.mashupDetailedElementLibrary[concreteElement.Id] = concreteElement
					helloApp.elementLoaderIndex[concreteElement.Name] = concreteElement.Id

					if concreteElement.GetName() == "Outside" {
						helloApp.fyneWidgetElements["Outside"].MashupDetailedElement = concreteElement
					}
				}

				for _, concreteElement := range concreteElementBundle.DetailedElements {
					if strings.HasPrefix(concreteElement.GetName(), "TorusEntity") {
						for _, childId := range concreteElement.Childids {
							helloApp.TorusParser(childId)
						}
					}
				}

				log.Printf("Mashup elements delivered.\n")

				helloApp.mashupDisplayContext.ApplySettled(mashupsdk.AppInitted, false)
			}
			helloApp.OnResize(helloApp.mashupDisplayContext.MainWinDisplay)
		})
		a.Lifecycle().SetOnResized(func(xpos int, ypos int, yoffset int, width int, height int) {
			log.Printf("Received resize: %d %d %d %d %d\n", xpos, ypos, yoffset, width, height)
			helloApp.mashupDisplayContext.ApplySettled(mashupsdk.Configured|mashupsdk.Position|mashupsdk.Frame, false)

			if helloApp.mashupDisplayContext.GetYoffset() == 0 {
				helloApp.mashupDisplayContext.SetYoffset(yoffset + 3)
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
		helloApp.mainWin.SetFixedSize(false)

		torusMenu := container.NewAppTabs(
			helloApp.fyneWidgetElements["Inside"].GuiComponent.(*container.TabItem),       // inside
			helloApp.fyneWidgetElements["Outside"].GuiComponent.(*container.TabItem),      // outside
			helloApp.fyneWidgetElements["It"].GuiComponent.(*container.TabItem),           // It
			helloApp.fyneWidgetElements["Up-Side-Down"].GuiComponent.(*container.TabItem), // Upside down
		)
		torusMenu.OnSelected = func(tabItem *container.TabItem) {
			// Too bad fyne doesn't have the ability for user to assign an id to TabItem...
			// Lookup by name instead and try to keep track of any name changes instead...
			log.Printf("Selected: %s\n", tabItem.Text)
			if mashupItemIndex, miOk := helloApp.elementLoaderIndex[tabItem.Text]; miOk {
				mashupDetailedElement := helloApp.mashupDetailedElementLibrary[mashupItemIndex]
				if mashupDetailedElement.Alias != "" {
					helloApp.fyneWidgetElements[mashupDetailedElement.Alias].OnClicked()
					return
				}
			}
			helloApp.fyneWidgetElements[tabItem.Text].OnClicked()
		}

		torusMenu.SetTabLocation(container.TabLocationTop)

		helloApp.mainWin.SetContent(torusMenu)
		helloApp.mainWin.SetCloseIntercept(func() {
			helloApp.HelloContext.mashupContext.Client.Shutdown(helloApp.HelloContext.mashupContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})
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
	log.Printf("Fyne OnResize - not implemented yet..\n")
	if helloApp.mainWin != nil {
		// TODO: Resize without infinite looping....
		// The moment fyne is resized, it'll want to resize g3n...
		// Which then wants to resize fyne ad-infinitum
		//helloApp.mainWin.PosResize(int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height))
		log.Printf("Fyne Received onResize xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	} else {
		log.Printf("Fyne Could not apply xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	}
}

func (mSdk *fyneMashupApiHandler) UpsertMashupElements(detailedElementBundle *mashupsdk.MashupDetailedElementBundle) (*mashupsdk.MashupDetailedElementBundle, error) {
	log.Printf("Fyne UpsertMashupElements - not implemented\n")
	return &mashupsdk.MashupDetailedElementBundle{}, nil
}

func (mSdk *fyneMashupApiHandler) UpsertMashupElementsState(elementStateBundle *mashupsdk.MashupElementStateBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("Fyne UpsertMashupElementsState called\n")
	for _, es := range elementStateBundle.ElementStates {
		detailedElement := helloApp.mashupDetailedElementLibrary[es.GetId()]

		fyneComponent := helloApp.fyneWidgetElements[detailedElement.GetAlias()]
		fyneComponent.MashupDetailedElement = detailedElement
		fyneComponent.MashupDetailedElement.State.State = es.State

		if mashupsdk.DisplayElementState(es.State) == mashupsdk.Clicked {
			for _, childId := range detailedElement.GetChildids() {
				if childDetailedElement, childDetailOk := helloApp.mashupDetailedElementLibrary[childId]; childDetailOk {
					if childFyneComponent, childFyneOk := helloApp.fyneWidgetElements[childDetailedElement.GetAlias()]; childFyneOk {
						childFyneComponent.MashupDetailedElement = childDetailedElement
						childFyneComponent.GuiComponent.(*container.TabItem).Text = childDetailedElement.Name
					}
				}
			}
			for _, parentId := range detailedElement.GetParentids() {
				if parentDetailedElement, parentDetailOk := helloApp.mashupDetailedElementLibrary[parentId]; parentDetailOk {
					if parentFyneComponent, parentFyneOk := helloApp.fyneWidgetElements[parentDetailedElement.GetAlias()]; parentFyneOk {
						parentFyneComponent.MashupDetailedElement = parentDetailedElement
						parentFyneComponent.GuiComponent.(*container.TabItem).Text = parentDetailedElement.Name
					}
				}
			}

			torusMenu := helloApp.mainWin.Content().(*container.AppTabs)
			// Select the item.
			fyneComponent.GuiComponent.(*container.TabItem).Text = detailedElement.Name
			torusMenu.Select(fyneComponent.GuiComponent.(*container.TabItem))
		}
	}
	log.Printf("Fyne UpsertMashupElementsState complete\n")
	return &mashupsdk.MashupElementStateBundle{}, nil
}
