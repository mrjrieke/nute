//go:build fyneboot
// +build fyneboot

package fyneboot

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/widget"
)

func InitMainWindow(initHandler interface{}, runtimeHandler interface{}) {
	a := app.New()

	fyneInitHandler := initHandler.(func(fyne.App))
	fyneInitHandler(a)

	fyneRuntimeHandler := runtimeHandler.(func())
	fyneRuntimeHandler()
}
