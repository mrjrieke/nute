//go:build gomobileboot
// +build gomobileboot

package gomobileboot

import (
	"golang.org/x/mobile/app"
)

func InitMainWindow(initHandler interface{}, runtimeHandler interface{}) interface{} {
	var mainWin interface{}
	var winCompleteChan chan bool

	gomobileRuntimeHandler := runtimeHandler.(func(app.App))

	app.Main(func(a app.App) {
		// ok so long as there is only one call to InitMainWindow
		gomobileRuntimeHandler(a)
		mainWin = &a

		// App starting... ok to let main exit.
		winCompleteChan <- true
	})

	<-winCompleteChan
	return mainWin
}
