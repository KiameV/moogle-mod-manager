package util

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

func ShowErrorLong(err error) {
	var text = widget.NewRichTextWithText(err.Error())
	text.Wrapping = fyne.TextWrapBreak

	button := widget.NewButton("Copy To Clipboard", func() {
		_ = clipboard.WriteAll(err.Error())
	})

	errDialog := dialog.NewCustom("Error", "OK", container.NewBorder(button, nil, nil, nil, container.NewVScroll(text)), ui.ActiveWindow())
	errDialog.Resize(fyne.NewSize(500, 400))
	errDialog.Show()
}
