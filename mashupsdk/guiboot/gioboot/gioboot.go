//go:build gioboot
// +build gioboot

package gioboot

import (
	"gioui.org/app"
	"log"
)

func InitMainWindow(initHandler interface{}, runtimeHandler interface{}) {
	log.Printf("Gio Sdk Boot init")

	gioInit := initHandler.(func())
	// This is a handler routine.
	gioInit()
	gioRuntimeHandler := runtimeHandler.(func())

	// Run the application -- this will not return.
	go gioRuntimeHandler()

	app.Main()
}
