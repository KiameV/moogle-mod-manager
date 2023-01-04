package util

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"net/url"
)

func CreateUrlRow(value string) *fyne.Container {
	t := widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
		_ = clipboard.WriteAll(value)
	})
	if u, err := url.Parse(value); err == nil {
		return container.NewBorder(nil, nil, nil, widget.NewToolbar(t), widget.NewHyperlink(value, u))
	}
	return container.NewBorder(nil, nil, nil, widget.NewToolbar(t), widget.NewLabel(value))
}
