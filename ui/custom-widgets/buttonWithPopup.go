package custom_widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ButtonWithPopups struct {
	*widget.Button
	items []*fyne.MenuItem
}

func NewButtonWithPopups(label string, items ...*fyne.MenuItem) *ButtonWithPopups {
	b := &ButtonWithPopups{items: items}
	b.Button = widget.NewButton(label, func() {
		b.popUp()
	})
	return b
}

func (b *ButtonWithPopups) popUp() {
	c := fyne.CurrentApp().Driver().CanvasForObject(b)
	p := widget.NewPopUpMenu(fyne.NewMenu("", b.items...), c)
	p.ShowAtPosition(b.popUpPos())
}

func (b *ButtonWithPopups) popUpPos() fyne.Position {
	buttonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(b)
	return buttonPos.Add(fyne.NewPos(0, b.Size().Height-theme.InputBorderSize()))
}
