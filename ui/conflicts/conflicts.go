package conflicts

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/files"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

type item struct {
	*widget.Select
	Conflict *files.Conflict
}

func ShowConflicts(mod *mods.Mod, conflicts []*files.Conflict, done func(mods.Result)) {
	var (
		d     dialog.Dialog
		f     = widget.NewForm()
		items = createItems(f, mod, conflicts)
		c     = container.NewBorder(
			container.NewHBox(
				widget.NewButton("Skip All", func() {
					for _, i := range items {
						i.skip()
					}
				}),
				widget.NewButton("Overwrite All", func() {
					for _, i := range items {
						i.overwrite(mod)
					}
				})), nil, nil, nil,
			container.NewVScroll(f))
	)
	d = dialog.NewCustomConfirm("Conflicts", "ok", "cancel", c, func(ok bool) {
		for _, i := range items {
			if i.Conflict.Selection == nil {
				ShowConflicts(mod, conflicts, done)
				dialog.ShowInformation("Error", "Please select an option for all conflicts", ui.ActiveWindow())
				return
			}
		}
		r := mods.Ok
		if !ok {
			r = mods.Cancel
		}
		done(r)
	}, ui.Window)
	d.Resize(fyne.NewSize(400, 400))
	d.Show()
}

func createItems(form *widget.Form, mod *mods.Mod, conflicts []*files.Conflict) (items []*item) {
	items = make([]*item, len(conflicts))
	for i, c := range conflicts {
		j := &item{
			Select: widget.NewSelect([]string{string(mod.Name), string(c.Owner.Name)}, func(s string) {
				if s == string(mod.Name) {
					c.Selection = mod
				} else {
					c.Selection = c.Owner
				}
			}),
			Conflict: c,
		}
		items[i] = j
		form.AppendItem(widget.NewFormItem(filepath.Base(c.Path), j.Select))
	}
	return
}

func (i item) skip() {
	i.Conflict.Selection = i.Conflict.Owner
	i.Select.SetSelected(string(i.Conflict.Owner.Name))
}

func (i item) overwrite(mod *mods.Mod) {
	i.Conflict.Selection = mod
	i.Select.SetSelected(string(mod.Name))
}
