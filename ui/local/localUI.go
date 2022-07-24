package local

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	ci "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/ncruces/zenity"
)

type LocalUI interface {
	state.Screen
	GetSelected() *model.TrackedMod
}

func New() LocalUI {
	return &localUI{}
}

type localUI struct {
	selectedMod *model.TrackedMod
	data        binding.UntypedList
}

func (ui *localUI) OnClose() {

}

func (ui *localUI) GetSelected() *model.TrackedMod {
	return ui.selectedMod
}

func (ui *localUI) Draw(w fyne.Window) {
	ui.data = binding.NewUntypedList()
	var (
		modList = widget.NewListWithData(
			ui.data,
			func() fyne.CanvasObject {
				return container.NewBorder(nil, nil, nil, widget.NewCheck("", func(b bool) {}), widget.NewLabel(""))
			},
			func(item binding.DataItem, co fyne.CanvasObject) {
				var tm *model.TrackedMod
				if i, ok := cw.GetValueFromDataItem(item); ok {
					if tm, ok = i.(*model.TrackedMod); ok {
						c := co.(*fyne.Container)
						c.Objects[0].(*widget.Label).SetText(tm.Mod.Name)
						c.Objects[1].(*widget.Check).Bind(newEnableBind(ui, tm))
					}
				}
			})
		addButton = cw.NewButtonWithPopups("Add",
			fyne.NewMenuItem("From File", func() {
				ui.addFromFile()
			}),
			fyne.NewMenuItem("From URL", func() {
				ui.addFromUrl()
			}))
		removeButton = widget.NewButton("Remove", func() {
			if ui.selectedMod != nil {
				if err := managed.RemoveMod(*state.CurrentGame, ui.selectedMod.GetModID()); err != nil {
					dialog.ShowError(err, state.Window)
					return
				}
				ui.selectedMod = nil
				ui.Draw(w)
			}
		})
		modDetails = container.NewScroll(container.NewMax())
	)

	for _, mod := range managed.GetMods(*state.CurrentGame) {
		ui.addModToList(mod)
	}

	removeButton.Disable()
	modList.OnSelected = func(id widget.ListItemID) {
		data, err := ui.data.GetItem(id)
		if err != nil {
			return
		}
		if i, ok := cw.GetValueFromDataItem(data); ok {
			ui.selectedMod = i.(*model.TrackedMod)
			removeButton.Enable()
			modDetails.Content = container.NewCenter(widget.NewLabel("Loading..."))
			modDetails.Refresh()
			modDetails.Content = ui.createPreview(ui.selectedMod.Mod)
			modDetails.Refresh()
		}
	}
	modList.OnUnselected = func(id widget.ListItemID) {
		ui.selectedMod = nil
		removeButton.Disable()
		modDetails.Hide()
	}

	buttons := container.NewHBox(addButton, widget.NewSeparator(), removeButton)

	split := container.NewHSplit(
		modList,
		modDetails)
	split.SetOffset(0.25)

	w.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(config.GameNameString(*state.CurrentGame), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			buttons,
		), nil, nil, nil,
		split))
}

func (ui *localUI) createPreview(mod *mods.Mod) fyne.CanvasObject {
	c := container.NewVBox(
		ui.createField("Name", mod.Name),
		ui.createMultiLineField("Description", mod.Description),
		ui.createField("Version", mod.Version),
		ui.createField("Link", mod.Link),
		ui.createField("Author", mod.Author),
		ui.createField("Category", mod.Category),
		ui.createField("Release Date", mod.ReleaseDate),
	)
	if mod.ReleaseNotes != "" {
		c.Add(ui.createMultiLineField("Release Notes", mod.ReleaseDate))
	}
	if mod.ModCompatibility != nil && mod.ModCompatibility.HasItems() {
		c.Add(ui.createCompatibility(mod.ModCompatibility))
	}

	if img := mod.Preview.Get(); img != nil {
		c = container.NewBorder(img, nil, nil, nil, c)
	}
	return c
}

func (ui *localUI) createField(name, value string) *fyne.Container {
	return container.NewHBox(
		widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(value),
	)
}

func (ui *localUI) createMultiLineField(name, value string) *fyne.Container {
	return container.NewHBox(
		widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewRichTextFromMarkdown(value),
	)
}

func (ui *localUI) createCompatibility(compatibility *mods.ModCompatibility) fyne.CanvasObject {
	c := container.NewVBox(
		widget.NewLabelWithStyle("Compatibility", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	if len(compatibility.Requires) > 0 {
		c.Add(widget.NewLabelWithStyle("  Requires", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for _, r := range compatibility.Requires {
			c.Add(widget.NewLabel("  - " + r.Name + ": " + r.Source))
		}
	}
	if len(compatibility.Requires) > 0 {
		c.Add(widget.NewLabelWithStyle("  Forbids", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for _, r := range compatibility.Requires {
			c.Add(widget.NewLabel("  - " + r.Name + ": " + r.Source))
		}
	}
	return c
}

func (ui *localUI) addFromFile() {
	var tm *model.TrackedMod
	if file, err := zenity.SelectFile(
		zenity.Title("Select a mod file"),
		zenity.FileFilter{
			Name:     "mod file",
			Patterns: []string{"*.xml", "*.json"},
		}); err == nil {
		if tm, err = managed.AddModFromFile(*state.CurrentGame, file); err != nil {
			dialog.ShowError(err, state.Window)
			return
		} else {
			ui.addModToList(tm)
		}
	}
}

func (ui *localUI) addFromUrl() {
	e := widget.NewEntry()
	dialog.ShowForm("Add Remote mod file", "Add", "Cancel",
		[]*widget.FormItem{widget.NewFormItem("URL", e)},
		func(ok bool) {
			if ok && e.Text != "" {
				if tm, err := managed.AddModFromUrl(*state.CurrentGame, e.Text); err != nil {
					dialog.ShowError(err, state.Window)
					return
				} else {
					ui.addModToList(tm)
				}
			}
		}, state.Window)
}

func (ui *localUI) addModToList(mod *model.TrackedMod) {
	u := binding.NewUntyped()
	if err := u.Set(mod); err == nil {
		_ = ui.data.Append(u)
	}
}

func (ui *localUI) toggleEnabled(mod *model.TrackedMod) bool {
	if mod.Enabled {
		return ui.enableMod(mod)
	}
	return ui.disableMod(mod)
}

func (ui *localUI) enableMod(tm *model.TrackedMod) bool {
	if len(tm.Mod.Configurations) > 0 {
		if err := state.GetScreen(state.ConfigInstaller).(ci.ConfigInstaller).Setup(tm.Mod, tm.GetDir(), func(tis []*mods.ToInstall) {
			managed.EnableMod(*state.CurrentGame, tm.Mod, tis)
		}); err != nil {
			dialog.ShowError(err, state.Window)
			return false
		}
		state.ShowScreen(state.ConfigInstaller)
	} else {
		tis, err := mods.NewToInstallForMod(ui.selectedMod.Mod, ui.selectedMod.Mod.AlwaysDownload)
		if err != nil {
			dialog.ShowError(err, state.Window)
			return false
		}
		managed.EnableMod(*state.CurrentGame, tm.Mod, tis)
	}
	return true
}

func (ui *localUI) disableMod(mod *model.TrackedMod) bool {
	return true
}
