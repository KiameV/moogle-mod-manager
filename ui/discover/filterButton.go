package discover

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type supportedKind string

const (
	unsupported supportedKind = "Unverified"
	supported   supportedKind = "Verified"
	all         supportedKind = "All"
)

var supportedKinds = []string{string(unsupported), string(supported), string(all)}

var filters = findFilter{
	supportedKind: supported,
	category:      nil,
}

type findFilter struct {
	supportedKind supportedKind
	category      *mods.Category
}

/*
func NewFilterButton(callback func(bool), window fyne.Window) *filterButton {
	b := &filterButton{
		callback: callback,
		window:   window,
	}
	b.Text = "Filters"
	return b
}

type filterButton struct {
	widget.Button
	callback func(bool)
	window   fyne.Window
}

func (b *filterButton) Tapped(e *fyne.PointEvent) {
	options := make([]string, len(mods.Categories)+1)
	options[0] = ""
	for i, c := range mods.Categories {
		options[i+1] = c
	}
	category := widget.NewSelect(options, func(s string) {
		if s == "" {
			filters.category = nil
		} else {
			c := mods.Category(s)
			filters.category = &c
		}
	})
	if filters.category != nil {
		category.SetSelected(string(*filters.category))
	} else {
		category.SetSelected("")
	}

	include := widget.NewSelect(supportedKinds, func(s string) {
		filters.supportedKind = supportedKind(s)
	})
	include.SetSelected(string(filters.supportedKind))

	d := dialog.NewForm("Filters", "Apply", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Category", category),
		widget.NewFormItem("Show", include),
	}, b.callback, b.window)
	d.Resize(fyne.NewSize(300, 200))
	d.Show()
}
*/

func newCategoryFilter(onChange func()) fyne.CanvasObject {
	options := make([]string, len(mods.Categories)+1)
	options[0] = ""
	for i, c := range mods.Categories {
		options[i+1] = c
	}
	category := widget.NewSelect(options, func(s string) {
		if s == "" {
			filters.category = nil
		} else {
			c := mods.Category(s)
			filters.category = &c
		}
		onChange()
	})
	if filters.category != nil {
		category.SetSelected(string(*filters.category))
	} else {
		category.SetSelected("")
	}
	return category
}

func newIncludeFilter(onChange func()) fyne.CanvasObject {
	include := widget.NewSelect(supportedKinds, func(s string) {
		filters.supportedKind = supportedKind(s)
		onChange()
	})
	include.SetSelected(string(filters.supportedKind))
	return include
}
