//go:build g3nboot
// +build g3nboot

package g3nboot

import (
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/renderer"
)

func InitMainWindow(initHandler interface{}, runtimeHandler interface{}) interface{} {
	a := app.App()

	g3nInit := initHandler.(func(a app.App))
	g3nInit(a)
	g3nRuntimeHandler := runtimeHandler.(func(renderer *renderer.Renderer, deltaTime time.Duration))

	// Run the application
	a.Run(g3nRuntimeHandler)
	return a
}
