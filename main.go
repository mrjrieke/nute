//go:build darwin || linux
// +build darwin linux

package main

// Nute is a basic gomobile app.
import (
	"flag"
	"log"
	"os"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
	"tini.com/nute/mashupsdk/server"
)

var (
	engine sprite.Engine
)

// CREDS={\"callerToken\":\"742bc42264f857dc68331cc5c26d0f89474fb499a17ac35d8c84cf8491906b54\",\"port\":46733}
func main() {
	callerCreds := flag.String("CREDS", "", "Credentials of caller")
	insecure := flag.Bool("insecure", false, "Skip server validation")
	flag.Parse()
	nuteLog, err := os.OpenFile("nute.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(nuteLog)

	server.InitServer(*callerCreds, *insecure)
	app.Main(func(a app.App) {
		var glCtx gl.Context
		var szEvent size.Event

		for event := range a.Events() {
			switch filteredEvent := a.Filter(event).(type) {
			case lifecycle.Event:
				switch filteredEvent.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glCtx, _ = filteredEvent.DrawContext.(gl.Context)
					onStart(glCtx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
				}
			case size.Event:
				// Capture screen sizing.
				szEvent = filteredEvent
			case paint.Event:
				// Do some painting.
				onPaint(glCtx, szEvent)
			case touch.Event:
				// Capture a screen touched event.
			case key.Event:
				// Capture general keyboard events.

			}

		}

	})

}

func onStart(glCtx gl.Context) {
	engine = glsprite.Engine(nil)
}

func onPaint(glCtx gl.Context, szEvent size.Event) {
	glCtx.ClearColor(1, 1, 1, 1)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
}
