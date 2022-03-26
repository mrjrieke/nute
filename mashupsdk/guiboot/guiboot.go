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

func InitMainWindow(guiType GuiProvider, initHandler interface{}, runtimeHandle interface{}) interface{} {

	switch guiType {
	case Fyne:
		return fyneboot.InitMainWindow(initHandler, runtimeHandle)
	case G3n:
		return g3nboot.InitMainWindow(initHandler, runtimeHandle)
	case Gio:
		return gioboot.InitMainWindow(initHandler, runtimeHandle)
	case Gomobile:
		return gomobileboot.InitMainWindow(initHandler, runtimeHandle)
	default:
		return nil
	}
}
