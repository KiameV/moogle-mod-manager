package author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/pr-modsync/mods"
)

func createCompatibilities(w fyne.Window) fyne.CanvasObject {
	var c fyne.CanvasObject
	c = container.NewVScroll(container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Requires", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {

			})),
	))
	return c
}

func showModCompatForm(w fyne.Window, parent fyne.CanvasObject, mc *mods.ModCompat) {
	if mc == nil {
		mc = &mods.ModCompat{}
	}
	setFormItem("compatModID", "")
	setFormItem("compatVersions", "")
	setFormMultiLine("compatSources", "")
	setFormSelect("compatOrder", []string{"", string(mods.Before), string(mods.After)}, "")

	dialog.ShowForm("Required Mod", "Add", "Cancel", []*widget.FormItem{
		getFormItem("Mod ID", "compatModID"),
		getFormItem("Versions (comma seperated)", "compatVersions"),
		getFormItem("Sources (newline seperated)", "compatSources"),
		getFormItem("Order (optional)", "compatOrder"),
	}, func(ok bool) {
		if ok {
			if mod.ModCompatibility == nil {
				mod.ModCompatibility = &mods.ModCompatibility{}
			}
			mod.ModCompatibility.Requires = append(mod.ModCompatibility.Requires, modCompat)
			parent.Refresh()
		}
	}, w)
}

func newModCompatForm() []*widget.FormItem {
	return []*widget.FormItem{}
}
