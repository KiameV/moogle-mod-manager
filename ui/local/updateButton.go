package local

import (
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type UpdateButton struct {
	*widget.Button
	tm     mods.TrackedMod
	update func(tm mods.TrackedMod)
}

func NewUpdateButton(update func(tm mods.TrackedMod)) *UpdateButton {
	b := &UpdateButton{update: update}
	b.Button = widget.NewButton("Update", func() {
		if b.tm != nil && b.tm.UpdatedMod() != nil {
			b.update(b.tm)
		}
	})
	return b
}

func (b *UpdateButton) Refresh() {
	if b.tm != nil {
		b.Hidden = b.tm.UpdatedMod() == nil
	} else {
		b.Hidden = true
	}
}

func (b *UpdateButton) SetTrackedMod(tm mods.TrackedMod) {
	b.tm = tm
	b.Refresh()
}
