package ui

import "fyne.io/fyne/v2"

var (
	App          fyne.App
	Window       fyne.Window
	PopupWindow  fyne.Window
	ShowingPopup bool
)

func ActiveWindow() fyne.Window {
	if ShowingPopup {
		return PopupWindow
	}
	return Window
}
