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

type HelloApp struct {
	HelloContext   *HelloContext
	mainWin        *app.Window
	mainWinDisplay *mashupsdk.MashupDisplayHint
	settled        int
	yOffset        int
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
		ha.mainWinDisplay.Ypos = displayHint.Ypos
		resize = true
	}
	if displayHint.Width != 0 && (*ha.mainWinDisplay).Width != displayHint.Width {
		ha.mainWinDisplay.Width = displayHint.Width
		resize = true
	}
	if displayHint.Height != 0 && (*ha.mainWinDisplay).Height != displayHint.Height+int64(ha.yOffset) {
		ha.mainWinDisplay.Height = displayHint.Height + int64(ha.yOffset)
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

	helloApp := HelloApp{}

	go func() {
		helloApp.HelloContext = &HelloContext{client.BootstrapInit("worldg3n", nil, nil, insecure)}
		helloApp.settled |= 8
		helloApp.OnResize(helloApp.mainWinDisplay)
	}()

	// Sync initialization.
	initHandler := func() {
		options := []app.Option{
			app.Size(unit.Dp(800), unit.Dp(600)),
			app.Title("Hello"),
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

				if helloApp.settled > 15 || (ce.Center && (helloApp.settled&1) == 0) {
					if ce.YOffset != 0 {
						helloApp.yOffset = ce.YOffset + 3
					}
					helloApp.settled |= 1
					helloApp.OnResize(&mashupsdk.MashupDisplayHint{
						Xpos:   int64(ce.Position.X),
						Ypos:   int64(ce.Position.Y),
						Width:  int64(ce.Size.X),
						Height: int64(ce.Size.Y),
					})
				}
			case system.PositionEvent:
				if e.YOffset != 0 {
					helloApp.yOffset = e.YOffset + 3
				}
				helloApp.settled |= 2
				helloApp.OnResize(&mashupsdk.MashupDisplayHint{
					Xpos:   int64(e.X),
					Ypos:   int64(e.Y),
					Width:  int64(e.Width),
					Height: int64(e.Height),
				})

			case app.X11ViewEvent:
				// display := e.Display
				// spew.Dump(display)

			case system.StageEvent:
				//stage := e.Stage
				//spew.Dump(stage)

			case key.FocusEvent:
				//fe := e.Focus
				//spew.Dump(fe)

			case pointer.Event:
				// Position of like a cursor.
				//pos := e.Position
				//spew.Dump(pos)

			case system.DestroyEvent:
				helloApp.HelloContext.MashContext.Client.Shutdown(helloApp.HelloContext.MashContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})
				return

			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				helloApp.settled |= 4
				helloApp.OnResize(&mashupsdk.MashupDisplayHint{
					Xpos:   int64(helloApp.mainWinDisplay.Xpos),
					Ypos:   int64(helloApp.mainWinDisplay.Ypos),
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
