package main

import (
	"flag"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/client"
)

type HelloContext struct {
	MshContext *mashupsdk.MashupContext
}

var helloContext HelloContext

func main() {
	insecure := flag.Bool("insecure", false, "Skip server validation")
	flag.Parse()

	helloLog, err := os.OpenFile("hello.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(helloLog)

	a := app.New()
	a.Lifecycle().SetOnEnteredForeground(func() {
		helloContext = HelloContext{client.BootstrapInit("./world", nil, nil, insecure)}
	})
	a.Lifecycle().SetOnResized(func(xpos int, ypos int, width int, height int) {
		// TODO: Notification to world of resize.
	})

	w := a.NewWindow("Hello World")
	w.Resize(fyne.NewSize(800, 30))
	w.SetContent(widget.NewLabel("The world of hello"))
	w.SetCloseIntercept(func() {
		helloContext.MshContext.Client.Shutdown(helloContext.MshContext, &mashupsdk.MashupEmpty{AuthToken: "TODO"})
		os.Exit(0)
	})
	w.ShowAndRun()
}
