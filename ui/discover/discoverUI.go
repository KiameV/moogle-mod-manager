package discover

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/repo"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	mp "github.com/kiamev/moogle-mod-manager/ui/mod-preview"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util"
)

func New() state.Screen {
	return &discoverUI{}
}

type discoverUI struct {
	selectedMod *mods.Mod
	data        binding.UntypedList
	split       *container.Split
	mods        []*mods.Mod
	localMods   map[string]bool
}

func (ui *discoverUI) OnClose() {}

func (ui *discoverUI) PreDraw(w fyne.Window, args ...interface{}) (err error) {
	var (
		d        = dialog.NewInformation("", "Finding Mods...", w)
		repoMods []*mods.Mod
		ok       bool
	)
	defer d.Hide()
	d.Show()

	ui.localMods = make(map[string]bool)
	for _, tm := range args[0].([]interface{})[0].([]*mods.TrackedMod) {
		ui.localMods[tm.GetModID()] = true
	}

	if repoMods, err = repo.GetMods(*state.CurrentGame); err != nil {
		return
	}

	ui.mods = make([]*mods.Mod, 0, len(repoMods))
	for _, m := range repoMods {
		if _, ok = ui.localMods[m.ID]; !ok {
			ui.mods = append(ui.mods, m)
		}
	}
	return
}

func (ui *discoverUI) DrawAsDialog(w fyne.Window) {
	ui.draw(w, true)
}

func (ui *discoverUI) Draw(w fyne.Window) {
	ui.draw(w, false)
}

func (ui *discoverUI) draw(w fyne.Window, isPopup bool) {
	if len(ui.mods) == 0 {
		// TODO
		return
	}
	ui.data = binding.NewUntypedList()
	modList := widget.NewListWithData(
		ui.data,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item binding.DataItem, co fyne.CanvasObject) {
			var m *mods.Mod
			if i, ok := cw.GetValueFromDataItem(item); ok {
				if m, ok = i.(*mods.Mod); ok {
					co.(*widget.Label).SetText(m.Name)
				}
			}
		})
	for _, m := range ui.mods {
		if err := ui.data.Append(m); err != nil {
			util.ShowErrorLong(err)
			// TODO
			return
		}
	}

	ui.split = container.NewHSplit(modList, container.NewMax())
	ui.split.SetOffset(0.25)

	modList.OnSelected = func(id widget.ListItemID) {
		data, err := ui.data.GetItem(id)
		if err != nil {
			util.ShowErrorLong(err)
			return
		}
		if i, ok := cw.GetValueFromDataItem(data); ok {
			ui.selectedMod = i.(*mods.Mod)
		}
		ui.split.Trailing = container.NewCenter(widget.NewLabel("Loading..."))
		ui.split.Refresh()
		ui.split.Trailing = container.NewBorder(
			container.NewHBox(widget.NewButton("Include Mod", func() {
				mod := ui.selectedMod
				if err := managed.AddMod(*state.CurrentGame, mods.NewTrackerMod(mod, *state.CurrentGame)); err != nil {
					util.ShowErrorLong(err)
					return
				}
				for i, m := range ui.mods {
					if m == mod {
						ui.mods = append(ui.mods[:i], ui.mods[i+1:]...)
						break
					}
				}
				sl := make([]interface{}, len(ui.mods))
				for i, m := range ui.mods {
					sl[i] = m
				}
				if err := ui.data.Set(sl); err != nil {
					util.ShowErrorLong(err)
					return
				}
				ui.selectedMod = nil
				ui.split.Trailing = container.NewMax()
				ui.split.Refresh()
				state.UpdateCurrentScreen()
			})), nil, nil, nil,
			mp.CreatePreview(ui.selectedMod))
		ui.split.Refresh()
	}

	w.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(config.GameNameString(*state.CurrentGame), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
		), nil, nil, nil, container.NewBorder(
			container.NewHBox(widget.NewButton("Back", func() {
				if isPopup {
					state.ClosePopupWindow()
				} else {
					state.ShowPreviousScreen()
				}
			})), nil, nil, nil,
			ui.split)))
}
