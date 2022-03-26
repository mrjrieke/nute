//go:build !gomobileboot
// +build !gomobileboot

package gomobileboot

func InitMainWindow(initHandler interface{}, runtimeHandler interface{}) interface{} {
	return nil
}
