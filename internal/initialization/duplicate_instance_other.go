//go:build !windows

package initialization

import (
	"log"
	"os/exec"
	"runtime"
)

func notifyDuplicateInstance(url string) bool {
	const (
		title   = "برنامه قبلاً اجرا شده است"
		message = "یک نمونه از برنامه در حال اجراست.\nبرای باز کردن برنامه در مرورگر، دکمه «باز کردن برنامه» را انتخاب کنید."
	)

	// Linux desktop environments do not provide one universal native dialog.
	// Prefer the commonly installed GTK/KDE helpers and wait for the user's
	// choice before returning to Start.
	if runtime.GOOS == "linux" {
		if zenity, err := exec.LookPath("zenity"); err == nil {
			err := exec.Command(
				zenity,
				"--question",
				"--title="+title,
				"--text="+message,
				"--ok-label=باز کردن برنامه",
				"--cancel-label=انصراف",
			).Run()
			return err == nil
		}

		if kdialog, err := exec.LookPath("kdialog"); err == nil {
			err := exec.Command(
				kdialog,
				"--title", title,
				"--yes-label", "باز کردن برنامه",
				"--no-label", "انصراف",
				"--yesno", message,
			).Run()
			return err == nil
		}

		if xmessage, err := exec.LookPath("xmessage"); err == nil {
			err := exec.Command(
				xmessage,
				"-title", title,
				"-buttons", "باز کردن برنامه:0,انصراف:1",
				"-default", "انصراف",
				message,
			).Run()
			return err == nil
		}
	}

	log.Printf("%s؛ نمونه‌ی فعال در %s در دسترس است (برای نمایش dialog، zenity یا kdialog نصب کنید)", title, url)
	return false
}
