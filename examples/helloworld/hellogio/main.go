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
	"gioui.org/widget/material"
)

type HelloContext struct {
	MashContext *mashupsdk.MashupContext
}

var helloContext *HelloContext

func main() {
	insecure := flag.Bool("insecure", false, "Skip server validation")
	flag.Parse()

	helloLog, err := os.OpenFile("hellogio.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(helloLog)
	go func() {
		helloContext = &HelloContext{client.BootstrapInit("./world", nil, nil, insecure)}
	}()

	go func() {
		w := app.NewWindow()

		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()

}

func run(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	var imgSize *image.Point
	for {
		e := <-w.Events()
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
			if imgSize == nil {
				if helloContext != nil && helloContext.MashContext != nil {
					helloContext.MashContext.Client.OnResize(helloContext.MashContext,
						&mashupsdk.MashupDisplayBundle{
							AuthToken: client.GetServerAuthToken(),
							MashupDisplayHint: &mashupsdk.MashupDisplayHint{
								Xpos:   int64(0),
								Ypos:   int64(0),
								Width:  int64(e.Size.X),
								Height: int64(e.Size.Y),
							},
						})
				}
				imgSize = &e.Size
			} else {
				if (*imgSize).X != e.Size.X || (*imgSize).Y != e.Size.Y {
					if helloContext != nil && helloContext.MashContext != nil {
						helloContext.MashContext.Client.OnResize(helloContext.MashContext,
							&mashupsdk.MashupDisplayBundle{
								AuthToken: client.GetServerAuthToken(),
								MashupDisplayHint: &mashupsdk.MashupDisplayHint{
									Xpos:   int64(0),
									Ypos:   int64(0),
									Width:  int64(e.Size.X),
									Height: int64(e.Size.Y),
								},
							})
					}
				}
			}

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
