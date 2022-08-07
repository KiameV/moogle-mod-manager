package util

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"golang.design/x/clipboard"
)

func ShowErrorLong(err error) {
	text := widget.NewRichTextWithText(err.Error())
	text.Wrapping = fyne.TextWrapBreak

	button := widget.NewButton("Copy Error Text", func() {
		clipboard.Write(clipboard.FmtText, []byte(err.Error()))
	})

	errDialog := dialog.NewCustom("Error", "OK", container.NewBorder(button, nil, nil, nil, container.NewVScroll(text)), state.Window)
	errDialog.Resize(fyne.NewSize(500, 400))
	errDialog.Show()
}
