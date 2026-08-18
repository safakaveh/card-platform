//go:build windows

package initialization

import (
	"syscall"
	"unsafe"
)

// notifyDuplicateInstance displays a native Windows dialog. The OK button is
// intentionally the action that opens the already running application.
func notifyDuplicateInstance(url string) bool {
	const (
		mbOKCancel    = 0x00000001
		mbIconWarning = 0x00000030
		idOK          = 1
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	messageBox := user32.NewProc("MessageBoxW")
	title, _ := syscall.UTF16PtrFromString("برنامه قبلاً اجرا شده است")
	message, _ := syscall.UTF16PtrFromString("یک نمونه از برنامه در حال اجراست.\nبرای باز کردن برنامه در مرورگر، تأیید را انتخاب کنید.")
	result, _, _ := messageBox.Call(
		0,
		uintptr(unsafe.Pointer(message)),
		uintptr(unsafe.Pointer(title)),
		mbOKCancel|mbIconWarning,
	)
	return result == idOK
}
