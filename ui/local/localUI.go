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
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"github.com/ncruces/zenity"
	"net/url"
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
			fyne.NewMenuItem("From File", func() {
				ui.addFromFile()
			}),
			fyne.NewMenuItem("From URL", func() {
				ui.addFromUrl()
			}))
		removeButton = widget.NewButton("Remove", func() {
			if ui.selectedMod != nil {
				if err := managed.RemoveMod(*state.CurrentGame, ui.selectedMod); err != nil {
					util.ShowErrorLong(err)
					return
				}
				ui.removeModFromList(ui.selectedMod)
				ui.selectedMod = nil
				ui.split.Trailing = container.NewMax()
			}
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
			ui.split.Trailing = ui.createPreview(ui.selectedMod)
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

func (ui *localUI) createPreview(tm *model.TrackedMod) fyne.CanvasObject {
	mod := tm.Mod
	c := container.NewVBox()
	if tm.UpdatedMod != nil {
		c.Add(widget.NewButton("Update", func() {
			if err := managed.UpdateMod(*state.CurrentGame, tm); err != nil {
				util.ShowErrorLong(err)
				return
			}
			ui.enableMod(*state.CurrentGame, tm)
			tm.DisplayName = tm.Mod.Name
		}))
	}
	c.Add(ui.createField("Name", mod.Name))
	c.Add(ui.createMultiLineField("Description", mod.Description))
	c.Add(ui.createLink("Link", mod.Link))
	c.Add(ui.createField("Author", mod.Author))
	c.Add(ui.createField("Category", mod.Category))
	k := mod.ModKind
	if k.Kind == mods.Hosted && k.Hosted != nil {
		c.Add(ui.createField("Version", k.Hosted.Version))
	} else if k.Nexus != nil {
		c.Add(ui.createField("Nexus Mod ID", k.Nexus.ID))
	}
	c.Add(ui.createField("Release Date", mod.ReleaseDate))

	if mod.ReleaseNotes != "" {
		c.Add(ui.createMultiLineField("Release Notes", mod.ReleaseDate))
	}
	if mod.ModCompatibility != nil && mod.ModCompatibility.HasItems() {
		c.Add(ui.createCompatibility(mod.ModCompatibility))
	}
	if mod.DonationLinks != nil && len(mod.DonationLinks) > 0 {
		c.Add(ui.createDonationLinks(mod.DonationLinks))
	}

	if img := mod.Preview.Get(); img != nil {
		c = container.NewBorder(img, nil, nil, nil, c)
	}
	return container.NewScroll(c)
}

func (ui *localUI) createField(name, value string) *fyne.Container {
	return container.NewHBox(
		widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(value),
	)
}

func (ui *localUI) createLink(name, value string) *fyne.Container {
	url, err := url.ParseRequestURI(value)
	if err != nil {
		return ui.createField(name, value)
	}
	return container.NewHBox(
		widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewHyperlink(value, url),
	)
}

func (ui *localUI) createMultiLineField(name, value string) *fyne.Container {
	rt := widget.NewRichTextFromMarkdown(value)
	s := container.NewScroll(rt)
	s.SetMinSize(fyne.NewSize(600, 300))
	return container.NewHBox(
		widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		s,
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

func (ui *localUI) createDonationLinks(links []*mods.DonationLink) fyne.CanvasObject {
	c := container.NewVBox(
		widget.NewLabelWithStyle("Support Project", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	for _, r := range links {
		if u, err := url.Parse(r.Link); err != nil {
			c.Add(widget.NewLabel("  - " + r.Name + ": " + r.Link))
		} else {
			c.Add(container.NewHBox(widget.NewLabel("  - "+r.Name), widget.NewHyperlink(r.Link, u)))
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
