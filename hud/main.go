package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"tini.com/nute/mashupsdk/client"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	go client.BootstrapInit("./nute", nil, nil)

	w.SetContent(widget.NewLabel("Nute Hud."))
	w.ShowAndRun()
}