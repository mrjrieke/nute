package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/client"

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
	HelloContext *HelloContext
	mainWin      *app.Window
	mainWinDims  *image.Point
}

func (ha *HelloApp) OnResize(frameEvent *system.FrameEvent) {
	resize := false
	if ha.mainWinDims == nil {
		resize = true
		ha.mainWinDims = &frameEvent.Size
	} else {
		if (*ha.mainWinDims).X != frameEvent.Size.X || (*ha.mainWinDims).Y != frameEvent.Size.Y {
			resize = true
		}
	}
	if resize {
		if ha.HelloContext == nil || ha.HelloContext.MashContext == nil {
			return
		}

		if ha.HelloContext.MashContext != nil {
			ha.HelloContext.MashContext.Client.OnResize(ha.HelloContext.MashContext,
				&mashupsdk.MashupDisplayBundle{
					AuthToken: client.GetServerAuthToken(),
					MashupDisplayHint: &mashupsdk.MashupDisplayHint{
						Xpos:   int64(0),
						Ypos:   int64(0),
						Width:  int64(frameEvent.Size.X),
						Height: int64(frameEvent.Size.Y),
					},
				})
		}
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
	}()

	go func() {
		options := []app.Option{
			app.Size(unit.Dp(800), unit.Dp(600)),
			app.Title("Hello World"),
		}
		helloApp.mainWin = app.NewWindow(options...)
		helloApp.mainWin.Center()

		err := run(&helloApp)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()

}

func run(appContext *HelloApp) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	for {
		e := <-appContext.mainWin.Events()
		// Event handler for main window.
		switch e := e.(type) {
		case app.ConfigEvent:
			ce := e.Config
			spew.Dump(ce)

		case app.X11ViewEvent:
			display := e.Display
			spew.Dump(display)

		case system.StageEvent:
			stage := e.Stage
			spew.Dump(stage)

		case key.FocusEvent:
			fe := e.Focus
			spew.Dump(fe)

		case pointer.Event:
			pos := e.Position
			spew.Dump(pos)

		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			appContext.OnResize(&e)

			title := material.H1(th, "Hello, Gio")
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon
			title.Alignment = text.Middle
			title.Layout(gtx)

			e.Frame(gtx.Ops)
		default:
			fmt.Println("In here.")
		}
	}
}
