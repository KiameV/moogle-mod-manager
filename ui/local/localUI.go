package local

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/actions"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	mp "github.com/kiamev/moogle-mod-manager/ui/mod-preview"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	u "github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"github.com/ncruces/zenity"
)

type LocalUI interface {
	state.Screen
	GetSelected() mods.TrackedMod
}

func New() LocalUI {
	return &localUI{}
}

type localUI struct {
	selectedMod   mods.TrackedMod
	data          binding.UntypedList
	split         *container.Split
	checkAll      *widget.Button
	ModList       *widget.List
	workingDialog dialog.Dialog
	mods          []mods.TrackedMod
}

func (ui *localUI) PreDraw(fyne.Window, ...interface{}) error { return nil }

func (ui *localUI) OnClose() {}

func (ui *localUI) GetSelected() mods.TrackedMod {
	return ui.selectedMod
}

func (ui *localUI) DrawAsDialog(fyne.Window) {}

func (ui *localUI) Draw(w fyne.Window) {
	ui.data = binding.NewUntypedList()
	ui.mods = make([]mods.TrackedMod, 0)
	ui.ModList = widget.NewListWithData(
		ui.data,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, widget.NewLabel(""), widget.NewCheck("", func(b bool) {}), NewUpdateButton(ui.updateMod))
		},
		func(item binding.DataItem, co fyne.CanvasObject) {
			var tm mods.TrackedMod
			if i, ok := cw.GetValueFromDataItem(item); ok {
				if tm, ok = i.(mods.TrackedMod); ok {
					if tm.DisplayName() == "" {
						tm.SetDisplayName(string(tm.Mod().Name))
					}
					c := co.(*fyne.Container)
					c.Objects[0].(*UpdateButton).SetTrackedMod(tm)
					c.Objects[1].(*widget.Label).Bind(binding.BindString(tm.DisplayNamePtr()))
					c.Objects[2].(*widget.Check).Bind(newEnableBind(ui, tm, ui.startEnableDisableCallback, ui.showWorkingDialog, ui.hideDialog))
				}
			}
		})

	addButton := cw.NewButtonWithPopups("Add",
		fyne.NewMenuItem("From File", func() {
			ui.addFromFile()
		}),
		fyne.NewMenuItem("From URL", func() {
			ui.addFromUrl()
		}))
	findButton := widget.NewButton("Find", func() {
		state.ShowScreen(state.DiscoverMods, ui.mods)
	})
	removeButton := widget.NewButton("Remove", func() {
		dialog.NewConfirm("Delete?", "Are you sure you want to delete this mod?", func(ok bool) {
			if ok && ui.selectedMod != nil {
				if err := managed.RemoveMod(state.CurrentGame, ui.selectedMod); err != nil {
					util.ShowErrorLong(err)
					return
				}
				for i, m := range ui.mods {
					if m == ui.selectedMod {
						ui.mods = append(ui.mods[:i], ui.mods[i+1:]...)
						break
					}
				}
				ui.removeModFromList(ui.selectedMod)
				ui.selectedMod = nil
				ui.ModList.UnselectAll()
				ui.split.Trailing = container.NewMax()
				ui.split.Refresh()
			}
		}, u.Window).Show()
	})

	ui.checkAll = widget.NewButton("Check For Updates", func() {
		ui.checkAll.Disable()
		defer func() {
			ui.split.Refresh()
			ui.checkAll.Enable()
		}()
		managed.CheckForUpdates(state.CurrentGame, func(err error) {
			if err != nil {
				util.ShowErrorLong(err)
			} else {
				dialog.ShowInformation("Check for updates", "Done checking for updates.", u.Window)
			}
		})
	})

	for _, mod := range managed.GetMods(state.CurrentGame) {
		ui.addModToList(mod)
	}

	removeButton.Disable()
	ui.ModList.OnSelected = func(id widget.ListItemID) {
		data, err := ui.data.GetItem(id)
		if err != nil {
			return
		}
		if i, ok := cw.GetValueFromDataItem(data); ok {
			ui.selectedMod = i.(mods.TrackedMod)
			removeButton.Enable()
			ui.split.Trailing = container.NewCenter(widget.NewLabel("Loading..."))
			ui.split.Refresh()
			ui.split.Trailing = mp.CreatePreview(ui.selectedMod.Mod(), mp.ModPreviewOptions{
				UpdateCallback: func(tm mods.TrackedMod) {
					ui.updateMod(tm)

					/*if err = managed.UpdateMod(state.CurrentGame, tm); err != nil {
						util.ShowErrorLong(err)
						return
					}
					if err = newEnableBind(ui.selectedMod, ui.startEnableDisableCallback, ui.showWorkingDialog, ui.endEnableDisableCallback).EnableMod(); err != nil {
						util.ShowErrorLong(err)
						return
					}*/
				},
				TrackedMod: ui.selectedMod,
			})
			ui.split.Refresh()
		}
	}
	ui.ModList.OnUnselected = func(id widget.ListItemID) {
		ui.selectedMod = nil
		removeButton.Disable()
		ui.split.Trailing = container.NewMax()
	}

	buttons := container.NewHBox(findButton, addButton, removeButton, ui.checkAll)
	ui.split = container.NewHSplit(
		ui.ModList,
		container.NewMax())
	ui.split.SetOffset(0.25)

	w.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(string(state.CurrentGame.Name()), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			buttons,
		), nil, nil, nil,
		ui.split))
}

