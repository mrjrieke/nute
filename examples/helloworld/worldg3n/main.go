//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"flag"
	"image"
	"log"
	"os"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/window"
	sdk "tini.com/nute/mashupsdk"
	"tini.com/nute/mashupsdk/guiboot"
	"tini.com/nute/mashupsdk/server"
)

var worldCompleteChan chan bool

type worldApiHandler struct {
}

type WorldApp struct {
	wApiHandler *worldApiHandler
	mainWin     *app.Application
	mainWinDims *image.Point
}

var worldApp WorldApp

func (w *WorldApp) InitMainWindow() {
	worldApp.mainWin = guiboot.InitMainWindow(guiboot.G3n, nil, nil).(*app.Application)
}

func (w *worldApiHandler) OnResize(displayHint *sdk.MashupDisplayHint) {
	log.Printf("Received onResize.")
	if worldApp.mainWin == nil {
		worldApp.InitMainWindow()
		(*worldApp.mainWin).IWindow.(*window.GlfwWindow).Window.SetPos(int(displayHint.Xpos), int(displayHint.Ypos+displayHint.Height))
	} else {
		// TODO: This will probably crash
		//(*worldApp.mainWin).IWindow.(*window.GlfwWindow).Window.SetPos(int(displayHint.Xpos), int(displayHint.Ypos+displayHint.Height))
	}
}

func main() {
	callerCreds := flag.String("CREDS", "", "Credentials of caller")
	insecure := flag.Bool("insecure", false, "Skip server validation")
	flag.Parse()
	worldLog, err := os.OpenFile("world.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(worldLog)

	worldApp = WorldApp{wApiHandler: &worldApiHandler{}}

	server.InitServer(*callerCreds, *insecure, worldApp.wApiHandler)

	<-worldCompleteChan
}
