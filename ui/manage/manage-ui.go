package local

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

func New() state.Screen {
	return &manageMods{}
}

type manageMods struct {
}

func (m *manageMods) Draw(w fyne.Window) {
	modList := container.NewVBox()
	for _, mod := range managed.GetMods(*state.CurrentGame) {
		modList.Objects = append(modList.Objects, widget.NewLabel(mod.Mod.Name))
	}
	modDetails := container.NewScroll(widget.NewRichText())

	split := container.NewHSplit(
		container.NewVScroll(modList),
		modDetails)

	w.SetContent(split)
}
