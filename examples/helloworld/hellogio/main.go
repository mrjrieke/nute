package main

import (
	"flag"
	"image/color"
	"log"
	"os"

	"tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/client"
	"tini.com/nute/mashupsdk/guiboot"

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
	MashContext *mashupsdk.MashupContext
}
type gioMashupApiHandler struct {
}

type HelloApp struct {
	gioMashupApiHandler  *gioMashupApiHandler
	HelloContext         *HelloContext
	mainWin              *app.Window
	mashupDisplayContext *mashupsdk.MashupDisplayContext
	DetailedElements     []*mashupsdk.MashupDetailedElement
	elementIndex         map[int64]*mashupsdk.MashupElementState // g3n indexes by string...
	elementStateBundle   *mashupsdk.MashupElementStateBundle
}

var helloApp HelloApp

func (ha *HelloApp) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	resize := ha.mashupDisplayContext.OnResize(displayHint)

	if resize {
		if ha.HelloContext == nil || ha.HelloContext.MashContext == nil {
			return
		}

		if ha.HelloContext.MashContext != nil {
			ha.HelloContext.MashContext.Client.OnResize(ha.HelloContext.MashContext,
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
	insecure := flag.Bool("insecure", false, "Skip server validation")
	flag.Parse()

	helloLog, err := os.OpenFile("hellogio.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(helloLog)

	helloApp = HelloApp{
		mashupDisplayContext: &mashupsdk.MashupDisplayContext{MainWinDisplay: &mashupsdk.MashupDisplayHint{}},
		DetailedElements: []*mashupsdk.MashupDetailedElement{
			&mashupsdk.MashupDetailedElement{
				Id:          1,
				State:       &mashupsdk.MashupElementState{Id: 1, State: mashupsdk.Init},
				Name:        "Inside",
				Description: "",
				Genre:       "Space",
				Subgenre:    "Ento",
				Parentids:   nil,
				Childids:    nil,
			},
			&mashupsdk.MashupDetailedElement{
				Id:          2,
				State:       &mashupsdk.MashupElementState{Id: 2, State: mashupsdk.Init},
				Name:        "Outside",
				Description: "",
				Genre:       "Space",
				Subgenre:    "Exo",
				Parentids:   nil,
				Childids:    nil,
			},
			&mashupsdk.MashupDetailedElement{
				Id:          3,
				State:       &mashupsdk.MashupElementState{Id: 3, State: mashupsdk.Init},
				Name:        "torus",
				Description: "",
				Genre:       "Solid",
				Subgenre:    "Ento",
				Parentids:   nil,
				Childids:    []int64{4},
			},
			&mashupsdk.MashupDetailedElement{
				Id:          4,
				State:       &mashupsdk.MashupElementState{Id: 4, State: mashupsdk.Init},
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

	go func() {
		helloApp.HelloContext = &HelloContext{client.BootstrapInit("worldg3n", helloApp.gioMashupApiHandler, nil, nil, insecure)}
		helloApp.mashupDisplayContext.ApplySettled(mashupsdk.AppInitted, false)
		helloApp.OnResize(helloApp.mashupDisplayContext.MainWinDisplay)
	}()

	// Sync initialization.
	initHandler := func() {
		options := []app.Option{
			app.Size(unit.Dp(800), unit.Dp(100)),
			app.Title("Hello Gio World"),
			Center(),
		}
		helloApp.mainWin = app.NewWindow(options...)
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
					helloApp.OnResize(&mashupsdk.MashupDisplayHint{
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
				helloApp.OnResize(&mashupsdk.MashupDisplayHint{
					Xpos:   int64(e.X),
					Ypos:   int64(e.Y),
					Width:  int64(e.Width),
					Height: int64(e.Height),
				})

			case app.X11ViewEvent:
				// display := e.Display

			case system.StageEvent:
				//stage := e.Stage

			case key.FocusEvent:
				//fe := e.Focus

			case pointer.Event:
				// Position of like a cursor.
				//pos := e.Position

			case system.DestroyEvent:
				helloApp.HelloContext.MashContext.Client.Shutdown(helloApp.HelloContext.MashContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})
				os.Exit(0)
				return

			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				helloApp.mashupDisplayContext.ApplySettled(mashupsdk.Frame, false)
				helloApp.OnResize(&mashupsdk.MashupDisplayHint{
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

func (mSdk *gioMashupApiHandler) OnResize(displayHint *mashupsdk.MashupDisplayHint) {
	if helloApp.mainWin != nil {
		log.Printf("Gio Received onResize xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	} else {
		log.Printf("Gio Could not apply xpos: %d ypos: %d width: %d height: %d ytranslate: %d\n", int(displayHint.Xpos), int(displayHint.Ypos), int(displayHint.Width), int(displayHint.Height), int(displayHint.Ypos+displayHint.Height))
	}
}

func (mSdk *gioMashupApiHandler) UpsertMashupElements(detailedElementBundle *mashupsdk.MashupDetailedElementBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("Gio UpsertMashupElements - not implemented\n")
	return nil, nil
}

func (mSdk *gioMashupApiHandler) UpsertMashupElementsState(elementStateBundle *mashupsdk.MashupElementStateBundle) (*mashupsdk.MashupElementStateBundle, error) {
	log.Printf("Gio UpsertMashupElementsState called\n")
	for _, es := range elementStateBundle.ElementStates {
		helloApp.elementIndex[es.GetId()].State = es.State
	}
	log.Printf("Gio UpsertMashupElementsState complete\n")
	return nil, nil
}
