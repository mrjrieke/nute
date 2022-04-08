package main

import (
	"flag"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/client"
	"tini.com/nute/mashupsdk/guiboot"
)

type HelloContext struct {
	MashContext *mashupsdk.MashupContext
}

var helloContext HelloContext

type HelloApp struct {
	HelloContext   *HelloContext
	mainWin        fyne.Window
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

func main() {
	insecure := flag.Bool("insecure", false, "Skip server validation")
	flag.Parse()

	helloLog, err := os.OpenFile("hellofyne.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(helloLog)

	helloApp := HelloApp{}

	// Sync initialization.
	initHandler := func(a fyne.App) {
		a.Lifecycle().SetOnEnteredForeground(func() {
			if helloApp.HelloContext == nil {
				helloApp.HelloContext = &HelloContext{client.BootstrapInit("worldg3n", nil, nil, insecure)}
				helloApp.settled |= 8
			}
			helloApp.OnResize(helloApp.mainWinDisplay)
		})
		a.Lifecycle().SetOnResized(func(xpos int, ypos int, width int, height int) {
			log.Printf("Received resize: %d %d %d %d\n", xpos, ypos, width, height)
			helloApp.settled |= 1
			helloApp.settled |= 2
			helloApp.settled |= 4

			helloApp.OnResize(&mashupsdk.MashupDisplayHint{
				Xpos:   int64(xpos),
				Ypos:   int64(ypos),
				Width:  int64(width),
				Height: int64(height),
			})
		})
		helloApp.mainWin = a.NewWindow("Hello Fyne World")

		helloApp.mainWin.Resize(fyne.NewSize(800, 30))
		helloApp.mainWin.SetContent(widget.NewLabel("The world of hello"))
		helloApp.mainWin.SetCloseIntercept(func() {
			helloContext.MashContext.Client.Shutdown(helloContext.MashContext, &mashupsdk.MashupEmpty{AuthToken: client.GetServerAuthToken()})
			os.Exit(0)
		})
	}

	// Async handler.
	runtimeHandler := func() {
		helloApp.mainWin.ShowAndRun()
	}

	guiboot.InitMainWindow(guiboot.Fyne, initHandler, runtimeHandler)

}
