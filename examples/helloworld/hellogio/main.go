package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mrjrieke/nute-core/mashupsdk"
	"github.com/mrjrieke/nute/mashupsdk/client"
	"github.com/mrjrieke/nute/mashupsdk/guiboot"
	"google.golang.org/protobuf/types/known/emptypb"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type HelloContext struct {
	mashupContext *mashupsdk.MashupContext // Needed for callbacks to other mashups
}
type gioMashupApiHandler struct {
}

type GioWidgetBundle struct {
	mashupsdk.GuiWidgetBundle
}

type HelloApp struct {
	gioMashupApiHandler  *gioMashupApiHandler
	HelloContext         *HelloContext
	mainWin              *app.Window
	mashupDisplayContext *mashupsdk.MashupDisplayContext
	gioWidgetElements    []*GioWidgetBundle
	gioComponentCache    map[int64]*GioWidgetBundle // g3n indexes by string...

}

var helloApp HelloApp

func (ha *HelloApp) OnDisplayChange(displayHint *mashupsdk.MashupDisplayHint) {
	resize := ha.mashupDisplayContext.OnDisplayChange(displayHint)

	if resize {
		if ha.HelloContext.mashupContext == nil {
			return
		}

		if ha.HelloContext.mashupContext != nil {
			ha.HelloContext.mashupContext.Client.OnDisplayChange(ha.HelloContext.mashupContext,
				&mashupsdk.MashupDisplayBundle{
					AuthToken:         client.GetServerAuthToken(),
					MashupDisplayHint: ha.mashupDisplayContext.MainWinDisplay,
				})
		}
	}
}

func Center() app.Option {
	return func(_ unit.Metric, cnf *app.Config) {
		cnf.Center = true
	}
}

func main() {
	insecure := flag.Bool("tls-skip-validation", false, "Skip server validation")
	flag.Parse()

	helloLog, err := os.OpenFile("hellogio.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(helloLog)

	helloApp = HelloApp{
		gioMashupApiHandler: &gioMashupApiHandler{},
		HelloContext:        &HelloContext{},
		mainWin: app.NewWindow([]app.Option{
			app.Size(unit.Dp(800), unit.Dp(100)),
			app.Title("Hello Gio World"),
			Center(),
		}...),
		mashupDisplayContext: &mashupsdk.MashupDisplayContext{MainWinDisplay: &mashupsdk.MashupDisplayHint{}},
		gioWidgetElements: []*GioWidgetBundle{
			{
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent: container.NewTabItem("Inside", widget.NewLabel("The magnetic field inside a toroid is always tangential to the circular closed path.  These magnetic field lines are concentric circles.")),
					MashupDetailedElement: &mashupsdk.MashupDetailedElement{
						Id:          1,
						State:       &mashupsdk.MashupElementState{Id: 1, State: int64(mashupsdk.Init)},
						Name:        "Inside",
						Description: "",
						Genre:       "Space",
						Subgenre:    "Ento",
						Parentids:   []int64{3},
						Childids:    nil,
					},
				},
			},
			{
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent: container.NewTabItem("Outside", widget.NewLabel("The magnetic field at any point outside the toroid is zero.")),
					MashupDetailedElement: &mashupsdk.MashupDetailedElement{
						Id:          2,
						State:       &mashupsdk.MashupElementState{Id: 2, State: int64(mashupsdk.Init)},
						Name:        "Outside",
						Description: "",
						Genre:       "Space",
						Subgenre:    "Exo",
						Parentids:   nil,
						Childids:    nil,
					},
				},
			},
			{
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent: container.NewTabItem("It", widget.NewLabel("The magnetic field inside the empty space surrounded by the toroid is zero.")),
					MashupDetailedElement: &mashupsdk.MashupDetailedElement{
						Id:          3,
						State:       &mashupsdk.MashupElementState{Id: 3, State: int64(mashupsdk.Init)},
						Name:        "torus",
						Description: "",
						Genre:       "Solid",
						Subgenre:    "Ento",
						Parentids:   nil,
						Childids:    []int64{4},
					},
				},
			},
			{
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent: container.NewTabItem("Up-side-down", widget.NewLabel("Torus is up-side-down")),
					MashupDetailedElement: &mashupsdk.MashupDetailedElement{
						Id:          4,
						State:       &mashupsdk.MashupElementState{Id: 4, State: int64(mashupsdk.Init)},
						Name:        "Up-Side-Down",
						Description: "",
						Genre:       "Attitude",
						Subgenre:    "",
						Parentids:   []int64{3},
						Childids:    nil,
					},
				},
			},
			{
				GuiWidgetBundle: mashupsdk.GuiWidgetBundle{
					GuiComponent: container.NewTabItem("Hide", widget.NewLabel("Poof...")),
					MashupDetailedElement: &mashupsdk.MashupDetailedElement{
						Id:          5,
						State:       &mashupsdk.MashupElementState{Id: 5, State: int64(mashupsdk.Hidden)},
						Name:        "Hide",
						Description: "",
						Genre:       "Hidden",
						Subgenre:    "",
						Parentids:   []int64{3},
						Childids:    nil,
					},
				},
			},
		},
		gioComponentCache: map[int64]*GioWidgetBundle{},
	}

	// Build G3nDetailedElement cache.
	for _, gc := range helloApp.gioWidgetElements {
		helloApp.gioComponentCache[gc.MashupDetailedElement.Id] = gc
	}

	go func() {
		helloApp.HelloContext.mashupContext = client.BootstrapInit("worldg3n", helloApp.gioMashupApiHandler, nil, nil, insecure)
		helloApp.mashupDisplayContext.ApplySettled(mashupsdk.AppInitted, false)
		helloApp.OnDisplayChange(helloApp.mashupDisplayContext.MainWinDisplay)

		DetailedElements := []*mashupsdk.MashupDetailedElement{}
		for _, fyneComponent := range helloApp.gioComponentCache {
			DetailedElements = append(DetailedElements, fyneComponent.MashupDetailedElement)
		}
		log.Printf("Delivering mashup elements: %d\n", len(DetailedElements))

		var upsertErr error

		// Connection with mashup fully established.  Initialize mashup elements.
		_, upsertErr = helloApp.HelloContext.mashupContext.Client.UpsertElements(helloApp.HelloContext.mashupContext,
			&mashupsdk.MashupDetailedElementBundle{
				AuthToken:        client.GetServerAuthToken(),
				DetailedElements: DetailedElements,
			})

		if upsertErr != nil {
			log.Printf("Element state initialization failure: %s\n", upsertErr.Error())
		}
		log.Printf("Mashup elements delivered.\n")

	}()

	// Sync initialization.
	initHandler := func() {
		helloApp.mainWin.Center()
	}

	// Async handler.
	runtimeHandler := func() {
		th := material.NewTheme(gofont.Collection())
		var ops op.Ops
		for {
			e := <-helloApp.mainWin.Events()
			// Event handler for main window.
			switch e := e.(type) {
			case app.ConfigEvent:
				ce := e.Config

				if helloApp.mashupDisplayContext.GetSettled() > 15 || (ce.Center && (helloApp.mashupDisplayContext.GetSettled()&1) == 0) {
					if ce.YOffset != 0 {
						helloApp.mashupDisplayContext.SetYoffset(ce.YOffset + 3)
					}
					helloApp.mashupDisplayContext.ApplySettled(mashupsdk.Configured, false)
					helloApp.OnDisplayChange(&mashupsdk.MashupDisplayHint{
						Xpos:   int64(ce.Position.X),
						Ypos:   int64(ce.Position.Y),
						Width:  int64(ce.Size.X),
						Height: int64(ce.Size.Y),
					})
				}
			case system.PositionEvent:
				if e.YOffset != 0 {
					helloApp.mashupDisplayContext.SetYoffset(e.YOffset + 3)
				}
				helloApp.mashupDisplayContext.ApplySettled(mashupsdk.Position, false)
				helloApp.OnDisplayChange(&mashupsdk.MashupDisplayHint{
					Xpos:   int64(e.X),
					Ypos:   int64(e.Y),
					Width:  int64(e.Width),
					Height: int64(e.Height),
				})

				//			case app.X11ViewEvent:
				// display := e.Display

			case system.StageEvent:
				//stage := e.Stage

			case key.FocusEvent:
				//fe := e.Focus

			case pointer.Event:
				// Position of like a cursor.
				//pos := e.Position

			case system.DestroyEvent:
				helloApp.HelloContext.mashupContext.Client.Shutdown(helloApp.HelloContext.mashupContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})
				os.Exit(0)
				return

			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				helloApp.mashupDisplayContext.ApplySettled(mashupsdk.Frame, false)
				helloApp.OnDisplayChange(&mashupsdk.MashupDisplayHint{
					Xpos:   int64(helloApp.mashupDisplayContext.MainWinDisplay.Xpos),
					Ypos:   int64(helloApp.mashupDisplayContext.MainWinDisplay.Ypos),
					Width:  int64(e.Size.X),
					Height: int64(e.Size.Y),
				})

				title := material.H1(th, "Hello, Gio")
				maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
				title.Color = maroon
				title.Alignment = text.Middle
				title.Layout(gtx)

				e.Frame(gtx.Ops)
			}

		}
	}

	guiboot.InitMainWindow(guiboot.Gio, initHandler, runtimeHandler)
}

