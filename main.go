// +build darwin linux

package main

// Nute is a basic gomobile app.
import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
	"nute/mashupsdk/server"
)

var (
	engine sprite.Engine
)

func main() {
	server.InitServer()
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