func (ui *localUI) addFromFile() {
	var tm mods.TrackedMod
	if file, err := zenity.SelectFile(
		zenity.Title("Select a mod file"),
		zenity.FileFilter{
			Name:     "mod file",
			Patterns: []string{"*.xml", "*.json"},
		}); err == nil {
		if tm, err = managed.AddModFromFile(state.CurrentGame, file); err != nil {
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
				if tm, err := managed.AddModFromUrl(state.CurrentGame, e.Text); err != nil {
					util.ShowErrorLong(err)
					return
				} else {
					ui.addModToList(tm)
				}
			}
		}, u.Window)
}

func (ui *localUI) addModToList(mod mods.TrackedMod) {
	u := binding.NewUntyped()
	if err := u.Set(mod); err == nil {
		_ = ui.data.Append(u)
	}
	ui.mods = append(ui.mods, mod)
}

func (ui *localUI) removeModFromList(mod mods.TrackedMod) {
	var item binding.DataItem
	sl, err := ui.data.Get()
	if err != nil {
		util.ShowErrorLong(err)
		return
	}
	for i := 0; i < len(sl); i++ {
		if item, err = ui.data.GetItem(i); err != nil {
			util.ShowErrorLong(err)
			return
		}
		if j, ok := cw.GetValueFromDataItem(item); ok {
			if j == mod {
				sl = append(sl[:i], sl[i+1:]...)
				if err = ui.data.Set(sl); err != nil {
					util.ShowErrorLong(err)
				}
				return
			}
		}
	}
	for i, m := range ui.mods {
		if m == mod {
			ui.mods = append(ui.mods[:i], ui.mods[i+1:]...)
			return
		}
	}
}

func (ui *localUI) startEnableDisableCallback() bool {
	return ui.workingDialog == nil
}

func (ui *localUI) showWorkingDialog() {
	if ui.workingDialog == nil {
		ui.workingDialog = dialog.NewInformation("Working", "Working...", u.Window)
		ui.workingDialog.Show()
	}
}

func (ui *localUI) hideDialog() {
	if ui.workingDialog != nil {
		ui.workingDialog.Hide()
		ui.workingDialog = nil
	}
}

/*func (ui *localUI) endEnableDisableCallback(result mods.Result, err ...error) {
	if ui.workingDialog != nil {
		ui.workingDialog.Hide()
		ui.workingDialog = nil
	}
	if result == mods.Error && len(err) == 0 {
		util.ShowErrorLong(errors.New("result is Error but no error messages received"))
	} else if result != mods.Error && len(err) > 0 && err[0] != nil {
		util.ShowErrorLong(fmt.Errorf("result was not Error but received error: %v", err[0]))
	} else if result == mods.Error && len(err) > 0 {
		util.ShowErrorLong(err[0])
	}
	ui.split.Leading.Refresh()
}*/

func (ui *localUI) updateMod(tm mods.TrackedMod) {
	if action, err := actions.New(actions.Update, actions.Params{
		Game: state.CurrentGame,
		Mod:  tm,
		WorkingDialog: actions.WorkingDialog{
			Show: ui.showWorkingDialog,
			Hide: ui.hideDialog,
		},
	}, nil); err != nil {
		util.ShowErrorLong(err)
	} else if err = action.Run(); err != nil {
		util.ShowErrorLong(err)
	} else {
		tm.SetDisplayName(string(tm.Mod().Name))
	}
	/*if err := managed.UpdateMod(state.CurrentGame, tm); err != nil {
		util.ShowErrorLong(err)
	}*/
	ui.split.Leading.Refresh()
}