func (mSdk *gioMashupApiHandler) OnDisplayChange(displayHint *mashupsdk.MashupDisplayHint) {
	if helloApp.mainWin != nil {
		log.Printf("Gio Received OnDisplayChange xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	} else {
		log.Printf("Gio Could not apply xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	}
}

func (mSdk *gioMashupApiHandler) GetElements() (*mashupsdk.MashupDetailedElementBundle, error) {
	log.Printf("Gio GetElements - not implemented\n")
	return &mashupsdk.MashupDetailedElementBundle{}, nil
}

func (mSdk *gioMashupApiHandler) UpsertElements(detailedElementBundle *mashupsdk.MashupDetailedElementBundle) (*mashupsdk.MashupDetailedElementBundle, error) {
	log.Printf("Gio UpsertElements - not implemented\n")
	return nil, nil
}

func (mSdk *gioMashupApiHandler) TweakStates(elementStateBundle *mashupsdk.MashupElementStateBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("Gio TweakStates called\n")
	for _, es := range elementStateBundle.ElementStates {
		fyneComponent := helloApp.gioComponentCache[es.GetId()]
		fyneComponent.MashupDetailedElement.State.State = es.State
		if (mashupsdk.DisplayElementState(es.State) & mashupsdk.Clicked) == mashupsdk.Clicked {
			// TODO: Select the item.
		}
	}
	log.Printf("Gio TweakStates complete\n")
	return &mashupsdk.MashupElementStateBundle{}, nil
}

func (mSdk *gioMashupApiHandler) TweakStatesByMotiv(motivIn *mashupsdk.Motiv) (*emptypb.Empty, error) {
	log.Printf("Gio Received TweakStatesByMotiv\n")
	// TODO: Find and TweakStates...
	fmt.Println(motivIn.Code)

	log.Printf("Gio finished TweakStatesByMotiv handle.\n")
	return &emptypb.Empty{}, nil
}

func (mSdk *gioMashupApiHandler) ResetStates() {
	// Not implemented.
}
