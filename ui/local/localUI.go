package local

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	ci "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	mp "github.com/kiamev/moogle-mod-manager/ui/mod-preview"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"github.com/ncruces/zenity"
	"path/filepath"
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
	split       *container.Split
	checkAll    *widget.Button
}

func (ui *localUI) PreDraw() error { return nil }

func (ui *localUI) OnClose() {}

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
						if tm.DisplayName == "" {
							tm.DisplayName = tm.Mod.Name
						}
						c := co.(*fyne.Container)
						c.Objects[0].(*widget.Label).Bind(binding.BindString(&tm.DisplayName))
						c.Objects[1].(*widget.Check).Bind(newEnableBind(ui, tm))
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
			ui.selectedMod = i.(*model.TrackedMod)
			removeButton.Enable()
			ui.split.Trailing = container.NewCenter(widget.NewLabel("Loading..."))
			ui.split.Refresh()
			ui.split.Trailing = mp.CreatePreview(ui.selectedMod.Mod, mp.ModPreviewOptions{
				UpdateCallback: func(tm *model.TrackedMod) {
					if err := managed.UpdateMod(*state.CurrentGame, tm); err != nil {
						util.ShowErrorLong(err)
						return
					}
					ui.enableMod(*state.CurrentGame, tm)
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
	var tm *model.TrackedMod
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

func (ui *localUI) addModToList(mod *model.TrackedMod) {
	u := binding.NewUntyped()
	if err := u.Set(mod); err == nil {
		_ = ui.data.Append(u)
	}
}

func (ui *localUI) removeModFromList(mod *model.TrackedMod) {
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

func (ui *localUI) toggleEnabled(game config.Game, mod *model.TrackedMod) bool {
	if mod.Enabled {
		return ui.enableMod(game, mod)
	}
	return ui.disableMod(mod)
}

func (ui *localUI) enableMod(game config.Game, tm *model.TrackedMod) bool {
	if len(tm.Mod.Configurations) > 0 {
		ui.showInputs(false)
		var modPath = filepath.Join(config.Get().GetModsFullPath(game), tm.GetDirSuffix())
		if err := state.GetScreen(state.ConfigInstaller).(ci.ConfigInstaller).Setup(tm.Mod, modPath, func(tis []*model.ToInstall) error {
			result := managed.EnableMod(*state.CurrentGame, tm, tis)
			ui.showInputs(true)
			return result
		}); err != nil {
			ui.showInputs(true)
			return false
		}
		state.ShowScreen(state.ConfigInstaller)
	} else {
		tis, err := model.NewToInstallForMod(tm.Mod.ModKind.Kind, tm.Mod, tm.Mod.AlwaysDownload)
		if err != nil {
			ui.showInputs(true)
			util.ShowErrorLong(err)
			return false
		}
		if err = managed.EnableMod(*state.CurrentGame, tm, tis); err != nil {
			ui.showInputs(true)
			return false
		}
	}
	return true
}

func (ui *localUI) disableMod(mod *model.TrackedMod) bool {
	if err := managed.DisableMod(*state.CurrentGame, mod); err != nil {
		util.ShowErrorLong(err)
		return false
	}
	return true
}

func (ui *localUI) showInputs(yes bool) {
	if yes {
		ui.split.Leading.Show()
	} else {
		ui.split.Leading.Hide()
	}
	ui.split.Refresh()
}
