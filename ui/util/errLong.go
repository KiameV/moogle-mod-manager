package util

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

func ShowErrorLong(err error) {
	label := widget.NewLabel(err.Error())
	label.Wrapping = fyne.TextWrapBreak
	scroll := container.NewVScroll(label)
	errDialog := dialog.NewCustom("Error", "OK", container.NewMax(scroll), state.Window)
	factor := float32(0.40)
	winSize := state.Window.Content().Size()
	errDialog.Resize(fyne.NewSize(winSize.Width*factor, winSize.Height*factor))
	errDialog.Show()
}
