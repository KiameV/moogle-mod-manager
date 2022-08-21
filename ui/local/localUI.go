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
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	mp "github.com/kiamev/moogle-mod-manager/ui/mod-preview"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"github.com/ncruces/zenity"
)

type LocalUI interface {
	state.Screen
	GetSelected() *mods.TrackedMod
}

func New() LocalUI {
	return &localUI{}
}

type localUI struct {
	selectedMod *mods.TrackedMod
	data        binding.UntypedList
	split       *container.Split
	checkAll    *widget.Button
}

func (ui *localUI) PreDraw() error { return nil }

func (ui *localUI) OnClose() {}

func (ui *localUI) GetSelected() *mods.TrackedMod {
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
				var tm *mods.TrackedMod
				if i, ok := cw.GetValueFromDataItem(item); ok {
					if tm, ok = i.(*mods.TrackedMod); ok {
						if tm.DisplayName == "" {
							tm.DisplayName = tm.Mod.Name
						}
						c := co.(*fyne.Container)
						c.Objects[0].(*widget.Label).Bind(binding.BindString(&tm.DisplayName))
						c.Objects[1].(*widget.Check).Bind(newEnableBind(tm, ui.enableDisableCallback))
					}
				}
			})
		addButton = cw.NewButtonWithPopups("Add",
			fyne.NewMenuItem("Find", func() {
				state.ShowScreen(state.DiscoverMods)
			}),
			fyne.NewMenuItem("From File", func() {
				ui.addFromFile()
			}),
			fyne.NewMenuItem("From URL", func() {
				ui.addFromUrl()
			}))
		removeButton = widget.NewButton("Remove", func() {
			dialog.NewConfirm("Delete?", "Are you sure you want to delete this mod?", func(ok bool) {
				if ok && ui.selectedMod != nil {
					if err := managed.RemoveMod(*state.CurrentGame, ui.selectedMod); err != nil {
						util.ShowErrorLong(err)
						return
					}
					ui.removeModFromList(ui.selectedMod)
					ui.selectedMod = nil
					ui.split.Trailing = container.NewMax()
				}
			}, state.Window).Show()
		})
	)
	ui.checkAll = widget.NewButton("Check All", func() {
		ui.checkAll.Disable()
		defer func() {
			ui.split.Refresh()
			ui.checkAll.Enable()
		}()
		managed.CheckForUpdates(*state.CurrentGame, func(err error) {
			if err != nil {
				util.ShowErrorLong(err)
			} else {
				dialog.ShowInformation("Check for updates", "Done checking for updates.", state.Window)
			}
		})
	})

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
			ui.selectedMod = i.(*mods.TrackedMod)
			removeButton.Enable()
			ui.split.Trailing = container.NewCenter(widget.NewLabel("Loading..."))
			ui.split.Refresh()
			ui.split.Trailing = mp.CreatePreview(ui.selectedMod.Mod, mp.ModPreviewOptions{
				UpdateCallback: func(tm *mods.TrackedMod) {
					var err error
					if err = managed.UpdateMod(*state.CurrentGame, tm); err != nil {
						util.ShowErrorLong(err)
						return
					}
					if err = newEnableBind(ui.selectedMod, ui.enableDisableCallback).EnableMod(); err != nil {
						util.ShowErrorLong(err)
						return
					}
					tm.DisplayName = tm.Mod.Name
				},
				TrackedMod: ui.selectedMod,
			})
			ui.split.Refresh()
		}
	}
	modList.OnUnselected = func(id widget.ListItemID) {
		ui.selectedMod = nil
		removeButton.Disable()
		ui.split.Trailing = container.NewMax()
	}

	buttons := container.NewHBox(addButton, removeButton, ui.checkAll)
	ui.split = container.NewHSplit(
		modList,
		container.NewMax())
	ui.split.SetOffset(0.25)

	w.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(config.GameNameString(*state.CurrentGame), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			buttons,
		), nil, nil, nil,
		ui.split))
}

func (ui *localUI) addFromFile() {
	var tm *mods.TrackedMod
	if file, err := zenity.SelectFile(
		zenity.Title("Select a mod file"),
		zenity.FileFilter{
			Name:     "mod file",
			Patterns: []string{"*.xml", "*.json"},
		}); err == nil {
		if tm, err = managed.AddModFromFile(*state.CurrentGame, file); err != nil {
			util.ShowErrorLong(err)
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
					util.ShowErrorLong(err)
					return
				} else {
					ui.addModToList(tm)
				}
			}
		}, state.Window)
}

func (ui *localUI) addModToList(mod *mods.TrackedMod) {
	u := binding.NewUntyped()
	if err := u.Set(mod); err == nil {
		_ = ui.data.Append(u)
	}
}

func (ui *localUI) removeModFromList(mod *mods.TrackedMod) {
	var item binding.DataItem
	sl, err := ui.data.Get()
	if err != nil {
		// TODO message
		return
	}
	for i := 0; i < len(sl); i++ {
		if item, err = ui.data.GetItem(i); err != nil {
			// TODO message
			return
		}
		if j, ok := cw.GetValueFromDataItem(item); ok {
			if j == mod {
				sl = append(sl[:i], sl[i+1:]...)
				if err = ui.data.Set(sl); err != nil {
					// TODO message
				}
				return
			}
		}
	}
	return
}

func (ui *localUI) showInputs(yes bool) {
	if yes {
		ui.split.Leading.Show()
	} else {
		ui.split.Leading.Hide()
	}
	ui.split.Refresh()
}

func (ui *localUI) enableDisableCallback(err error) {

}
