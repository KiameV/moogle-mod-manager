package local

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/pr-modsync/mods/managed"
	"github.com/kiamev/pr-modsync/ui/menu"
	"github.com/kiamev/pr-modsync/ui/state"
)

func Draw(w fyne.Window) {
	menu.Add(w)
	modList := container.NewVBox()
	for _, mod := range managed.GetMods(*state.CurrentGame) {
		modList.Objects = append(modList.Objects, widget.NewLabel(mod.Mod.Name))
	}
	modDetails := container.NewScroll(widget.NewRichText())

	split := container.NewHSplit(
		container.NewVScroll(modList),
		modDetails)
	//split.Offset
	w.SetContent(split)
}
