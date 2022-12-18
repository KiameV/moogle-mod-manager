package discover

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/discover"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	mp "github.com/kiamev/moogle-mod-manager/ui/mod-preview"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"strings"
)

func New() state.Screen {
	return &discoverUI{}
}

type discoverUI struct {
	selectedMod *mods.Mod
	data        binding.UntypedList
	split       *container.Split
	mods        []*mods.Mod
	localMods   map[mods.ModID]bool
	prevSearch  string
	modList     *widget.List
}

func (ui *discoverUI) OnClose() {}

func (ui *discoverUI) PreDraw(w fyne.Window, args ...interface{}) (err error) {
	var (
		d      = dialog.NewInformation("", "Finding Mods...", w)
		lookup mods.ModLookup[*mods.Mod]
		ok     bool
	)
	defer d.Hide()
	d.Show()

	ui.prevSearch = ""

	ui.localMods = make(map[mods.ModID]bool)
	for _, tm := range args[0].([]interface{})[0].([]mods.TrackedMod) {
		ui.localMods[tm.ID()] = true
	}

	if lookup, err = discover.GetModsAsLookup(state.CurrentGame); err != nil {
		return
	}

	ui.mods = make([]*mods.Mod, 0, lookup.Len())
	for _, m := range lookup.All() {
		if _, ok = ui.localMods[m.ID()]; !ok {
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
		w.SetContent(container.NewBorder(
			container.NewVBox(
				widget.NewLabelWithStyle(string(state.CurrentGame.Name()), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
			), nil, nil, nil, container.NewBorder(
				container.NewAdaptiveGrid(8, container.NewHBox(widget.NewButton("Back", func() {
					if isPopup {
						state.ClosePopupWindow()
					} else {
						state.ShowPreviousScreen()
					}
				}))), nil, nil, nil,
				container.NewCenter(widget.NewLabel("No mods found")))))
		return
	}
	ui.data = binding.NewUntypedList()
	ui.modList = widget.NewListWithData(
		ui.data,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item binding.DataItem, co fyne.CanvasObject) {
			var m *mods.Mod
			if i, ok := cw.GetValueFromDataItem(item); ok {
				if m, ok = i.(*mods.Mod); ok {
					co.(*widget.Label).SetText(string(m.Name))
				}
			}
		})
	if err := ui.showSorted(ui.mods); err != nil {
		util.ShowErrorLong(err)
		return
	}

	ui.split = container.NewHSplit(ui.modList, container.NewMax())
	ui.split.SetOffset(0.25)

	ui.modList.OnSelected = func(id widget.ListItemID) {
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
				if err := managed.AddMod(state.CurrentGame, mods.NewTrackerMod(mod, state.CurrentGame)); err != nil {
					util.ShowErrorLong(err)
					return
				}
				for i, m := range ui.mods {
					if m == mod {
						ui.mods = append(ui.mods[:i], ui.mods[i+1:]...)
						break
					}
				}
				filtered := ui.applyFilters(ui.mods)
				sl := make([]interface{}, len(filtered))
				for i, m := range filtered {
					sl[i] = m
				}
				if err := ui.data.Set(sl); err != nil {
					util.ShowErrorLong(err)
					return
				}
				ui.selectedMod = nil
				ui.modList.UnselectAll()
				ui.split.Trailing = container.NewMax()
				ui.split.Refresh()
				state.UpdateCurrentScreen()
			})), nil, nil, nil,
			mp.CreatePreview(ui.selectedMod))
		ui.split.Refresh()
	}

	searchTb := widget.NewEntry()
	searchTb.OnChanged = func(s string) {
		if err := ui.search(s); err != nil {
			util.ShowErrorLong(err)
		}
	}

	w.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(string(state.CurrentGame.Name()), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
		), nil, nil, nil, container.NewBorder(
			container.NewHBox(
				widget.NewButton("Back", func() {
					if isPopup {
						state.ClosePopupWindow()
					} else {
						state.ShowPreviousScreen()
					}
				}),
				//NewFilterButton(ui.filterCallback, w),
				container.New(layout.NewGridLayout(6),
					widget.NewLabelWithStyle("Search", fyne.TextAlignTrailing, fyne.TextStyle{}), searchTb,
					widget.NewLabelWithStyle("Category", fyne.TextAlignTrailing, fyne.TextStyle{}), newCategoryFilter(ui.filterCallback),
					widget.NewLabelWithStyle("Include", fyne.TextAlignTrailing, fyne.TextStyle{}), newIncludeFilter(ui.filterCallback))),
			nil, nil, nil,
			ui.split)))
}

func (ui *discoverUI) filterCallback() {
	m := ui.applyFilters(ui.mods)
	m = ui.applySearch(ui.prevSearch, m)
	_ = ui.showSorted(m)
}

func (ui *discoverUI) search(s string) error {
	m := ui.applyFilters(ui.mods)
	if s != ui.prevSearch {
		m = ui.applySearch(s, m)
	}
	return ui.showSorted(m)
}

func (ui *discoverUI) applyFilters(orig []*mods.Mod) (result []*mods.Mod) {
	result = make([]*mods.Mod, 0, len(orig))
	for _, m := range orig {
		if filters.category != nil && m.Category != *filters.category {
			continue
		}
		if filters.supportedKind == supported && !m.IsManuallyCreated {
			continue
		}
		if filters.supportedKind == unsupported && m.IsManuallyCreated {
			continue
		}
		result = append(result, m)
	}
	return result
}

func (ui *discoverUI) applySearch(s string, orig []*mods.Mod) (result []*mods.Mod) {
	if len(s) < 3 && ui.prevSearch == "" {
		s = ""
		if ui.data.Length() == len(ui.mods) {
			return orig
		}
	}

	s = strings.ToLower(s)
	ui.prevSearch = s

	for _, m := range orig {
		if strings.Contains(strings.ToLower(string(m.Name)), s) ||
			strings.Contains(strings.ToLower(string(m.Category)), s) ||
			strings.Contains(strings.ToLower(m.Description), s) ||
			strings.Contains(strings.ToLower(m.Author), s) {
			result = append(result, m)
		}
	}
	return result
}

func (ui *discoverUI) showSorted(ms []*mods.Mod) error {
	_ = ui.data.Set(nil)
	for _, m := range mods.Sort(ms) {
		if err := ui.data.Append(m); err != nil {
			return err
		}
	}
	return nil
}
