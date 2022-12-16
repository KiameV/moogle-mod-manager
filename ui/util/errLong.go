package util

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

func ShowErrorLong(err error, w ...fyne.Window) {
	var (
		window = state.Window
		text   = widget.NewRichTextWithText(err.Error())
	)
	if len(w) > 0 && w[0] != nil {
		window = w[0]
	}
	text.Wrapping = fyne.TextWrapBreak

	button := widget.NewButton("Copy To Clipboard", func() {
		_ = clipboard.WriteAll(err.Error())
	})

	errDialog := dialog.NewCustom("Error", "OK", container.NewBorder(button, nil, nil, nil, container.NewVScroll(text)), window)
	errDialog.Resize(fyne.NewSize(500, 400))
	errDialog.Show()
}
