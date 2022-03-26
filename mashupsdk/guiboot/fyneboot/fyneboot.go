//go:build fyneboot
// +build fyneboot

package fyneboot

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func InitMainWindow(initHandler interface{}, runtimeHandler interface{}) interface{} {
	a := app.New()

	fyneInitHandler := initHandler.(func(app.App) fyne.Window)
	w := fyneInitHandler(app)

	w.ShowAndRun()

}
