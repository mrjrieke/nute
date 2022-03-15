package main

import (
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/client"
)

type HudContext struct {
	MshContext *mashupsdk.MashupContext
}

var hudContext HudContext

func main() {
	a := app.New()
	a.Lifecycle().SetOnEnteredForeground(func() {
		hudContext = HudContext{client.BootstrapInit("./nute", nil, nil)}
	})
	w := a.NewWindow("Hello World")
	w.SetContent(widget.NewLabel("Nute Hud."))
	w.SetCloseIntercept(func() {
		hudContext.MshContext.Client.Shutdown(hudContext.MshContext, &mashupsdk.MashupEmpty{})
		os.Exit(0)
	})
	w.ShowAndRun()

}
