package conflicts

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/files"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"path/filepath"
)

func ShowConflicts(mod *mods.Mod, conflicts []*files.Conflict, done func(mods.Result)) {
	f := widget.NewForm()
	for _, c := range conflicts {
		f.Items = append(f.Items, createItem(mod, c))
	}
	d := dialog.NewCustomConfirm("Conflicts", "ok", "cancel", container.NewVScroll(f), func(ok bool) {
		r := mods.Ok
		if !ok {
			r = mods.Cancel
		}
		done(r)
	}, ui.Window)
	d.Resize(fyne.NewSize(400, 400))
	d.Show()
}

func createItem(mod *mods.Mod, c *files.Conflict) *widget.FormItem {
	c.Selection = mod
	return widget.NewFormItem(
		filepath.Base(c.Path),
		widget.NewSelect([]string{string(mod.Name), string(c.Owner.Name)}, func(s string) {
			if s == string(mod.Name) {
				c.Selection = mod
			} else {
				c.Selection = c.Owner
			}
		}))
}
