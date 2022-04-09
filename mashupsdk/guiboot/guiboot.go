package guiboot

import (
	"tini.com/nute/mashupsdk/guiboot/fyneboot"
	"tini.com/nute/mashupsdk/guiboot/g3nboot"
	"tini.com/nute/mashupsdk/guiboot/gioboot"
	"tini.com/nute/mashupsdk/guiboot/gomobileboot"
)

type GuiProvider int64

const (
	Fyne GuiProvider = iota
	G3n
	Gio
	Gomobile
)

func InitMainWindow(guiType GuiProvider, initHandler interface{}, runtimeHandle interface{}) {

	switch guiType {
	case Fyne:
		fyneboot.InitMainWindow(initHandler, runtimeHandle)
	case G3n:
		g3nboot.InitMainWindow(initHandler, runtimeHandle)
	case Gio:
		gioboot.InitMainWindow(initHandler, runtimeHandle)
	case Gomobile:
		gomobileboot.InitMainWindow(initHandler, runtimeHandle)
	default:
	}
}
