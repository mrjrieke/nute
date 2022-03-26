//go:build !gioboot
// +build !gioboot

package gioboot

func InitMainWindow(initHandler interface{}, runtimeHandler interface{}) interface{} {
	return nil
}
